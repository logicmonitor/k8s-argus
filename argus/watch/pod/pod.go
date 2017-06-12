package pod

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/logicmonitor/argus/argus/config"
	"github.com/logicmonitor/argus/constants"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "pods"
)

// Watcher represents a watcher type that watches pods.
type Watcher struct {
	LMClient  *lmv1.DefaultApi
	K8sClient *kubernetes.Clientset
	Config    *config.Config
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
			log.Infof("Adding pod %s", pod.Name)
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
			log.Infof("Adding pod %s", newPod.Name)
			w.addDevice(newPod)
		} else if oldPod.Status.PodIP != "" && newPod.Status.PodIP != "" {
			// Covers the case when a pod has been terminated (new ip doesn't exist)
			// and if a pod needs to be added.
			filter := fmt.Sprintf("displayName:%s", oldPod.Name)
			restResponse, _, err := w.LMClient.GetDeviceList("", -1, 0, filter)
			if err != nil {
				log.Errorf("Failed searching for pod %s: %s", oldPod.Name, restResponse.Errmsg)
			}
			if restResponse.Data.Total == 1 {
				id := restResponse.Data.Items[0].Id
				log.Infof("Updating pod %s with id %d", oldPod.Name, id)
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
		restResponse, _, err := w.LMClient.GetDeviceList("", -1, 0, filter)
		if err != nil {
			log.Errorf("Failed searching for pod %s: %s", pod.Name, restResponse.Errmsg)
		}
		if restResponse.Data.Total == 1 {
			id := restResponse.Data.Items[0].Id
			log.Infof("Deleting pod %s with id %d", pod.Name, id)
			w.LMClient.DeleteDevice(id)
		}
	}
}

func (w Watcher) addDevice(pod *v1.Pod) {
	device := w.makeDeviceObject(pod)
	restResponse, _, err := w.LMClient.AddDevice(device, false)
	if err != nil {
		log.Error(restResponse.Errmsg)
	}
}

func (w Watcher) updateDevice(pod *v1.Pod, id int32) {
	device := w.makeDeviceObject(pod)
	restResponse, _, err := w.LMClient.UpdateDevice(device, id, "")
	if err != nil {
		log.Error(restResponse.Errmsg)
	}
}

func (w Watcher) makeDeviceObject(pod *v1.Pod) (device lmv1.RestDevice) {
	categories := constants.PodCategory
	for k, v := range pod.Labels {
		categories += "," + k + "=" + v

	}

	device = lmv1.RestDevice{
		Name:                 pod.Name,
		DisplayName:          pod.Name,
		DisableAlerting:      w.Config.DisableAlerting,
		HostGroupIds:         "1",
		PreferredCollectorId: w.Config.PreferredCollector,
		CustomProperties: []lmv1.NameAndValue{
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
