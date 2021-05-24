package namespace

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	resource = "namespaces"
)

// Watcher represents a watcher type that watches namespaces.
type Watcher struct {
	*types.Base
	// nolint: godox
	// TODO: This should be thread safe.
	DeviceGroups  map[string]int32
	DeviceGroups2 map[string]int32
	*types.WConfig
}

// APIVersion is a function that implements the Watcher interface.
func (w *Watcher) APIVersion() string {
	return constants.K8sAPIVersionV1
}

// Enabled is a function that check the resource can watch.
func (w *Watcher) Enabled() bool {
	return true
}

// Namespaced returns true if resource is namespaced
func (w *Watcher) Namespaced() bool {
	return true
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &corev1.Namespace{} // nolint: exhaustivestruct
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		namespace := obj.(*corev1.Namespace) // nolint: forcetypeassert
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": namespace.Name, "event": "add"}))
		logrus.Debugf("Handling add namespace event: %s", namespace.Name)

		for name, parentID := range w.DeviceGroups {
			resourceType, err2 := enums.ParseResourceType(name)
			if err2 != nil {
				continue
			}
			appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(resourceType.GetCategory()).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)

			// nolint: exhaustivestruct
			opts := &devicegroup.Options{
				AppliesTo:        appliesTo,
				Client:           w.LMClient,
				DisableAlerting:  w.Config.DisableAlerting,
				Name:             namespace.Name,
				ParentID:         parentID,
				CustomProperties: devicegroup.NewPropertyBuilder(),
			}

			logrus.Debugf("devicegroup create options: %v", opts)

			_, err := devicegroup.Create(lctx, opts, nil)
			if err != nil {
				logrus.Errorf("Failed to add namespace to %q: %v", name, err)
				return
			}

			logrus.Infof("Added namespace %q to %q", namespace.Name, name)
		}

		appliesTo := devicegroup.NewAppliesToBuilder().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)

		// nolint: exhaustivestruct
		opts := &devicegroup.Options{
			AppliesTo:        appliesTo,
			Client:           w.LMClient,
			DisableAlerting:  w.Config.DisableAlerting,
			Name:             namespace.Name,
			ParentID:         w.DeviceGroups2["Namespaces"],
			CustomProperties: devicegroup.NewPropertyBuilder(),
		}

		logrus.Debugf("Namespace create options: %v", opts)

		_, err := devicegroup.Create(lctx, opts, nil)
		if err != nil {
			logrus.Errorf("Failed to add namespace to %q: %v", namespace.Name, err)
			return
		}

		logrus.Infof("Added new structure namespace %q", namespace.Name)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		logrus.Debugf("Ignoring update namespace event")
		// oldNamespace := oldObj.(*v1.Namespace)
		// newNamespace := newObj.(*v1.Namespace)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		namespace := obj.(*corev1.Namespace) // nolint: forcetypeassert
		logrus.Debugf("Handle deleting namespace event: %s", namespace.Name)

		deviceGroups, err := devicegroup.FindDeviceGroupsByName(namespace.Name, w.LMClient)
		if err != nil {
			logrus.Errorf("Failed to get device group for namespace:\"%s\" with error: %v", namespace.Name, err)
			return
		}

		reversedDeviceGroups := getDeviceGroupsReversed(w.DeviceGroups)
		for _, d := range deviceGroups {
			if _, ok := reversedDeviceGroups[d.ParentID]; ok {
				err = devicegroup.DeleteGroup(d, w.LMClient)
				if err != nil {
					logrus.Errorf("Failed to delete device group of namespace:\"%s\" having ID:\"%d\" with error: %v", namespace.Name, d.ID, err)
				}
			}
		}
	}
}

func getDeviceGroupsReversed(deviceGroups map[string]int32) map[int32]string {
	reversedDeviceGroups := make(map[int32]string)
	for key, value := range deviceGroups {
		reversedDeviceGroups[value] = key
	}

	return reversedDeviceGroups
}
