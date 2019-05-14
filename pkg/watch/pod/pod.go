// Package pod provides the logic for mapping a Kubernetes Pod to a
// LogicMonitor w.
package pod

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/err"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

const (
	resource = "pods"
)

// Watcher represents a watcher type that watches pods.
type Watcher struct {
	types.DeviceManager
}

// ApiVersion is a function that implements the Watcher interface.
func (w *Watcher) ApiVersion() string {
	return constants.K8sApiVersion_v1
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &v1.Pod{}
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		// Due to panic error in this call stack will crash the application; recovering those panics here could make our application robust.
		defer err.RecoverError("Add pod")
		pod := obj.(*v1.Pod)

		log.Debugf("Handling add pod event: %s", pod.Name)

		// Require an IP address.
		if pod.Status.PodIP == "" {
			return
		}
		w.add(pod)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		// Due to panic error in this call stack will crash the application; recovering those panics here could make our application robust.
		defer err.RecoverError("Update pod")
		old := oldObj.(*v1.Pod)
		new := newObj.(*v1.Pod)

		log.Debugf("Handling update pod event: %s", old.Name)

		// If the old pod does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new w.
		if old.Status.PodIP == "" && new.Status.PodIP != "" {
			w.add(new)
			return
		}

		if new.Status.Phase == v1.PodSucceeded {
			if err := w.DeleteByDisplayName(old.Name); err != nil {
				log.Errorf("Failed to delete pod: %v", err)
				return
			}
			log.Infof("Deleted pod %s", old.Name)
			return
		}

		if old.Status.PodIP != new.Status.PodIP {
			w.update(old, new)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		// Due to panic error in this call stack will crash the application; recovering those panics here could make our application robust.
		defer err.RecoverError("Delete pod")
		pod := obj.(*v1.Pod)

		log.Debugf("Handling delete pod event: %s", pod.Name)

		// Delete the pod.
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(pod.Name); err != nil {
				log.Errorf("Failed to delete pod: %v", err)
				return
			}
			log.Infof("Deleted pod %s", pod.Name)
			return
		}

		// Move the pod.
		w.move(pod)
	}
}

// nolint: dupl
func (w *Watcher) add(pod *v1.Pod) {
	if _, err := w.Add(
		w.args(pod, constants.PodCategory)...,
	); err != nil {
		log.Errorf("Failed to add pod %q: %v", pod.Name, err)
		return
	}
	log.Infof("Added pod %q", pod.Name)
}

func (w *Watcher) update(old, new *v1.Pod) {
	if _, err := w.UpdateAndReplaceByDisplayName(
		old.Name,
		w.args(new, constants.PodCategory)...,
	); err != nil {
		log.Errorf("Failed to update pod %q: %v", new.Name, err)
		return
	}
	log.Infof("Updated pod %q", old.Name)
}

// nolint: dupl
func (w *Watcher) move(pod *v1.Pod) {
	if _, err := w.UpdateAndReplaceFieldByDisplayName(pod.Name, constants.CustomPropertiesFieldName, w.args(pod, constants.PodDeletedCategory)...); err != nil {
		log.Errorf("Failed to move pod %q: %v", pod.Name, err)
		return
	}
	log.Infof("Moved pod %q", pod.Name)
}

func (w *Watcher) args(pod *v1.Pod, category string) []types.DeviceOption {
	categories := utilities.BuildSystemCategoriesFromLabels(category, pod.Labels)
	options := []types.DeviceOption{
		w.Name(getPodDNSName(pod)),
		w.ResourceLabels(pod.Labels),
		w.DisplayName(pod.Name),
		w.SystemCategories(categories),
		w.Auto("name", pod.Name),
		w.Auto("namespace", pod.Namespace),
		w.Auto("nodename", pod.Spec.NodeName),
		w.Auto("selflink", pod.SelfLink),
		w.Auto("uid", string(pod.UID)),
		w.System("ips", pod.Status.PodIP),
	}
	if pod.Spec.HostNetwork {
		options = append(options, w.Custom("kubernetes.pod.hostNetwork", "true"))
	}
	return options
}

func getPodDNSName(pod *v1.Pod) string {
	// if the pod is configured as "hostnetwork=true", we will use the pod name as the IP/DNS name of the pod device
	if pod.Spec.HostNetwork {
		return pod.Name
	}
	return pod.Status.PodIP
}

// GetPodsMap implements the getting pods map info from k8s
func GetPodsMap(k8sClient *kubernetes.Clientset, namespace string) (map[string]string, error) {
	podsMap := make(map[string]string)
	podList, err := k8sClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil || podList == nil {
		return nil, err
	}
	for _, podInfo := range podList.Items {
		// TODO: we should improve the value of the map to the ip of the pod when changing the name of the device to the ip
		podsMap[podInfo.Name] = getPodDNSName(&podInfo)
	}

	return podsMap, nil
}
