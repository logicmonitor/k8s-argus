package namespace

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// OldWatcher represents a watcher type that watches namespaces.
type OldWatcher struct {
	Resource enums.ResourceType
	*types.Base
	// nolint: godox
	// TODO: This should be thread safe.
	DeviceGroups  map[string]int32
	DeviceGroups2 map[string]int32
	ResourceCache *devicecache.ResourceCache
}

// ResourceType resource
func (w *OldWatcher) ResourceType() enums.ResourceType {
	return w.Resource
}

// GetConfig get
func (w *OldWatcher) GetConfig() *types.WConfig {
	return nil
}

// AddFunc is a function that implements the Watcher interface.
func (w *OldWatcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		namespace := obj.(*corev1.Namespace) // nolint: forcetypeassert
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": namespace.Name, "event": "add"}))
		logrus.Debugf("Handling add namespace event: %s", namespace.Name)
		// resource wise separate static groups and underneath namespace groups in each
		w.createPreviousDeviceTree(namespace, lctx)

		// this will be based on new device tree where namespace groups will be created and all resources to put under it
		// w.createNewDeviceTree(namespace, lctx)
	}
}

// nolint: unused
func (w *OldWatcher) createNewDeviceTree(namespace *corev1.Namespace, lctx *lmctx.LMContext) {
	appliesTo := devicegroup.NewAppliesToBuilder().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)

	// nolint: exhaustivestruct
	opts := &devicegroup.Options{
		AppliesTo:       appliesTo,
		Client:          w.LMClient,
		DisableAlerting: w.Config.DisableAlerting,
		Name:            namespace.Name,
		ParentID:        w.DeviceGroups2["Namespaces"],
	}

	logrus.Debugf("Namespace create options: %v", opts)

	_, err := devicegroup.Create(lctx, opts, w.ResourceCache)
	if err != nil {
		logrus.Errorf("Failed to add namespace to %q: %v", namespace.Name, err)
		return
	}

	logrus.Infof("Added new structure namespace %q", namespace.Name)
}

func (w *OldWatcher) createPreviousDeviceTree(namespace *corev1.Namespace, lctx *lmctx.LMContext) {
	for deviceGroupName, parentID := range w.DeviceGroups {
		var appliesTo devicegroup.AppliesToBuilder
		// Ensure that we are creating namespaces for namespaced resources.
		switch deviceGroupName {
		case constants.ServiceDeviceGroupName:
			appliesTo = devicegroup.NewAppliesToBuilder().HasCategory(constants.ServiceCategory).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)
		case constants.PodDeviceGroupName:
			appliesTo = devicegroup.NewAppliesToBuilder().HasCategory(constants.PodCategory).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)
		case constants.DeploymentDeviceGroupName:
			appliesTo = devicegroup.NewAppliesToBuilder().HasCategory(constants.DeploymentCategory).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)
		case constants.HorizontalPodAutoscalerDeviceGroupName:
			appliesTo = devicegroup.NewAppliesToBuilder().HasCategory(constants.HorizontalPodAutoscalerCategory).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(w.Config.ClusterName)
		default:

			continue
		}

		// nolint: exhaustivestruct
		opts := &devicegroup.Options{
			AppliesTo:        appliesTo,
			Client:           w.LMClient,
			DisableAlerting:  w.Config.DisableAlerting,
			Name:             namespace.Name,
			ParentID:         parentID,
			CustomProperties: devicegroup.NewPropertyBuilder(),
		}

		logrus.Debugf("namespace group options: %v", opts)
		key := cache.ResourceName{
			Name:     fmt.Sprintf("ns-%s", namespace.Name),
			Resource: enums.Namespaces,
		}

		_, ok := w.ResourceCache.Exists(lctx, key, deviceGroupName)
		if !ok {
			dgID, err := devicegroup.Create(lctx, opts, w.ResourceCache)
			if err != nil {
				logrus.Errorf("Failed to add %q namespace group under %q device group. Error: %v", namespace.Name, deviceGroupName, err)
				// continue to add remaining groups
				continue
			}
			w.ResourceCache.Set(key, cache.ResourceMeta{
				Container: deviceGroupName,
				LMID:      dgID,
			})
			logrus.Infof("Added namespace %q to %q", namespace.Name, deviceGroupName)
		} else {
			logrus.Debugf("Found device group in cache %s in %s, ignoring add event", namespace.Name, deviceGroupName)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *OldWatcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		logrus.Debugf("Ignoring update namespace event")
		// oldNamespace := oldObj.(*v1.Namespace)
		// newNamespace := newObj.(*v1.Namespace)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *OldWatcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		namespace := obj.(*corev1.Namespace) // nolint: forcetypeassert
		logrus.Debugf("Handle deleting namespace event: %s", namespace.Name)

		deviceGroups, err := devicegroup.FindDeviceGroupsByName(namespace.Name, w.LMClient)
		if err != nil {
			logrus.Errorf("Failed to get device group for namespace:\"%s\" with error: %v", namespace.Name, err)
			return
		}

		reversedDeviceGroups := getReversedDeviceGroups(w.DeviceGroups)
		for _, d := range deviceGroups {
			if parentDG, ok := reversedDeviceGroups[d.ParentID]; ok {
				err = devicegroup.DeleteGroup(d, w.LMClient)
				if err != nil {
					logrus.Errorf("Failed to delete device group of namespace:\"%s\" having ID:\"%d\" with error: %v", namespace.Name, d.ID, err)
				}
				key := cache.ResourceName{
					Name:     fmt.Sprintf("ns-%s", namespace.Name),
					Resource: enums.Namespaces,
				}
				w.ResourceCache.Unset(key, parentDG)
			}
		}
	}
}

func getReversedDeviceGroups(deviceGroups map[string]int32) map[int32]string {
	reversedDeviceGroups := make(map[int32]string)
	for key, value := range deviceGroups {
		reversedDeviceGroups[value] = key
	}

	return reversedDeviceGroups
}

// GetNamespaceList Fetches list of namespaces name
func GetNamespaceList(lctx *lmctx.LMContext, kubeClient kubernetes.Interface) []string {
	log := lmlog.Logger(lctx)
	namespaceList := make([]string, 0)
	namespaces, err := kubeClient.CoreV1().Namespaces().List(constants.DefaultListOptions)
	if err != nil || namespaces == nil {
		log.Warnf("Failed to get namespaces from k8s. Error: %v", err)

		return namespaceList
	}
	for i := range namespaces.Items {
		namespaceList = append(namespaceList, namespaces.Items[i].GetName())
	}

	return namespaceList
}
