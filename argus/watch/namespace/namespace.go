package namespace

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/logicmonitor/k8s-argus/argus/config"
	"github.com/logicmonitor/k8s-argus/argus/constants"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "namespaces"
)

// Watcher represents a watcher type that watches namespaces.
type Watcher struct {
	LMClient     *lmv1.DefaultApi
	K8sClient    *kubernetes.Clientset
	Config       *config.Config
	DeviceGroups map[string]int32
}

// Resource is a function that implements the Watcher interface.
func (w Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w Watcher) ObjType() runtime.Object {
	return &v1.Namespace{}
}

// AddFunc is a function that implements the Watcher interface.
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		namespace := obj.(*v1.Namespace)
		for k, parentID := range w.DeviceGroups {
			var appliesTo string
			switch k {
			case "services":
				appliesTo = "hasCategory(\"" + constants.ServiceCategory + "\") && auto.namespace == \"" + namespace.Name + "\" && auto.clustername == \"" + w.Config.ClusterName + "\""
			case "pods":
				appliesTo = "hasCategory(\"" + constants.PodCategory + "\") && auto.namespace == \"" + namespace.Name + "\" && auto.clustername == \"" + w.Config.ClusterName + "\""
			default:
				return
			}

			parentDeviceGroup, err := w.findDeviceGroup(parentID)
			if err != nil {
				log.Printf("Failed searching for device group: %v", err)
			}
			deviceGroup, err := w.createDeviceGroup(namespace.Name, appliesTo, parentDeviceGroup.Id)
			if err != nil {
				log.Printf("Failed creating device group: %v", err)
			} else {
				log.Printf("Created device group with id %d", deviceGroup.Id)
			}
			log.Printf("Using device group with id %d for namespace %q", deviceGroup.Id, namespace.Name)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		// oldNamespace := oldObj.(*v1.Namespace)
		// newNamespace := newObj.(*v1.Namespace)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		// namespace := obj.(*v1.Namespace)
	}
}

func (w *Watcher) findDeviceGroup(parentID int32) (deviceGroup *lmv1.RestDeviceGroup, err error) {
	restDeviceGroupResponse, _, err := w.LMClient.GetDeviceGroupById(parentID, "")
	if err != nil {
		err = fmt.Errorf("Failed to find device group: %s", restDeviceGroupResponse.Errmsg)
		return
	}

	deviceGroup = &restDeviceGroupResponse.Data
	log.Debugf("%#v", restDeviceGroupResponse)

	return
}

func (w *Watcher) createDeviceGroup(name, appliesTo string, parentID int32) (deviceGroup *lmv1.RestDeviceGroup, err error) {
	log.Infof("Creating device group %q", name)
	restDeviceGroupResponse, _, err := w.LMClient.AddDeviceGroup(lmv1.RestDeviceGroup{
		Name:            name,
		Description:     "A dynamic device group for Kubernetes.",
		ParentId:        parentID,
		AppliesTo:       appliesTo,
		DisableAlerting: w.Config.DisableAlerting,
	})
	if err != nil {
		err = fmt.Errorf("Failed to add device group %q", restDeviceGroupResponse.Errmsg)
		return
	}

	deviceGroup = &restDeviceGroupResponse.Data
	log.Infof("Created device group with id %d", deviceGroup.Id)

	return
}
