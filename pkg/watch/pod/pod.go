// Package pod provides the logic for mapping a Kubernetes Pod to a
// LogicMonitor w.
package pod

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "pods"
)

// Watcher represents a watcher type that watches pods.
type Watcher struct {
	types.DeviceManager
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
		old := oldObj.(*v1.Pod)
		new := newObj.(*v1.Pod)

		// If the old pod does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new w.
		if old.Status.PodIP == "" && new.Status.PodIP != "" {
			w.add(new)
		}

		// Covers the case when the old pod is in the process of terminating
		// and the new pod is coming up to replace it.
		if old.Status.PodIP != "" && new.Status.PodIP != "" {
			w.update(old, new)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		pod := obj.(*v1.Pod)

		// Delete the pod.
		if w.Config().DeleteDevices {
			if err := w.DeleteByName(pod.Name); err != nil {
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
	if _, err := w.UpdateAndReplaceByName(
		old.Name,
		w.args(new, constants.PodCategory)...,
	); err != nil {
		log.Errorf("Failed to update pod %q: %v", new.Name, err)
		return
	}
	log.Infof("Updated pod %q", old.Name)
}

func (w *Watcher) move(pod *v1.Pod) {
	if _, err := w.UpdateAndReplaceFieldByName(pod.Name, constants.CustomPropertiesFieldName, w.args(pod, constants.PodDeletedCategory)...); err != nil {
		log.Errorf("Failed to move pod %q: %v", pod.Name, err)
		return
	}
	log.Infof("Moved pod %q", pod.Name)
}

func (w *Watcher) args(pod *v1.Pod, category string) []types.DeviceOption {
	categories := utilities.BuildSystemCategoriesFromLabels(category, pod.Labels)
	return []types.DeviceOption{
		w.Name(pod.Name),
		w.DisplayName(pod.Name),
		w.SystemCategories(categories),
		w.Auto("name", pod.Name),
		w.Auto("namespace", pod.Namespace),
		w.Auto("nodename", pod.Spec.NodeName),
		w.Auto("selflink", pod.SelfLink),
		w.Auto("uid", string(pod.UID)),
		w.System("ips", pod.Status.PodIP),
	}
}
