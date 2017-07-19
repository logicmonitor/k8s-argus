package argus

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/logicmonitor/argus/argus/config"
	"github.com/logicmonitor/argus/argus/watch"
	"github.com/logicmonitor/argus/argus/watch/namespace"
	"github.com/logicmonitor/argus/argus/watch/node"
	"github.com/logicmonitor/argus/argus/watch/pod"
	"github.com/logicmonitor/argus/argus/watch/service"
	"github.com/logicmonitor/argus/constants"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Argus represents the Argus cli.
type Argus struct {
	LMClient  *lmv1.DefaultApi
	K8sClient *kubernetes.Clientset
	Watchers  []watcher.Watcher
	Config    *config.Config
}

func newLMClient(id, key, company string) *lmv1.DefaultApi {
	config := lmv1.NewConfiguration()
	config.APIKey = map[string]map[string]string{
		"Authorization": map[string]string{
			"AccessID":  id,
			"AccessKey": key,
		},
	}
	config.BasePath = "https://" + company + ".logicmonitor.com/santaba/rest"

	api := lmv1.NewDefaultApi()
	api.Configuration = config

	return api
}

func newK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// NewArgus instantiates and returns argus.
func NewArgus(config *config.Config) (argus *Argus, err error) {
	// LogicMonitor API client.
	lmClient := newLMClient(config.ID, config.Key, config.Company)
	id, err := discoverCollectorID(lmClient, config.CollectorDescription)
	if err != nil {
		return nil, err
	}
	config.PreferredCollector = id
	log.Infof("Using collector %d", config.PreferredCollector)

	// Kubernetes API client.
	k8sClient, err := newK8sClient()
	if err != nil {
		return
	}

	argus = &Argus{
		LMClient:  lmClient,
		K8sClient: k8sClient,
		Config:    config,
	}

	deviceGroups, err := argus.CreateDeviceTree()
	if err != nil {
		return
	}

	argus.Watchers = []watcher.Watcher{
		namespace.Watcher{
			LMClient:     lmClient,
			K8sClient:    k8sClient,
			Config:       config,
			DeviceGroups: deviceGroups,
		},
		node.Watcher{
			LMClient:  lmClient,
			K8sClient: k8sClient,
			Config:    config,
		},
		service.Watcher{
			LMClient:  lmClient,
			K8sClient: k8sClient,
			Config:    config,
		},
		pod.Watcher{
			LMClient:  lmClient,
			K8sClient: k8sClient,
			Config:    config,
		},
	}

	return
}

// CreateDeviceTree creates the Device tree that will represent the cluster in the LogicMonitor portal.
func (a *Argus) CreateDeviceTree() (deviceGroups map[string]int32, err error) {
	deviceGroups, err = createDeviceGroups(a)
	if err != nil {
		log.Error(err)
	}

	return
}

