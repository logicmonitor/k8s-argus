package argus

import (
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/etcd"
	"github.com/logicmonitor/k8s-argus/pkg/tree"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/watch/namespace"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	lm "github.com/logicmonitor/lm-sdk-go"
	"k8s.io/api/core/v1"
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

func newLMClient(id, key, company string) *lm.DefaultApi {
	config := lm.NewConfiguration()
	config.APIKey = map[string]map[string]string{
		"Authorization": {
			"AccessID":  id,
			"AccessKey": key,
		},
	}
	config.BasePath = "https://" + company + ".logicmonitor.com/santaba/rest"
	config.UserAgent = constants.UserAgentBase + constants.Version

	api := lm.NewDefaultApi()
	api.Configuration = config

	return api
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
	lmClient := newLMClient(config.ID, config.Key, config.Account)

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
