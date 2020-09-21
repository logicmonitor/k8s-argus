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
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/etcd"
	"github.com/logicmonitor/k8s-argus/pkg/facade"
	"github.com/logicmonitor/k8s-argus/pkg/lmexec"
	"github.com/logicmonitor/k8s-argus/pkg/sync"
	"github.com/logicmonitor/k8s-argus/pkg/tree"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/watch/deployment"
	"github.com/logicmonitor/k8s-argus/pkg/watch/namespace"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/logicmonitor/k8s-argus/pkg/worker"
	log "github.com/sirupsen/logrus"
	"github.com/vkumbhar94/lm-sdk-go/client"
	"github.com/vkumbhar94/lm-sdk-go/client/lm"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Argus represents the Argus cli.
type Argus struct {
	*types.Base
	types.LMFacade
	types.DeviceManager
	Watchers []types.Watcher
}

func newLMClient(argusConfig *config.Config) (*client.LMSdkGo, error) {
	config := client.NewConfig()
	config.SetAccessID(&argusConfig.ID)
	config.SetAccessKey(&argusConfig.Key)
	domain := argusConfig.Account + ".logicmonitor.com"
	config.SetAccountDomain(&domain)
	//config.UserAgent = constants.UserAgentBase + constants.Version
	if argusConfig.ProxyURL == "" {
		return client.New(config), nil
	}
	return newLMClientWithProxy(config, argusConfig)
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
	log.Infof("Using http/s proxy: %s", argusConfig.ProxyURL)
	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	transport := httptransport.NewWithClient(config.TransportCfg.Host, config.TransportCfg.BasePath, config.TransportCfg.Schemes, &httpClient)
	authInfo := client.LMv1Auth(*config.AccessID, *config.AccessKey)
	client := new(client.LMSdkGo)
	client.Transport = transport
	client.LM = lm.New(transport, strfmt.Default, authInfo)
	return client, nil
}

func newK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.UserAgent = constants.UserAgentBase + constants.Version

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// NewArgus instantiates and returns argus.
func NewArgus(base *types.Base) (*Argus, error) {
	facadeObj := facade.NewFacade()
	argus := &Argus{
		Base:     base,
		LMFacade: facadeObj,
	}
	lmExecObj := &lmexec.LMExec{
		Base: base,
	}

	dcache := devicecache.NewDeviceCache(base, 5)
	dcache.Run()

	deviceManager := &device.Manager{
		Base:       base,
		LMExecutor: lmExecObj,
		LMFacade:   facadeObj,
		DC:         dcache,
	}
	argus.DeviceManager = deviceManager

	deviceTree := &tree.DeviceTree{
		Base: base,
	}

	deviceGroups, err := deviceTree.CreateDeviceTree()
	if err != nil {
		return nil, err
	}

	podChannel := make(chan types.ICommand)
	serviceChannel := make(chan types.ICommand)
	deploymentChannel := make(chan types.ICommand)
	nodeChannel := make(chan types.ICommand)
	argus.Watchers = []types.Watcher{
		&namespace.Watcher{
			Base:         base,
			DeviceGroups: deviceGroups,
		},
		&node.Watcher{
			DeviceManager: deviceManager,
			DeviceGroups:  deviceGroups,
			LMClient:      base.LMClient,
			WConfig: &types.WConfig{
				MethodChannels: map[string]chan types.ICommand{
					"GET":    nodeChannel,
					"POST":   nodeChannel,
					"DELETE": nodeChannel,
					"PUT":    nodeChannel,
					"PATCH":  nodeChannel,
				},
				RetryLimit: 2,
				ID:         "nodes",
			},
		},
		&service.Watcher{
			DeviceManager: deviceManager,
			WConfig: &types.WConfig{
				MethodChannels: map[string]chan types.ICommand{
					"GET":    serviceChannel,
					"POST":   serviceChannel,
					"DELETE": serviceChannel,
					"PUT":    serviceChannel,
					"PATCH":  serviceChannel,
				},
				RetryLimit: 2,
				ID:         "services",
			},
		},
		&pod.Watcher{
			DeviceManager: deviceManager,
			WConfig: &types.WConfig{
				MethodChannels: map[string]chan types.ICommand{
					"GET":    podChannel,
					"POST":   podChannel,
					"DELETE": podChannel,
					"PUT":    podChannel,
					"PATCH":  podChannel,
				},
				RetryLimit: 2,
				ID:         "pods",
			},
		},
		&deployment.Watcher{
			DeviceManager: deviceManager,
			WConfig: &types.WConfig{
				MethodChannels: map[string]chan types.ICommand{
					"GET":    deploymentChannel,
					"POST":   deploymentChannel,
					"DELETE": deploymentChannel,
					"PUT":    deploymentChannel,
					"PATCH":  deploymentChannel,
				},
				RetryLimit: 2,
				ID:         "deployments",
			},
		},
	}

	// Start workers
	for _, w := range argus.Watchers {
		c := w.GetConfig()
		if c == nil {
			log.Warningf("Watcher %v doesn't have worker config, couldn't run worker for it", w.Resource())
			continue
		}
		wc := worker.NewWorker(c)
		b, err := argus.LMFacade.RegisterWorker(w.Resource(), wc)
		if err != nil {
			log.Errorf("Failed to register worker for resource for: %s", w.Resource())
		}
		if b {
			wc.Run()
		}
	}
	// init sync to delete the non-exist resource devices through logicmonitor API
	initSyncer := sync.InitSyncer{
		DeviceManager: deviceManager,
	}
	initSyncer.InitSync()

	if base.Config.EtcdDiscoveryToken != "" {
		etcdController := etcd.Controller{
			DeviceManager: deviceManager,
		}
		_, err = etcdController.DiscoverByToken()
		if err != nil {
			return nil, err
		}
	}
	log.Debugf("Initialized argus")
	//	podChannel := make(chan types.ICommand)
	//	serviceChannel := make(chan types.ICommand)
	//	deploymentChannel := make(chan types.ICommand)
	//	nodeChannel := make(chan types.ICommand)
	//	argus.Watchers = []types.Watcher{
	//		&namespace.Watcher{
	//			Base:         base,
	//			DeviceGroups: deviceGroups,
	//		},
	//		&node.Watcher{
	//			DeviceManager: deviceManager,
	//			DeviceGroups:  deviceGroups,
	//			LMClient:      base.LMClient,
	//			WConfig: types.WConfig{
	//				MethodChannels: map[string]chan types.ICommand{
	//					"GET":    nodeChannel,
	//					"POST":   nodeChannel,
	//					"DELETE": nodeChannel,
	//					"PUT":    nodeChannel,
	//					"PATCH":  nodeChannel,
	//				},
	//				RetryLimit: 2,
	//			},
	//		},
	//		&service.Watcher{
	//			DeviceManager: deviceManager,
	//			WConfig: types.WConfig{
	//				MethodChannels: map[string]chan types.ICommand{
	//					"GET":    serviceChannel,
	//					"POST":   serviceChannel,
	//					"DELETE": serviceChannel,
	//					"PUT":    serviceChannel,
	//					"PATCH":  serviceChannel,
	//				},
	//				RetryLimit: 2,
	//			},
	//		},
	//		&pod.Watcher{
	//			DeviceManager: deviceManager,
	//			WConfig: types.WConfig{
	//				MethodChannels: map[string]chan types.ICommand{
	//					"GET":    podChannel,
	//					"POST":   podChannel,
	//					"DELETE": podChannel,
	//					"PUT":    podChannel,
	//					"PATCH":  podChannel,
	//				},
	//				RetryLimit: 2,
	//			},
	//		},
	//		&deployment.Watcher{
	//			DeviceManager: deviceManager,
	//			WConfig: types.WConfig{
	//				MethodChannels: map[string]chan types.ICommand{
	//					"GET":    deploymentChannel,
	//					"POST":   deploymentChannel,
	//					"DELETE": deploymentChannel,
	//					"PUT":    deploymentChannel,
	//					"PATCH":  deploymentChannel,
	//				},
	//				RetryLimit: 2,
	//			},
	//		},
	//	}

	return argus, nil
}

