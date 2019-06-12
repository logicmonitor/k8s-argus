package argus

import (
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/utilities"

	"github.com/logicmonitor/k8s-argus/pkg/watch/deployment"

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func newLMClient(id, key, company string) *client.LMSdkGo {
	config := client.NewConfig()
	config.SetAccessID(&id)
	config.SetAccessKey(&key)
	domain := company + ".logicmonitor.com"
	config.SetAccountDomain(&domain)
	//config.UserAgent = constants.UserAgentBase + constants.Version
	api := client.New(config)
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

	hasDeploymentRbac := checkHasDeploymentRbac(base)
	log.Infof("Has deployment rbac privilege: %+v", hasDeploymentRbac)

	// init sync to delete the non-exist resource devices through logicmonitor API
	initSyncer := sync.InitSyncer{
		DeviceManager: deviceManager,
	}
	initSyncer.InitSync(hasDeploymentRbac)

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

	if hasDeploymentRbac {
		argus.Watchers = append(argus.Watchers, &deployment.Watcher{
			DeviceManager: deviceManager,
		})
	}

	return argus, nil
}

func checkHasDeploymentRbac(base *types.Base) bool {
	clusterRole, err := base.K8sClient.RbacV1beta1().ClusterRoles().Get("argus", metav1.GetOptions{})
	if err != nil {
		log.Errorf("Get clusterRoles failed: %+v", err)
	} else {
		for _, rule := range clusterRole.Rules {
			if len(rule.APIGroups) == 0 {
				continue
			}
			hasApps, _ := utilities.Contains("apps", rule.APIGroups)
			if hasApps {
				hasDeployment, _ := utilities.Contains("deployments", rule.Resources)
				if hasDeployment {
					return true
				}
			}
		}
	}
	return false
}

// NewBase instantiates and returns the base structure used throughout Argus.
func NewBase(config *config.Config) (*types.Base, error) {
	// LogicMonitor API client.
	lmClient := newLMClient(config.ID, config.Key, config.Account)

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
	for _, w := range a.Watchers {
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
		stop := make(chan struct{})
		go controller.Run(stop)
	}
}

// get the K8s RESTClient by apiVersion, use the default V1 version if there is no match
func getK8sRESTClient(clientset *kubernetes.Clientset, apiVersion string) rest.Interface {
	switch apiVersion {
	case constants.K8sAPIVersionV1:
		return clientset.CoreV1().RESTClient()
	case constants.K8sAPIVersionAppsV1beta2:
		return clientset.AppsV1beta2().RESTClient()
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
