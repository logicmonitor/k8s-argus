package pod

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/tree/device"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "pods"
)

// Watcher represents a watcher type that watches pods.
type Watcher struct {
	*types.Base
}

// Resource is a function that implements the Watcher interface.
func (w Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w Watcher) ObjType() runtime.Object {
	return &v1.Pod{}
}

// AddFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		pod := obj.(*v1.Pod)

		// We need an IP address.
		if pod.Status.PodIP == "" {
			return
		}

		// Check if the pod has already been added.
		d, err := device.FindByDisplayName(pod.Name, w.LMClient)
		if err != nil {
			log.Errorf("Failed to find pod %q: %v", pod.Name, err)
			return
		}

		// Add the pod.
		if d == nil {
			newDevice := w.makeDeviceObject(pod)
			err = device.Add(newDevice, w.LMClient)
			if err != nil {
				log.Errorf("Failed to add pod %q: %v", newDevice.DisplayName, err)
			}
			log.Infof("Added pod %s", newDevice.DisplayName)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		oldPod := oldObj.(*v1.Pod)
		newPod := newObj.(*v1.Pod)

		// If the old pod does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		if oldPod.Status.PodIP == "" && newPod.Status.PodIP != "" {
			d := w.makeDeviceObject(newPod)
			err := device.Add(d, w.LMClient)
			if err != nil {
				log.Errorf("Failed to add pod %s: %s", d.DisplayName, err)
				return
			}
			log.Infof("Added pod %s", d.DisplayName)
			return
		}

		// Covers the case when a pod has been terminated (new ip doesn't exist)
		// and if a pod needs to be added.
		if oldPod.Status.PodIP != "" && newPod.Status.PodIP != "" {
			oldDevice, err := device.FindByDisplayName(oldPod.Name, w.LMClient)
			if err != nil {
				log.Errorf("Failed to find pod %q: %v", oldPod.Name, err)
				return
			}

			// Update the pod.
			if oldDevice != nil {
				newDevice := w.makeDeviceObject(newPod)
				err := device.UpdateAndReplace(newDevice, oldDevice.Id, w.LMClient)
				if err != nil {
					log.Errorf("Failed to update pod %s: %v", oldDevice.DisplayName, err)
					return
				}
			}

			log.Infof("Updated pod %s with id %d", oldDevice.DisplayName, oldDevice.Id)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		pod := obj.(*v1.Pod)
		d, err := device.FindByDisplayName(pod.Name, w.LMClient)
		if err != nil {
			log.Errorf("Failed to find pod %q: %v", pod.Name, err)
			return
		}

		if d == nil {
			return
		}

		// Delete the device
		if w.Config.DeleteDevices {
			err = device.Delete(d, w.LMClient)
			if err != nil {
				log.Errorf("Failed to delete pod: %v", err)
				return
			}
			log.Infof("Deleted pod %s with id %d", d.DisplayName, d.Id)
			return
		}

		// Move the device

		categories := device.BuildSystemCategoriesFromLabels(constants.PodDeletedCategory, pod.Labels)
		newDevice := &lm.RestDevice{
			CustomProperties: []lm.NameAndValue{
				{
					Name:  "system.categories",
					Value: categories,
				},
			},
		}
		err = device.UpdateAndReplaceField(newDevice, d.Id, constants.CustomPropertiesFieldName, w.LMClient)
		if err != nil {
			log.Errorf("Failed to move pod %s: %s", d.DisplayName, err)
			return
		}

		log.Infof("Moved pod %s with id %d to deleted group", d.DisplayName, d.Id)
	}
}

func (w Watcher) makeDeviceObject(pod *v1.Pod) *lm.RestDevice {
	categories := device.BuildSystemCategoriesFromLabels(constants.PodCategory, pod.Labels)

	d := &lm.RestDevice{
		Name:                 pod.Name,
		DisplayName:          pod.Name,
		DisableAlerting:      w.Config.DisableAlerting,
		HostGroupIds:         "1",
		PreferredCollectorId: w.Config.PreferredCollector,
		CustomProperties: []lm.NameAndValue{
			{
				Name:  "system.categories",
				Value: categories,
			},
			{
				Name:  "auto.clustername",
				Value: w.Config.ClusterName,
			},
			{
				Name:  "auto.selflink",
				Value: pod.SelfLink,
			},
			{
				Name:  "auto.nodename",
				Value: pod.Spec.NodeName,
			},
			{
				Name:  "auto.name",
				Value: pod.Name,
			},
			{
				Name:  "auto.namespace",
				Value: pod.Namespace,
			},
			{
				Name:  "auto.uid",
				Value: string(pod.UID),
			},
			{
				Name:  "auto.ip",
				Value: pod.Status.PodIP,
			},
		},
	}

	return d
}