// NewBase instantiates and returns the base structure used throughout Argus.
func NewBase(config *config.Config) (*types.Base, error) {
	// LogicMonitor API client.
	lmClient, err := newLMClient(config)
	if err != nil {
		return nil, err
	}

	// check and update the params
	checkAndUpdateClusterGroup(config, lmClient)

	// Kubernetes API client.
	k8sClient, err := newK8sClient()
	if err != nil {
		return nil, err
	}

	base := &types.Base{
		LMClient:  lmClient,
		K8sClient: k8sClient,
		Config:    config,
	}

	return base, nil
}

// Watch watches the API for events.
func (a *Argus) Watch() {
	log.Debugf("Starting watchers")
	for _, w := range a.Watchers {
		if !w.Enabled() {
			log.Warnf("Have no permission for resource %s", w.Resource())
			continue
		}
		watchlist := cache.NewListWatchFromClient(getK8sRESTClient(a.K8sClient, w.APIVersion()), w.Resource(), v1.NamespaceAll, fields.Everything())
		_, controller := cache.NewInformer(
			watchlist,
			w.ObjType(),
			time.Minute*10,
			cache.ResourceEventHandlerFuncs{
				AddFunc:    w.AddFunc(),
				DeleteFunc: w.DeleteFunc(),
				UpdateFunc: w.UpdateFunc(),
			},
		)
		log.Debugf("Starting watcher of %v", w.Resource())
		stop := make(chan struct{})
		go controller.Run(stop)
		//		c := w.GetConfig()
		//		if c == nil {
		//			continue
		//		}
		//		wc := worker.NewWorker(c)
		//		b, err := a.Facade.RegisterWorker(w.Resource(), wc)
		//		if err != nil {
		//			log.Errorf("Failed to register worker for resource for: %s", w.Resource())
		//		}
		//		if b {
		//			wc.StartWorker()
		//		}
	}
}

// get the K8s RESTClient by apiVersion, use the default V1 version if there is no match
func getK8sRESTClient(clientset *kubernetes.Clientset, apiVersion string) rest.Interface {
	switch apiVersion {
	case constants.K8sAPIVersionV1:
		return clientset.CoreV1().RESTClient()
	case constants.K8sAPIVersionAppsV1beta2:
		return clientset.AppsV1beta2().RESTClient()
	case constants.K8sAPIVersionAppsV1:
		return clientset.AppsV1().RESTClient()
	default:
		return clientset.CoreV1().RESTClient()
	}
}

// check the cluster group ID, if the group does not exist, just use the root group
func checkAndUpdateClusterGroup(config *config.Config, lmClient *client.LMSdkGo) {
	// do not need to check the root group
	if config.ClusterGroupID == constants.RootDeviceGroupID {
		return
	}

	// if the group does not exist anymore, we will add the cluster to the root group
	if !devicegroup.ExistsByID(config.ClusterGroupID, lmClient) {
		log.Warnf("The device group (id=%v) does not exist, the cluster will be added to the root group", config.ClusterGroupID)
		config.ClusterGroupID = constants.RootDeviceGroupID
	}
}
