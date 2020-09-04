// Package service provides the logic for mapping a Kubernetes Service to a
// LogicMonitor w.
package service

import (
	"fmt"
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
	resource = "services"
)

// Watcher represents a watcher type that watches services.
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
	return &v1.Service{}
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + service.Name}))
		log := lmlog.Logger(lctx)

		log.Infof("Service type is %s", service.Spec.Type)

		// Require an IP address.
		if service.Spec.ClusterIP == "" {
			log.Warningf("Service clusterIP is empty, service name : %s", service.Name)
			return
		}
		w.add(lctx, service)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*v1.Service)
		new := newObj.(*v1.Service)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + old.Name}))

		// If the old service does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new w.
		if old.Spec.ClusterIP == "" && new.Spec.ClusterIP != "" {
			w.add(lctx, new)
			return
		}

		// Covers the case when the old service is in the process of terminating
		// and the new service is coming up to replace it.
		// if old.Spec.ClusterIP != new.Spec.ClusterIP {
		w.update(lctx, old, new)
		// }
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + service.Name}))
		log := lmlog.Logger(lctx)
		// Delete the service.
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(lctx, w.Resource(), fmtServiceDisplayName(service)); err != nil {
				log.Errorf("Failed to delete service: %v", err)
				return
			}
			log.Infof("Deleted service %s", service.Name)
			return
		}

		// Move the service.
		w.move(lctx, service)
	}
}

// nolint: dupl
func (w *Watcher) add(lctx *lmctx.LMContext, service *v1.Service) {
	log := lmlog.Logger(lctx)
	if _, err := w.Add(lctx, w.Resource(), service.Labels,
		w.args(service, constants.ServiceCategory)...,
	); err != nil {
		log.Errorf("Failed to add service %q: %v", fmtServiceDisplayName(service), err)
		return
	}
	log.Infof("Added service %q", fmtServiceDisplayName(service))
}

func (w *Watcher) serviceUpdateFilter(old, new *v1.Service) types.UpdateFilter {
	return func() bool {
		return old.Spec.ClusterIP != new.Spec.ClusterIP
	}
}

func (w *Watcher) update(lctx *lmctx.LMContext, old, new *v1.Service) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceByDisplayName(lctx, "services",
		fmtServiceDisplayName(old), w.serviceUpdateFilter(old, new),
		w.args(new, constants.ServiceCategory)...,
	); err != nil {
		log.Errorf("Failed to update service %q: %v", fmtServiceDisplayName(new), err)
		return
	}
	log.Infof("Updated service %q", fmtServiceDisplayName(old))
}

// nolint: dupl
func (w *Watcher) move(lctx *lmctx.LMContext, service *v1.Service) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceFieldByDisplayName(lctx, w.Resource(), fmtServiceDisplayName(service), constants.CustomPropertiesFieldName, w.args(service, constants.ServiceDeletedCategory)...); err != nil {
		log.Errorf("Failed to move service %q: %v", fmtServiceDisplayName(service), err)
		return
	}
	log.Infof("Moved service %q", fmtServiceDisplayName(service))
}

func (w *Watcher) args(service *v1.Service, category string) []types.DeviceOption {
	return []types.DeviceOption{
		w.Name(service.Spec.ClusterIP),
		w.ResourceLabels(service.Labels),
		w.DisplayName(fmtServiceDisplayName(service)),
		w.SystemCategories(category),
		w.Auto("name", service.Name),
		w.Auto("namespace", service.Namespace),
		w.Auto("selflink", service.SelfLink),
		w.Auto("uid", string(service.UID)),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(service.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, fmtServiceDisplayName(service)),
	}
}

// FmtServiceDisplayName implements the conversion for the service display name
func fmtServiceDisplayName(service *v1.Service) string {
	return fmt.Sprintf("%s.%s.svc-%s", service.Name, service.Namespace, string(service.UID))
}

// GetServicesMap implements the getting services map info from k8s
func GetServicesMap(lctx *lmctx.LMContext, k8sClient kubernetes.Interface, namespace string) (map[string]string, error) {
	log := lmlog.Logger(lctx)
	servicesMap := make(map[string]string)
	serviceList, err := k8sClient.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil || serviceList == nil {
		log.Warnf("Failed to get the services from k8s")
		return nil, err
	}
	for i := range serviceList.Items {
		servicesMap[fmtServiceDisplayName(&serviceList.Items[i])] = serviceList.Items[i].Spec.ClusterIP
	}

	return servicesMap, nil
}
