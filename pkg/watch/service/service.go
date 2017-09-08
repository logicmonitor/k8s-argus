// Package service provides the logic for mapping a Kubernetes Service to a
// LogicMonitor w.
package service

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "services"
)

// Watcher represents a watcher type that watches services.
type Watcher struct {
	types.DeviceManager
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
		// Only add the service if it is has a ClusterIP.
		if service.Spec.Type != v1.ServiceTypeClusterIP {
			return
		}
		// Require an IP address.
		if service.Spec.ClusterIP == "" {
			return
		}
		w.add(service)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*v1.Service)
		new := newObj.(*v1.Service)

		// Only add the service if it is has a ClusterIP.
		if new.Spec.Type != v1.ServiceTypeClusterIP {
			return
		}

		// If the old service does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new w.
		if old.Spec.ClusterIP == "" && new.Spec.ClusterIP != "" {
			w.add(new)
		}

		// Covers the case when the old service is in the process of terminating
		// and the new service is coming up to replace it.
		if old.Spec.ClusterIP != "" && new.Spec.ClusterIP != "" {
			w.update(old, new)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)

		// Delete the service.
		if w.Config().DeleteDevices {
			if err := w.DeleteByName(fmtServiceName(service)); err != nil {
				log.Errorf("Failed to delete service: %v", err)
				return
			}
			log.Infof("Deleted service %s", service.Name)
			return
		}

		// Move the service.
		w.move(service)
	}
}

// nolint: dupl
func (w *Watcher) add(service *v1.Service) {
	if _, err := w.Add(
		w.args(service, constants.ServiceCategory)...,
	); err != nil {
		log.Errorf("Failed to add service %q: %v", service.Name, err)
		return
	}
	log.Infof("Added service %q", service.Name)
}

func (w *Watcher) update(old, new *v1.Service) {
	if _, err := w.UpdateAndReplaceByName(
		old.Name,
		w.args(new, constants.ServiceCategory)...,
	); err != nil {
		log.Errorf("Failed to update service %q: %v", new.Name, err)
		return
	}
	log.Infof("Updated service %q", old.Name)
}

func (w *Watcher) move(service *v1.Service) {
	if _, err := w.UpdateAndReplaceFieldByName(
		service.Name,
		constants.CustomPropertiesFieldName,
		w.args(service, constants.ServiceDeletedCategory)...,
	); err != nil {
		log.Errorf("Failed to move service %q: %v", service.Name, err)
		return
	}
	log.Infof("Moved service %q", service.Name)
}

func (w *Watcher) args(service *v1.Service, category string) []types.DeviceOption {
	categories := utilities.BuildSystemCategoriesFromLabels(category, service.Labels)
	return []types.DeviceOption{
		w.Name(fmtServiceName(service)),
		w.DisplayName(fmtServiceDisplayName(service)),
		w.SystemCategories(categories),
		w.Auto("name", service.Name),
		w.Auto("namespace", service.Namespace),
		w.Auto("selflink", service.SelfLink),
		w.Auto("uid", string(service.UID)),
	}
}

func fmtServiceName(service *v1.Service) string {
	return service.Name + "." + service.Namespace + ".svc"
}

func fmtServiceDisplayName(service *v1.Service) string {
	return fmtServiceName(service) + "-" + string(service.UID)
}
