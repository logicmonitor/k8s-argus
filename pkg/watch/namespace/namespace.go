package namespace

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/err"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	resource = "namespaces"
)

// Watcher represents a watcher type that watches namespaces.
type Watcher struct {
	*types.Base
	// TODO: This should be thread safe.
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
		// Due to panic error in this call stack will crash the application; recovering those panics here could make our application robust.
		defer err.RecoverError("Add namespace")
		namespace := obj.(*v1.Namespace)
		log.Debugf("Handling add namespace event: %s", namespace.Name)
		for name, parentID := range w.DeviceGroups {
			var appliesTo devicegroup.AppliesToBuilder
			// Ensure that we are creating namespaces for namespaced resources.
			switch name {
			case constants.ServiceDeviceGroupName:
				appliesTo = devicegroup.NewAppliesToBuilder().HasCategory(constants.ServiceCategory).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)
			case constants.PodDeviceGroupName:
				appliesTo = devicegroup.NewAppliesToBuilder().HasCategory(constants.PodCategory).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)
			default:
				continue
			}

			opts := &devicegroup.Options{
				AppliesTo:       appliesTo,
				Client:          w.LMClient,
				DisableAlerting: w.Config.DisableAlerting,
				Name:            namespace.Name,
				ParentID:        parentID,
			}

			log.Debugf("%v", opts)

			_, err := devicegroup.Create(opts)
			if err != nil {
				log.Errorf("Failed to add namespace to %q: %v", name, err)
				return
			}

			log.Infof("Added namespace %q to %q", namespace.Name, name)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		// Due to panic error in this call stack will crash the application; recovering those panics here could make our application robust.
		defer err.RecoverError("Update namespace")
		log.Debugf("Ignoring update namespace event")
		// oldNamespace := oldObj.(*v1.Namespace)
		// newNamespace := newObj.(*v1.Namespace)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		// Due to panic error in this call stack will crash the application; recovering those panics here could make our application robust.
		defer err.RecoverError("Delete namespace")
		namespace := obj.(*v1.Namespace)
		log.Debugf("Handle deleting namespace event: %s", namespace.Name)

		for name, parentID := range w.DeviceGroups {
			deviceGroup, err := devicegroup.Find(parentID, name, w.LMClient)
			if err != nil {
				log.Warnf("Failed to find namespace %s: %v", name, err)
				return
			}
			// We should only be returned a device group if it is namespaced.
			if deviceGroup == nil {
				continue
			}
			err = devicegroup.DeleteSubGroup(deviceGroup, namespace.Name, w.LMClient)
			if err != nil {
				log.Errorf("Failed to delete namespace %q: %v", namespace.Name, err)
				return
			}
		}

	}
}
