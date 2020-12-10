// Package pod provides the logic for mapping a Kubernetes Pod to a
// LogicMonitor w.
package pod

import (
	"fmt"
	"strconv"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
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
		lctx = util.WatcherContext(lctx, w)
		log := lmlog.Logger(lctx)
		log.Debugf("Handling add pod event: %s", w.getDesiredDisplayName(pod))

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
		lctx = util.WatcherContext(lctx, w)
		log := lmlog.Logger(lctx)
		log.Debugf("Handling update pod event: %s", w.getDesiredDisplayName(new))

		// If the old pod does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new w.
		if old.Status.PodIP == "" && new.Status.PodIP != "" {
			w.add(lctx, new)
			return
		}

		// nolint: dupl
		if new.Status.Phase == v1.PodSucceeded {
			if err := w.DeleteByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(old),
				fmtPodDisplayName(old, w.Config().ClusterName)); err != nil {
				log.Errorf("Failed to delete pod: %v", err)
				return
			}
			log.Infof("Deleted pod %s", w.getDesiredDisplayName(old))
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

		log.Debugf("Handling delete pod event: %s", w.getDesiredDisplayName(pod))

		// Delete the pod.
		// nolint: dupl
		if w.Config().DeleteDevices {
			err := w.DeleteByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(pod),
				fmtPodDisplayName(pod, w.Config().ClusterName))
			if err != nil {
				log.Errorf("Failed to delete pod: %v", err)
				return
			}
			log.Infof("Deleted pod %s", w.getDesiredDisplayName(pod))
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
		log.Errorf("Failed to add pod %q: %v", w.getDesiredDisplayName(pod), err)
		return
	}

	if p == nil {
		log.Debugf("pod %q is not added as it is mentioned for filtering.", w.getDesiredDisplayName(pod))
		return
	}
	log.Infof("Added pod %q", *p.DisplayName)
}

func (w *Watcher) podUpdateFilter(old, new *v1.Pod) types.UpdateFilter {
	return func() bool {
		return old.Status.PodIP != new.Status.PodIP
	}
}

// nolint: dupl
func (w *Watcher) update(lctx *lmctx.LMContext, old, new *v1.Pod) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(old),
		fmtPodDisplayName(old, w.Config().ClusterName), w.podUpdateFilter(old, new), new.Labels,
		w.args(new, constants.PodCategory)...,
	); err != nil {
		log.Errorf("Failed to update pod %q: %v", w.getDesiredDisplayName(new), err)
		return
	}
	log.Infof("Updated pod %q", w.getDesiredDisplayName(old))
}

// nolint: dupl
func (w *Watcher) move(lctx *lmctx.LMContext, pod *v1.Pod) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceFieldByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(pod),
		fmtPodDisplayName(pod, w.Config().ClusterName), constants.CustomPropertiesFieldName,
		w.args(pod, constants.PodDeletedCategory)...); err != nil {
		log.Errorf("Failed to move pod %q: %v", w.getDesiredDisplayName(pod), err)
		return
	}
	log.Infof("Moved pod %q", w.getDesiredDisplayName(pod))
}

func (w *Watcher) args(pod *v1.Pod, category string) []types.DeviceOption {
	options := []types.DeviceOption{
		w.Name(getPodDNSName(pod)),
		w.ResourceLabels(pod.Labels),
		w.DisplayName(w.getDesiredDisplayName(pod)),
		w.SystemCategories(category),
		w.Auto("name", pod.Name),
		w.Auto("namespace", pod.Namespace),
		w.Auto("nodename", pod.Spec.NodeName),
		w.Auto("selflink", pod.SelfLink),
		w.Auto("uid", string(pod.UID)),
		w.System("ips", pod.Status.PodIP),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(pod.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, w.getDesiredDisplayName(pod)),
	}
	// Pod running on fargate doesn't support HostNetwork so check fargate profile label, if label exists then mark hostNetwork as true
	if pod.Spec.HostNetwork || pod.Labels[constants.LabelFargateProfile] != "" {
		options = append(options, w.Custom("kubernetes.pod.hostNetwork", "true"))
	}
	return options
}

func (w *Watcher) getDesiredDisplayName(pod *v1.Pod) string {
	return w.DeviceManager.GetDesiredDisplayName(pod.Name, pod.Namespace, constants.Pods)
}

// fmtPodDisplayName implements the conversion for the pod display name
func fmtPodDisplayName(pod *v1.Pod, clusterName string) string {
	return fmt.Sprintf("%s-pod-%s-%s", pod.Name, pod.Namespace, clusterName)
}

func getPodDNSName(pod *v1.Pod) string {
	// if the pod is configured as "hostnetwork=true" or running on fargate, we will use the pod name as the IP/DNS name of the pod device
	if pod.Spec.HostNetwork || pod.Labels[constants.LabelFargateProfile] != "" {
		return pod.Name
	}
	return pod.Status.PodIP
}

// GetPodsMap implements the getting pods map info from k8s
func GetPodsMap(k8sClient kubernetes.Interface, namespace string, clusterName string) (map[string]string, error) {
	podsMap := make(map[string]string)
	podList, err := k8sClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil || podList == nil {
		return nil, err
	}
	for i := range podList.Items {
		// TODO: we should improve the value of the map to the ip of the pod when changing the name of the device to the ip
		podsMap[fmtPodDisplayName(&podList.Items[i], clusterName)] = getPodDNSName(&podList.Items[i])
	}

	return podsMap, nil
}
