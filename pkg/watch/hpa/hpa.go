// Package hpa provides the logic for mapping a Kubernetes horizontalPodAutoscaler to a
// LogicMonitor w.
// nolint: dupl
package hpa

import (
	"fmt"
	"strconv"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	log "github.com/sirupsen/logrus"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

const (
	resource = "horizontalpodautoscalers"
)

// Watcher represents a watcher type that watches horizontalPodAutoscalers.
type Watcher struct {
	types.DeviceManager
}

// APIVersion is a function that implements the Watcher interface.
func (w *Watcher) APIVersion() string {
	return constants.K8sAutoscalingV1
}

// Enabled is a function that check the resource can watch.
func (w *Watcher) Enabled() bool {
	return permission.HasHorizontalPodAutoscalerPermissions()
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &autoscalingv1.HorizontalPodAutoscaler{}

}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		horizontalPodAutoscaler := obj.(*autoscalingv1.HorizontalPodAutoscaler)
		log.Infof("Handling add horizontalPodAutoscaler event: %s", horizontalPodAutoscaler.Name)
		w.add(horizontalPodAutoscaler)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*autoscalingv1.HorizontalPodAutoscaler)
		new := newObj.(*autoscalingv1.HorizontalPodAutoscaler)
		w.update(old, new)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		horizontalPodAutoscaler := obj.(*autoscalingv1.HorizontalPodAutoscaler)
		log.Debugf("Handling delete horizontalPodAutoscaler event: %s", horizontalPodAutoscaler.Name)
		// Delete the horizontalPodAutoscaler.
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler)); err != nil {
				log.Errorf("Failed to delete horizontalPodAutoscaler: %v", err)
				return
			}
			log.Infof("Deleted horizontalPodAutoscaler %s", horizontalPodAutoscaler.Name)
			return
		}

		// Move the horizontalPodAutoscaler.
		w.move(horizontalPodAutoscaler)
	}
}

// nolint: dupl
func (w *Watcher) add(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) {
	if _, err := w.Add(
		w.args(horizontalPodAutoscaler, constants.HorizontalPodAutoscalerCategory)...,
	); err != nil {
		log.Errorf("Failed to add horizontalPodAutoscaler %q: %v", fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler), err)
		return
	}
	log.Infof("Added horizontalPodAutoscaler %q", fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler))
}

func (w *Watcher) update(old, new *autoscalingv1.HorizontalPodAutoscaler) {
	if _, err := w.UpdateAndReplaceByDisplayName(
		fmtHorizontalPodAutoscalerDisplayName(old),
		w.args(new, constants.HorizontalPodAutoscalerCategory)...,
	); err != nil {
		log.Errorf("Failed to update horizontalPodAutoscaler %q: %v", fmtHorizontalPodAutoscalerDisplayName(new), err)
		return
	}
	log.Infof("Updated horizontalPodAutoscaler %q", fmtHorizontalPodAutoscalerDisplayName(old))
}

func (w *Watcher) move(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) {
	if _, err := w.UpdateAndReplaceFieldByDisplayName(fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler), constants.CustomPropertiesFieldName, w.args(horizontalPodAutoscaler, constants.HorizontalPodAutoscalerDeletedCategory)...); err != nil {
		log.Errorf("Failed to move horizontalPodAutoscaler %q: %v", fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler), err)
		return
	}
	log.Infof("Moved horizontalPodAutoscaler %q", fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler))
}

func (w *Watcher) args(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, category string) []types.DeviceOption {
	return []types.DeviceOption{
		w.Name(horizontalPodAutoscaler.Name),
		w.ResourceLabels(horizontalPodAutoscaler.Labels),
		w.DisplayName(fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler)),
		w.SystemCategories(category),
		w.Auto("name", horizontalPodAutoscaler.Name),
		w.Auto("namespace", horizontalPodAutoscaler.Namespace),
		w.Auto("selflink", horizontalPodAutoscaler.SelfLink),
		w.Auto("uid", string(horizontalPodAutoscaler.UID)),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(horizontalPodAutoscaler.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler)),
	}
}

// FmthorizontalPodAutoscalerDisplayName implements the conversion for the horizontalPodAutoscaler display name
func fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) string {
	return fmt.Sprintf("%s.%s.hpa-%s", horizontalPodAutoscaler.Name, horizontalPodAutoscaler.Namespace, string(horizontalPodAutoscaler.UID))
}

// GetHorizontalPodAutoscalersMap implements the getting horizontalPodAutoscaler map info from k8s
func GetHorizontalPodAutoscalersMap(k8sClient *kubernetes.Clientset, namespace string) (map[string]string, error) {

	horizontalPodAutoscalersMap := make(map[string]string)
	horizontalPodAutoscalerV1List, err := k8sClient.AutoscalingV1().HorizontalPodAutoscalers(namespace).List(v1.ListOptions{})
	if err != nil || horizontalPodAutoscalerV1List == nil {
		log.Warnf("Failed to get the horizontalPodAutoscalers from k8s")
		return nil, err
	}

	for _, horizontalPodAutoscalerInfo := range horizontalPodAutoscalerV1List.Items {
		horizontalPodAutoscalersMap[fmtHorizontalPodAutoscalerDisplayName(&horizontalPodAutoscalerInfo)] = horizontalPodAutoscalerInfo.Name
	}

	return horizontalPodAutoscalersMap, nil
}