// Watch watches the API for events.
func (a *Argus) Watch() {
	getter := a.K8sClient.Core().RESTClient()
	for _, w := range a.Watchers {
		watchlist := cache.NewListWatchFromClient(getter, w.Resource(), v1.NamespaceAll, fields.Everything())
		_, controller := cache.NewInformer(
			watchlist,
			w.ObjType(),
			time.Second*0,
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

func findDeviceGroup(client *lmv1.DefaultApi, parentID int32, name string) (deviceGroup *lmv1.RestDeviceGroup, err error) {
	// filter := fmt.Sprintf("parentId:%d", parentID)
	filter := fmt.Sprintf("name:%s", url.QueryEscape(name))
	restDeviceGroupPaginationResponse, _, err := client.GetDeviceGroupList("name,id,parentId", -1, 0, filter)
	if err != nil {
		err = fmt.Errorf("Failed to find device group: %s", restDeviceGroupPaginationResponse.Errmsg)
		return
	}

	log.Debugf("%#v", restDeviceGroupPaginationResponse)

	for _, d := range restDeviceGroupPaginationResponse.Data.Items {
		if d.ParentId == parentID {
			log.Infof("Found device group %q with id %d", name, parentID)
			deviceGroup = &d

			return
		}
	}

	return
}

func createDeviceGroup(client *lmv1.DefaultApi, name, appliesTo string, disableAlerting bool, parentID int32) (deviceGroup *lmv1.RestDeviceGroup, err error) {
	log.Infof("Creating device group %q", name)
	restDeviceGroupResponse, _, err := client.AddDeviceGroup(lmv1.RestDeviceGroup{
		Name:            name,
		Description:     "A dynamic device group for Kubernetes.",
		ParentId:        parentID,
		AppliesTo:       appliesTo,
		DisableAlerting: disableAlerting,
	})
	if err != nil {
		err = fmt.Errorf("Failed to add device group %q", restDeviceGroupResponse.Errmsg)
		return
	}

	deviceGroup = &restDeviceGroupResponse.Data
	log.Infof("Created device group with id %d", deviceGroup.Id)

	return
}

func createClusterDeviceGroup(argus *Argus) (clusterDeviceGroup *lmv1.RestDeviceGroup, err error) {
	name := "Kubernetes Cluster: " + argus.Config.ClusterName
	appliesTo := "hasCategory(\"" + constants.ClusterCategory + "\") && auto.clustername ==\"" + argus.Config.ClusterName + "\""

	clusterDeviceGroup, err = findDeviceGroup(argus.LMClient, 1, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = createDeviceGroup(argus.LMClient, name, appliesTo, argus.Config.DisableAlerting, 1)
		if err != nil {
			return
		}
	}

	return
}

func createServiceDeviceGroup(argus *Argus, parentDeviceGroup *lmv1.RestDeviceGroup) (clusterDeviceGroup *lmv1.RestDeviceGroup, err error) {
	name := "Services"
	appliesTo := ""

	clusterDeviceGroup, err = findDeviceGroup(argus.LMClient, parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = createDeviceGroup(argus.LMClient, name, appliesTo, argus.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func createServiceDeletedDeviceGroup(argus *Argus, parentDeviceGroup *lmv1.RestDeviceGroup) (clusterDeviceGroup *lmv1.RestDeviceGroup, err error) {
	name := "_deleted"
	appliesTo := "hasCategory(\"" + constants.NodeDeletedCategory + "\") && auto.clustername ==\"" + argus.Config.ClusterName + "\""

	clusterDeviceGroup, err = findDeviceGroup(argus.LMClient, parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = createDeviceGroup(argus.LMClient, name, appliesTo, argus.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func createNodeDeviceGroup(argus *Argus, parentDeviceGroup *lmv1.RestDeviceGroup) (clusterDeviceGroup *lmv1.RestDeviceGroup, err error) {
	name := "Nodes"
	appliesTo := "hasCategory(\"" + constants.NodeCategory + "\") && auto.clustername ==\"" + argus.Config.ClusterName + "\""

	clusterDeviceGroup, err = findDeviceGroup(argus.LMClient, parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = createDeviceGroup(argus.LMClient, name, appliesTo, argus.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func createNodeDeletedDeviceGroup(argus *Argus, parentDeviceGroup *lmv1.RestDeviceGroup) (clusterDeviceGroup *lmv1.RestDeviceGroup, err error) {
	name := "_deleted"
	appliesTo := "hasCategory(\"" + constants.NodeDeletedCategory + "\") && auto.clustername ==\"" + argus.Config.ClusterName + "\""

	clusterDeviceGroup, err = findDeviceGroup(argus.LMClient, parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = createDeviceGroup(argus.LMClient, name, appliesTo, true, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func createPodDeviceGroup(argus *Argus, parentDeviceGroup *lmv1.RestDeviceGroup) (clusterDeviceGroup *lmv1.RestDeviceGroup, err error) {
	name := "Pods"
	appliesTo := ""

	clusterDeviceGroup, err = findDeviceGroup(argus.LMClient, parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = createDeviceGroup(argus.LMClient, name, appliesTo, argus.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func createPodDeletedDeviceGroup(argus *Argus, parentDeviceGroup *lmv1.RestDeviceGroup) (clusterDeviceGroup *lmv1.RestDeviceGroup, err error) {
	name := "_deleted"
	appliesTo := "hasCategory(\"" + constants.PodDeletedCategory + "\") && auto.clustername ==\"" + argus.Config.ClusterName + "\""

	clusterDeviceGroup, err = findDeviceGroup(argus.LMClient, parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = createDeviceGroup(argus.LMClient, name, appliesTo, true, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func createDeviceGroups(argus *Argus) (deviceGroups map[string]int32, err error) {
	deviceGroups = make(map[string]int32)
	clusterDeviceGroup, err := createClusterDeviceGroup(argus)
	if err != nil {
		return
	}
	log.Infof("Using cluster device group with id %d", clusterDeviceGroup.Id)

	serviceDeviceGroup, err := createServiceDeviceGroup(argus, clusterDeviceGroup)
	if err != nil {
		return
	}
	deviceGroups["services"] = serviceDeviceGroup.Id
	log.Infof("Using service device group with id %d", serviceDeviceGroup.Id)
	_, err = createServiceDeletedDeviceGroup(argus, serviceDeviceGroup)
	if err != nil {
		return
	}

	nodeDeviceGroup, err := createNodeDeviceGroup(argus, clusterDeviceGroup)
	if err != nil {
		return
	}
	log.Infof("Using node device group with id %d", nodeDeviceGroup.Id)
	_, err = createNodeDeletedDeviceGroup(argus, nodeDeviceGroup)
	if err != nil {
		return
	}

	podDeviceGroup, err := createPodDeviceGroup(argus, clusterDeviceGroup)
	if err != nil {
		return
	}
	deviceGroups["pods"] = podDeviceGroup.Id
	log.Infof("Using pod device group with id %d", podDeviceGroup.Id)
	_, err = createPodDeletedDeviceGroup(argus, podDeviceGroup)
	if err != nil {
		return
	}

	return
}

func discoverCollectorID(client *lmv1.DefaultApi, description string) (int32, error) {
	attempts := 0
	for {
		restResponse, apiResponse, err := client.GetCollectorList("id", -1, 0, "description:"+description)
		if err != nil {
			return -1, err
		}
		if apiResponse.StatusCode != http.StatusOK {
			return -1, fmt.Errorf(apiResponse.Status)
		}
		if restResponse.Status != http.StatusOK {
			return -1, fmt.Errorf(restResponse.Errmsg)
		}
		switch restResponse.Data.Total {
		case 0:
			if attempts == 6 {
				return -1, fmt.Errorf("Timeout waiting for collector ID")
			}
			log.Infof("No collector found, waiting 10 seconds...")
			time.Sleep(10 * time.Second)
		case 1:
			return restResponse.Data.Items[0].Id, nil
		default:
			return -1, fmt.Errorf("Found %d collectors with description %q", restResponse.Data.Total, description)
		}
	}
}
