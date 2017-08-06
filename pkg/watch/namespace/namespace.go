package namespace

import (
	"github.com/logicmonitor/k8s-argus/pkg/utilities"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	lm "github.com/logicmonitor/lm-sdk-go"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "namespaces"
)

// Watcher represents a watcher type that watches namespaces.
type Watcher struct {
	*types.Base
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
				logrus.Errorf("Failed to find namespace: %v", err)

				return
			}
			_, err = w.createDeviceGroup(namespace.Name, appliesTo, parentDeviceGroup.Id)
			if err != nil {
				logrus.Errorf("Failed to add namespace: %v", err)

				return
			}

			logrus.Printf("Added namespace %s", namespace.Name)
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
		namespace := obj.(*v1.Namespace)
		for name, parentID := range w.DeviceGroups {
			deviceGroup, err := w.findDeviceGroup(parentID)
			if err != nil {
				logrus.Printf("Failed to find namespace %s: %v", name, err)

				return
			}
			for _, subGroup := range deviceGroup.SubGroups {
				if subGroup.Name == namespace.Name {
					restResponse, apiResponse, err := w.LMClient.DeleteDeviceGroupById(subGroup.Id, true)
					if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
						logrus.Errorf("Failed to delete namespace %q: %v", subGroup.Name, _err)

						return
					}
				}
			}
		}
	}
}

func (w *Watcher) findDeviceGroup(parentID int32) (deviceGroup *lm.RestDeviceGroup, err error) {
	restResponse, apiResponse, err := w.LMClient.GetDeviceGroupById(parentID, "")
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		logrus.Errorf("Failed to find namespace: %v", _err)

		return
	}

	deviceGroup = &restResponse.Data
	logrus.Debugf("%#v", restResponse)

	return
}

func (w *Watcher) createDeviceGroup(name, appliesTo string, parentID int32) (deviceGroup *lm.RestDeviceGroup, err error) {
	restResponse, apiResponse, err := w.LMClient.AddDeviceGroup(lm.RestDeviceGroup{
		Name:            name,
		Description:     "A dynamic device group for Kubernetes.",
		ParentId:        parentID,
		AppliesTo:       appliesTo,
		DisableAlerting: w.Config.DisableAlerting,
	})
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		logrus.Errorf("Failed to add namespace %q: %v", name, _err)

		return
	}

	deviceGroup = &restResponse.Data

	return
}
