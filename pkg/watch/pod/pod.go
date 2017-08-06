package pod

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
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
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		pod := obj.(*v1.Pod)

		if pod.Status.PodIP != "" {
			w.addDevice(pod)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		oldPod := oldObj.(*v1.Pod)
		newPod := newObj.(*v1.Pod)

		// If the old pod does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		if oldPod.Status.PodIP == "" && newPod.Status.PodIP != "" {
			w.addDevice(newPod)
		} else if oldPod.Status.PodIP != "" && newPod.Status.PodIP != "" {
			// Covers the case when a pod has been terminated (new ip doesn't exist)
			// and if a pod needs to be added.
			filter := fmt.Sprintf("displayName:%s", oldPod.Name)
			restResponse, apiResponse, err := w.LMClient.GetDeviceList("", -1, 0, filter)
			if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
				log.Errorf("Failed to find pod %s: %s", oldPod.Name, _err)

				return
			}
			if restResponse.Data.Total == 1 {
				id := restResponse.Data.Items[0].Id
				w.updateDevice(newPod, id)
			}
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		pod := obj.(*v1.Pod)
		filter := fmt.Sprintf("displayName:%s", pod.Name)
		restResponse, apiResponse, err := w.LMClient.GetDeviceList("", -1, 0, filter)
		if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
			log.Errorf("Failed to find pod %s: %s", pod.Name, _err)

			return
		}
		if restResponse.Data.Total == 1 {
			id := restResponse.Data.Items[0].Id
			if w.Config.DeleteDevices {
				restResponse, apiResponse, err := w.LMClient.DeleteDevice(id)
				if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
					log.Errorf("Failed to delete device with id %q: %s", id, _err)

					return
				}
				log.Infof("Deleted pod %s with id %d", pod.Name, id)
			} else {
				categories := constants.PodDeletedCategory
				for k, v := range pod.Labels {
					categories += "," + k + "=" + v

				}
				device := lm.RestDevice{
					CustomProperties: []lm.NameAndValue{
						{
							Name:  "system.categories",
							Value: categories,
						},
					},
				}
				restResponse, apiResponse, err := w.LMClient.PatchDeviceById(device, id, "replace", "customProperties")
				if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
					log.Errorf("Failed to move pod %s: %s", pod.Name, _err)

					return
				}
				log.Infof("Moved pod %s with id %d to deleted group", pod.Name, id)
			}
		}
	}
}

func (w Watcher) addDevice(pod *v1.Pod) {
	device := w.makeDeviceObject(pod)
	restResponse, apiResponse, err := w.LMClient.AddDevice(device, false)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to add pod %s: %s", pod.Name, _err)

		return
	}
	log.Infof("Added pod %s", pod.Name)
}

func (w Watcher) updateDevice(pod *v1.Pod, id int32) {
	device := w.makeDeviceObject(pod)
	restResponse, apiResponse, err := w.LMClient.UpdateDevice(device, id, "replace")
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to update pod %s: %s", pod.Name, _err)

		return
	}
	log.Infof("Updated pod %s with id %d", pod.Name, id)
}

func (w Watcher) makeDeviceObject(pod *v1.Pod) (device lm.RestDevice) {
	categories := constants.PodCategory
	for k, v := range pod.Labels {
		categories += "," + k + "=" + v

	}

	device = lm.RestDevice{
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

	return
}
