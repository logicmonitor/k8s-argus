package argus

import (
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	cache2 "github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/etcd"
	"github.com/logicmonitor/k8s-argus/pkg/facade"
	"github.com/logicmonitor/k8s-argus/pkg/lmexec"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/sync"
	"github.com/logicmonitor/k8s-argus/pkg/tree"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/k8s-argus/pkg/watch/namespace"
	"github.com/logicmonitor/k8s-argus/pkg/watch/resource"
	"github.com/logicmonitor/k8s-argus/pkg/worker"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

// Argus represents the Argus cli.
type Argus struct {
	*types.Base
	types.LMFacade
	types.DeviceManager
	Watchers               []types.ResourceWatcher
	controllerStateHolders map[enums.ResourceType]*types.ControllerInitSyncStateHolder
	NSWatcher              *namespace.OldWatcher
}

func newLMClient(argusConfig *config.Config) (*client.LMSdkGo, error) {
	conf := client.NewConfig()
	conf.SetAccessID(&argusConfig.ID)
	conf.SetAccessKey(&argusConfig.Key)
	domain := argusConfig.Account + ".logicmonitor.com"
	conf.SetAccountDomain(&domain)
	// conf.UserAgent = constants.UserAgentBase + constants.Version
	if argusConfig.ProxyURL == "" {
		if argusConfig.IgnoreSSL {
			return newLMClientWithoutSSL(conf)
		}

		return client.New(conf), nil
	}

	return newLMClientWithProxy(conf, argusConfig)
}

func newLMClientWithProxy(config *client.Config, argusConfig *config.Config) (*client.LMSdkGo, error) {
	proxyURL, err := url.Parse(argusConfig.ProxyURL)
	if err != nil {
		return nil, err
	}
	if argusConfig.ProxyUser != "" {
		if argusConfig.ProxyPass != "" {
			proxyURL.User = url.UserPassword(argusConfig.ProxyUser, argusConfig.ProxyPass)
		} else {
			proxyURL.User = url.User(argusConfig.ProxyUser)
		}
	}
	logrus.Infof("Using http/s proxy: %s", argusConfig.ProxyURL)
	httpClient := http.Client{
		Transport: &http.Transport{ // nolint: exhaustivestruct
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	transport := httptransport.NewWithClient(config.TransportCfg.Host, config.TransportCfg.BasePath, config.TransportCfg.Schemes, &httpClient)
	authInfo := client.LMv1Auth(*config.AccessID, *config.AccessKey)
	clientObj := new(client.LMSdkGo)
	clientObj.Transport = transport
	clientObj.LM = lm.New(transport, strfmt.Default, authInfo)

	return clientObj, nil
}

func newLMClientWithoutSSL(config *client.Config) (*client.LMSdkGo, error) {
	opts := httptransport.TLSClientOptions{InsecureSkipVerify: true}
	httpClient, err := httptransport.TLSClient(opts)
	if err != nil {
		return nil, err
	}
	transport := httptransport.NewWithClient(config.TransportCfg.Host, config.TransportCfg.BasePath, config.TransportCfg.Schemes, httpClient)
	authInfo := client.LMv1Auth(*config.AccessID, *config.AccessKey)
	cli := new(client.LMSdkGo)
	cli.Transport = transport
	cli.LM = lm.New(transport, strfmt.Default, authInfo)

	return cli, nil
}

// NewArgus instantiates and returns argus.
func NewArgus(base *types.Base) (*Argus, error) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"function": "new-argus"}))
	facadeObj := facade.NewFacade()
	argus := &Argus{ // nolint: exhaustivestruct
		Base:                   base,
		LMFacade:               facadeObj,
		controllerStateHolders: make(map[enums.ResourceType]*types.ControllerInitSyncStateHolder),
	}
	lmExecObj := &lmexec.LMExec{
		Base: base,
	}

	resourceCache := devicecache.NewResourceCache(base, base.Config.GetCacheSyncInterval())
	resourceCache.Run()

	deviceManager := &device.Manager{ // nolint: exhaustivestruct
		Base:          base,
		LMExecutor:    lmExecObj,
		LMFacade:      facadeObj,
		ResourceCache: resourceCache,
	}
	argus.DeviceManager = deviceManager

	deviceTree := &tree.DeviceTree{
		Base:          base,
		ResourceCache: resourceCache,
	}

	deviceGroups, err := deviceTree.CreateDeviceTree(lctx)
	if err != nil {
		return nil, err
	}

	// Graceful rebuild
	if devicegroup.GetClusterGroupProperty(lctx, constants.ResyncCacheProp, base.LMClient) == "true" {
		resourceCache.Rebuild(lctx)
		clusterGroupID := util.GetClusterGroupID(lctx, base.LMClient)
		devicegroup.DeleteDeviceGroupPropertyByName(lctx, clusterGroupID, &models.EntityProperty{Name: constants.ResyncCacheProp, Value: "true"}, base.LMClient)
	}

	storeDeviceGroupsInCache(base, deviceGroups, resourceCache)
	logrus.Infof("Device group tree: %v", deviceGroups)

	deviceTree2 := &tree.DeviceTree2{
		Base:          base,
		ResourceCache: resourceCache,
	}
	deviceGroups2, err2 := deviceTree2.CreateDeviceTree(lctx)
	if err2 != nil {
		return nil, err2
	}
	logrus.Infof("New Device group tree : %v", deviceGroups2)

	createWatchers(argus)

	argus.NSWatcher = &namespace.OldWatcher{
		Resource:      enums.Namespaces,
		Base:          base,
		DeviceGroups:  deviceGroups,
		DeviceGroups2: deviceGroups2,
		ResourceCache: resourceCache,
	}

	startWorkers(argus)
	// init sync to delete the non-exist resource devices through LogicMonitor API
	initSyncer := sync.InitSyncer{
		DeviceManager: deviceManager,
	}

	lctx2 := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "init-sync"}))
	initSyncer.Sync(lctx2)

	// periodically delete the non-exist resource devices through logicmonitor API based on specified time interval.
	initSyncer.RunPeriodicSync(base.Config.GetPeriodicDeleteInterval())

	a, err3 := discoverETCDNodes(base, deviceManager)
	if err3 != nil {
		return a, err3
	}
	logrus.Debugf("Initialized argus")

	return argus, nil
}

