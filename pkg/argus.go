package argus

import (
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/logicmonitor/lm-sdk-go/client/lm"

	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/etcd"
	"github.com/logicmonitor/k8s-argus/pkg/sync"
	"github.com/logicmonitor/k8s-argus/pkg/tree"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/watch/namespace"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	"github.com/logicmonitor/lm-sdk-go/client"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Argus represents the Argus cli.
type Argus struct {
	*types.Base
	Watchers []types.Watcher
}

func newLMClient(id, key, company string, proxyUrl string) (*client.LMSdkGo, error) {
	config := client.NewConfig()
	config.SetAccessID(&id)
	config.SetAccessKey(&key)
	domain := company + ".logicmonitor.com"
	config.SetAccountDomain(&domain)
	//config.UserAgent = constants.UserAgentBase + constants.Version
	log.Infof("proxyUrl = %s", proxyUrl)
	if proxyUrl == "" {
		return client.New(config), nil
	} else {
		return newLMClientProxy(config, proxyUrl)

	}
}

func newLMClientProxy(config *client.Config, proxyUrlStr string) (*client.LMSdkGo, error) {
	log.Infof("Use http proxy: %s", proxyUrlStr)
	proxyUrl, err := url.Parse(proxyUrlStr)
	if err != nil {
		return nil, err
	}
	// set proxy username and password
	//proxyUrl.User = url.UserPassword("","")

	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	transport := httptransport.NewWithClient(config.TransportCfg.Host, config.TransportCfg.BasePath, config.TransportCfg.Schemes, &httpClient)
	authInfo := client.LMv1Auth(*config.AccessID, *config.AccessKey)
	cli := new(client.LMSdkGo)
	cli.Transport = transport
	cli.LM = lm.New(transport, strfmt.Default, authInfo)
	return cli, nil
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
func NewArgus(base *types.Base, client api.CollectorSetControllerClient) (*Argus, error) {
	argus := &Argus{
		Base: base,
	}

	deviceManager := &device.Manager{
		Base:             base,
		ControllerClient: client,
	}

	deviceTree := &tree.DeviceTree{
		Base: base,
	}

	// init sync to delete the non-exist resource devices through logicmonitor API
	initSyncer := sync.InitSyncer{
		DeviceManager: deviceManager,
	}
	initSyncer.InitSync()

	deviceGroups, err := deviceTree.CreateDeviceTree()
	if err != nil {
		return nil, err
	}

	if base.Config.EtcdDiscoveryToken != "" {
		etcdController := etcd.Controller{
			DeviceManager: deviceManager,
		}
		_, err = etcdController.DiscoverByToken()
		if err != nil {
			return nil, err
		}
	}

	argus.Watchers = []types.Watcher{
		&namespace.Watcher{
			Base:         base,
			DeviceGroups: deviceGroups,
		},
		&node.Watcher{
			DeviceManager: deviceManager,
			DeviceGroups:  deviceGroups,
			LMClient:      base.LMClient,
		},
		&service.Watcher{
			DeviceManager: deviceManager,
		},
		&pod.Watcher{
			DeviceManager: deviceManager,
		},
	}

	return argus, nil
}

// NewBase instantiates and returns the base structure used throughout Argus.
func NewBase(config *config.Config) (*types.Base, error) {
	// LogicMonitor API client.
	lmClient, err := newLMClient(config.ID, config.Key, config.Account, config.ProxyUrl)
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
	getter := a.K8sClient.CoreV1().RESTClient()
	for _, w := range a.Watchers {
		watchlist := cache.NewListWatchFromClient(getter, w.Resource(), v1.NamespaceAll, fields.Everything())
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
		stop := make(chan struct{})
		go controller.Run(stop)
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
