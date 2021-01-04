// Package hpa provides the logic for mapping a Kubernetes horizontalPodAutoscaler to a
// LogicMonitor w.
// nolint: dupl
package hpa

import (
	"fmt"
	"strconv"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
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
	*types.WConfig
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
		lctx := lmlog.NewLMContextWith(log.WithFields(log.Fields{"device_id": resource + "-" + horizontalPodAutoscaler.Name}))
		lctx = util.WatcherContext(lctx, w)
		log := lmlog.Logger(lctx)
		log.Infof("Handling add horizontalPodAutoscaler event: %s", horizontalPodAutoscaler.Name)
		w.add(lctx, horizontalPodAutoscaler)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*autoscalingv1.HorizontalPodAutoscaler)
		lctx := lmlog.NewLMContextWith(log.WithFields(log.Fields{"device_id": resource + "-" + old.Name}))
		lctx = util.WatcherContext(lctx, w)
		log := lmlog.Logger(lctx)
		log.Debugf("Handling update horizontalPodAutoscaler event: %s", old.Name)
		new := newObj.(*autoscalingv1.HorizontalPodAutoscaler)
		w.update(lctx, old, new)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		horizontalPodAutoscaler := obj.(*autoscalingv1.HorizontalPodAutoscaler)
		lctx := lmlog.NewLMContextWith(log.WithFields(log.Fields{"device_id": resource + "-" + horizontalPodAutoscaler.Name}))
		log := lmlog.Logger(lctx)
		log.Debugf("Handling delete horizontalPodAutoscaler event: %s", horizontalPodAutoscaler.Name)

		// Delete the horizontalPodAutoscaler.
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(horizontalPodAutoscaler),
				fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler, w.Config().ClusterName)); err != nil {
				log.Errorf("Failed to delete horizontalPodAutoscaler: %v", err)
				return
			}
			log.Infof("Deleted horizontalPodAutoscaler %s", horizontalPodAutoscaler.Name)
			return
		}

		// Move the horizontalPodAutoscaler.
		w.move(lctx, horizontalPodAutoscaler)
	}
}

// nolint: dupl
func (w *Watcher) add(lctx *lmctx.LMContext, horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) {
	log := lmlog.Logger(lctx)
	hpa, err := w.Add(lctx, w.Resource(), horizontalPodAutoscaler.Labels,
		w.args(horizontalPodAutoscaler, constants.HorizontalPodAutoscalerCategory)...,
	)
	if err != nil {
		log.Errorf("Failed to add horizontalPodAutoscaler %q: %v", w.getDesiredDisplayName(horizontalPodAutoscaler), err)
		return
	}

	if hpa == nil {
		log.Debugf("horizontalPodAutoscaler %q is not added as it is mentioned for filtering.", w.getDesiredDisplayName(horizontalPodAutoscaler))
		return
	}

	log.Infof("Added horizontalPodAutoscaler %q", *hpa.DisplayName)
}

func (w *Watcher) update(lctx *lmctx.LMContext, old, new *autoscalingv1.HorizontalPodAutoscaler) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(old),
		fmtHorizontalPodAutoscalerDisplayName(old, w.Config().ClusterName), nil, new.Labels,
		w.args(new, constants.HorizontalPodAutoscalerCategory)...,
	); err != nil {
		log.Errorf("Failed to update horizontalPodAutoscaler %q: %v", w.getDesiredDisplayName(new), err)
		return
	}
	log.Infof("Updated horizontalPodAutoscaler %q", w.getDesiredDisplayName(old))
}

func (w *Watcher) move(lctx *lmctx.LMContext, horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceFieldByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(horizontalPodAutoscaler),
		fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler, w.Config().ClusterName), constants.CustomPropertiesFieldName,
		w.args(horizontalPodAutoscaler, constants.HorizontalPodAutoscalerDeletedCategory)...); err != nil {
		log.Errorf("Failed to move horizontalPodAutoscaler %q: %v", w.getDesiredDisplayName(horizontalPodAutoscaler), err)
		return
	}
	log.Infof("Moved horizontalPodAutoscaler %q", w.getDesiredDisplayName(horizontalPodAutoscaler))
}

func (w *Watcher) args(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, category string) []types.DeviceOption {
	return []types.DeviceOption{
		w.Name(w.getDesiredDisplayName(horizontalPodAutoscaler)),
		w.ResourceLabels(horizontalPodAutoscaler.Labels),
		w.DisplayName(w.getDesiredDisplayName(horizontalPodAutoscaler)),
		w.SystemCategories(category),
		w.Auto("name", horizontalPodAutoscaler.Name),
		w.Auto("namespace", horizontalPodAutoscaler.Namespace),
		w.Auto("selflink", horizontalPodAutoscaler.SelfLink),
		w.Auto("uid", string(horizontalPodAutoscaler.UID)),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(horizontalPodAutoscaler.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, w.getDesiredDisplayName(horizontalPodAutoscaler)),
	}
}

// FmthorizontalPodAutoscalerDisplayName implements the conversion for the horizontalPodAutoscaler display name
func fmtHorizontalPodAutoscalerDisplayName(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler, clusterName string) string {
	return fmt.Sprintf("%s-hpa-%s-%s", horizontalPodAutoscaler.Name, horizontalPodAutoscaler.Namespace, clusterName)
}

func (w *Watcher) getDesiredDisplayName(horizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler) string {
	return w.DeviceManager.GetDesiredDisplayName(horizontalPodAutoscaler.Name, horizontalPodAutoscaler.Namespace, constants.HorizontalPodAutoScalers)
}

// GetHorizontalPodAutoscalersMap implements the getting horizontalPodAutoscaler map info from k8s
func GetHorizontalPodAutoscalersMap(lctx *lmctx.LMContext, k8sClient *kubernetes.Clientset, namespace string, clusterName string) (map[string]string, error) {
	log := lmlog.Logger(lctx)
	horizontalPodAutoscalersMap := make(map[string]string)
	horizontalPodAutoscalerV1List, err := k8sClient.AutoscalingV1().HorizontalPodAutoscalers(namespace).List(v1.ListOptions{})
	if err != nil || horizontalPodAutoscalerV1List == nil {
		log.Warnf("Failed to get the horizontalPodAutoscalers from k8s")
		return nil, err
	}

	for _, horizontalPodAutoscalerInfo := range horizontalPodAutoscalerV1List.Items {
		horizontalPodAutoscalersMap[fmtHorizontalPodAutoscalerDisplayName(&horizontalPodAutoscalerInfo, clusterName)] = horizontalPodAutoscalerInfo.Name
	}

	return horizontalPodAutoscalersMap, nil
}
