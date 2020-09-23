// Package pod provides the logic for mapping a Kubernetes Pod to a
// LogicMonitor w.
package pod

import (
	"strconv"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
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
		pod := obj.(*v1.Pod)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + pod.Name}))
		log := lmlog.Logger(lctx)
		log.Debugf("Handling add pod event: %s", pod.Name)

		// Require an IP address.
		if pod.Status.PodIP == "" {
			return
		}
		w.add(lctx, pod)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*v1.Pod)
		new := newObj.(*v1.Pod)

		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + old.Name}))
		log := lmlog.Logger(lctx)
		log.Debugf("Handling update pod event: %s", old.Name)

		// If the old pod does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new w.
		if old.Status.PodIP == "" && new.Status.PodIP != "" {
			w.add(lctx, new)
			return
		}

		if new.Status.Phase == v1.PodSucceeded {
			if err := w.DeleteByDisplayName(lctx, w.Resource(), old.Name); err != nil {
				log.Errorf("Failed to delete pod: %v", err)
				return
			}
			log.Infof("Deleted pod %s", old.Name)
			return
		}

		// if old.Status.PodIP != new.Status.PodIP {
		w.update(lctx, old, new)
		// }
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		pod := obj.(*v1.Pod)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + pod.Name}))
		log := lmlog.Logger(lctx)

		log.Debugf("Handling delete pod event: %s", pod.Name)

		// Delete the pod.
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(lctx, w.Resource(), pod.Name); err != nil {
				log.Errorf("Failed to delete pod: %v", err)
				return
			}
			log.Infof("Deleted pod %s", pod.Name)
			return
		}

		// Move the pod.
		w.move(lctx, pod)
	}
}

// nolint: dupl

func (w *Watcher) add(lctx *lmctx.LMContext, pod *v1.Pod) {
	log := lmlog.Logger(lctx)

	p, err := w.Add(lctx, w.Resource(), pod.Labels,
		w.args(pod, constants.PodCategory)...,
	)
	if err != nil {
		log.Errorf("Failed to add pod %q: %v", pod.Name, err)
		return
	}

	if p == nil {
		log.Infof("pod %q is not added as it is mentioned for filtering.", pod.Name)
	}
	log.Infof("Added pod %q", pod.Name)
}

func (w *Watcher) podUpdateFilter(old, new *v1.Pod) types.UpdateFilter {
	return func() bool {
		return old.Status.PodIP != new.Status.PodIP
	}
}

func (w *Watcher) update(lctx *lmctx.LMContext, old, new *v1.Pod) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceByDisplayName(lctx, "pods",
		old.Name, w.podUpdateFilter(old, new), new.Labels,
		w.args(new, constants.PodCategory)...,
	); err != nil {
		log.Errorf("Failed to update pod %q: %v", new.Name, err)
		return
	}
	log.Infof("Updated pod %q", old.Name)
}

// nolint: dupl
func (w *Watcher) move(lctx *lmctx.LMContext, pod *v1.Pod) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceFieldByDisplayName(lctx, w.Resource(), pod.Name, constants.CustomPropertiesFieldName, w.args(pod, constants.PodDeletedCategory)...); err != nil {
		log.Errorf("Failed to move pod %q: %v", pod.Name, err)
		return
	}
	log.Infof("Moved pod %q", pod.Name)
}

func (w *Watcher) args(pod *v1.Pod, category string) []types.DeviceOption {
	options := []types.DeviceOption{
		w.Name(getPodDNSName(pod)),
		w.ResourceLabels(pod.Labels),
		w.DisplayName(pod.Name),
		w.SystemCategories(category),
		w.Auto("name", pod.Name),
		w.Auto("namespace", pod.Namespace),
		w.Auto("nodename", pod.Spec.NodeName),
		w.Auto("selflink", pod.SelfLink),
		w.Auto("uid", string(pod.UID)),
		w.System("ips", pod.Status.PodIP),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(pod.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, pod.Name),
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
func GetPodsMap(k8sClient kubernetes.Interface, namespace string) (map[string]string, error) {
	podsMap := make(map[string]string)
	podList, err := k8sClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil || podList == nil {
		return nil, err
	}
	for i := range podList.Items {
		// TODO: we should improve the value of the map to the ip of the pod when changing the name of the device to the ip
		podsMap[podList.Items[i].Name] = getPodDNSName(&podList.Items[i])
	}

	return podsMap, nil
}