func discoverETCDNodes(base *types.Base, deviceManager *device.Manager) (*Argus, error) {
	if base.Config.EtcdDiscoveryToken != "" {
		etcdController := etcd.Controller{
			DeviceManager: deviceManager,
		}
		_, err := etcdController.DiscoverByToken()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func startWorkers(argus *Argus) {
	// Start workers
	for _, w := range argus.Watchers {
		resourceType := w.ResourceType()
		c := w.GetConfig()
		if c == nil {
			logrus.Warningf("Watcher %v doesn't have worker config, couldn't run worker for it", w)

			continue
		}
		wc := worker.NewWorker(c)
		b, err := argus.LMFacade.RegisterWorker(c.ID, wc)
		if err != nil {
			logrus.Errorf("Failed to register worker for resource for: %s", resourceType.String())
		}
		if b {
			wc.Run()
		}
	}
}

func createWatchers(argus *Argus) {
	for _, rt := range enums.ALLResourceTypes {
		// These need special handling
		if rt == enums.Namespaces || rt == enums.ETCD {
			continue
		}
		argus.Watchers = append(argus.Watchers, &resource.Watcher{Resource: rt, WConfig: types.NewHTTPWConfig(rt)})
	}
}

func storeDeviceGroupsInCache(base *types.Base, deviceGroups map[string]int32, resourceCache *devicecache.ResourceCache) {
	// TODO: Temporary added here, move to its respective place
	for k, v := range deviceGroups {
		parent := util.ClusterGroupName(base.Config.ClusterName)
		switch k {
		case constants.AllNodeDeviceGroupName:
			// the all nodes group should be nested in 'Nodes'
			parent = constants.NodeDeviceGroupName
		case util.ClusterGroupName(base.Config.ClusterName):
			parent = ""
		}

		resourceCache.Set(cache2.ResourceName{
			Name:     k,
			Resource: enums.Namespaces,
		}, cache2.ResourceMeta{ // nolint: exhaustivestruct
			Container: parent,
			LMID:      v,
		})
	}
}

// NewBase instantiates and returns the base structure used throughout Argus.
func NewBase(conf *config.Config) (*types.Base, error) {
	// LogicMonitor API client.
	lmClient, err := newLMClient(conf)
	if err != nil {
		return nil, err
	}

	// check and update the params
	checkAndUpdateClusterGroup(conf, lmClient)

	// Kubernetes API client.
	k8sClient := config.GetClientSet()

	base := &types.Base{
		LMClient:  lmClient,
		K8sClient: k8sClient,
		Config:    conf,
	}

	return base, nil
}

// Watch watches the API for events.
func (a *Argus) Watch() {
	syncInterval := a.Base.Config.GetPeriodicSyncInterval()
	logrus.Debugf("Starting watchers")
	b := &builder.Builder{}

	nsRT, controller := a.RunNSWatcher(syncInterval)
	logrus.Debugf("Starting ns watcher of %v", nsRT.String())
	stop := make(chan struct{})
	go controller.Run(stop)

	for _, w := range a.Watchers {
		rt := w.ResourceType()
		// TODO: has permission and check for enabled flag in case if user wants to avoid all resource of specific type
		//  earlier all resources used to ignore from filter config but still it used to put pressure on k8s api-server to unnecessary polls
		if !permission.HasPermissions(rt) {
			logrus.Warnf("Have no permission for resource %s", rt.String())

			continue
		}
		watchlist := cache.NewListWatchFromClient(util.GetK8sRESTClient(a.K8sClient, rt.K8SAPIVersion()), rt.String(), corev1.NamespaceAll, fields.Everything())
		controller := a.createNewInformer(watchlist, rt, syncInterval, b)
		logrus.Debugf("Starting watcher of %v", rt.String())
		stop := make(chan struct{})
		stateHolder := types.NewControllerInitSyncStateHolder(controller)
		stateHolder.Run()
		a.controllerStateHolders[rt] = &stateHolder
		go controller.Run(stop)
	}
}

func (a *Argus) createNewInformer(watchlist cache.ListerWatcher, rt enums.ResourceType, syncInterval time.Duration, b *builder.Builder) cache.Controller {
	_, controller := cache.NewInformer(
		watchlist,
		rt.K8SObjectType(),
		syncInterval,
		cache.ResourceEventHandlerFuncs{
			AddFunc: resource.AddFuncDispatcher(
				resource.AddFuncWithExclude(
					resource.AddOrUpdateFunc(
						a.controllerStateHolders,
						b.AddFuncWithDefaults(
							a.DeviceManager.GetResourceCache(),
							resource.WatcherConfigurer(rt),
							a.DeviceManager,
						),
						b.UpdateFuncWithDefaults(
							resource.UpsertBasedOnCache(
								a.DeviceManager.GetResourceCache(),
								resource.WatcherConfigurer(rt),
								a.DeviceManager,
								b,
							),
						),
					),
					b.DeleteFuncWithDefaults(resource.WatcherConfigurer(rt), a.DeviceManager.DeleteFunc()),
				),
			),
			UpdateFunc: resource.UpdateFuncDispatcher(
				resource.UpdateFuncWithExclude(
					b.UpdateFuncWithDefaults(
						resource.UpsertBasedOnCache(
							a.DeviceManager.GetResourceCache(),
							resource.WatcherConfigurer(rt),
							a.DeviceManager,
							b,
						),
					),
					b.DeleteFuncWithDefaults(resource.WatcherConfigurer(rt), a.DeviceManager.DeleteFunc()),
				),
			),
			DeleteFunc: resource.DeleteFuncDispatcher(
				b.DeleteFuncWithDefaults(resource.WatcherConfigurer(rt), a.DeviceManager.DeleteFunc()),
			),
		},
	)
	return controller
}

func (a *Argus) RunNSWatcher(syncInterval time.Duration) (enums.ResourceType, cache.Controller) {
	nsRT := a.NSWatcher.ResourceType()
	// start ns watcher
	watchlist := cache.NewListWatchFromClient(util.GetK8sRESTClient(a.K8sClient, nsRT.K8SAPIVersion()), nsRT.String(), corev1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		nsRT.K8SObjectType(),
		syncInterval,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    a.NSWatcher.AddFunc(),
			UpdateFunc: a.NSWatcher.UpdateFunc(),
			DeleteFunc: a.NSWatcher.DeleteFunc(),
		},
	)
	return nsRT, controller
}

// check the cluster group ID, if the group does not exist, just use the root group
func checkAndUpdateClusterGroup(config *config.Config, lmClient *client.LMSdkGo) {
	// do not need to check the root group
	if config.ClusterGroupID == constants.RootDeviceGroupID {
		return
	}

	// if the group does not exist anymore, we will add the cluster to the root group
	if !devicegroup.ExistsByID(config.ClusterGroupID, lmClient) {
		logrus.Warnf("The device group (id=%v) does not exist, the cluster will be added to the root group", config.ClusterGroupID)
		config.ClusterGroupID = constants.RootDeviceGroupID
	}
}
