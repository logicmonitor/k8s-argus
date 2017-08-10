package service

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/tree/device"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "services"
)

// Watcher represents a watcher type that watches services.
type Watcher struct {
	*types.Base
}

// Resource is a function that implements the Watcher interface.
func (w Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w Watcher) ObjType() runtime.Object {
	return &v1.Service{}
}

// AddFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)

		// Only add the service if it is of the type 'ClusterIP'.
		if service.Spec.Type != v1.ServiceTypeClusterIP {
			return
		}

		// We need an IP address.
		if service.Spec.ClusterIP == "" {
			return
		}

		// Check if the service has already been added.
		d, err := device.FindByDisplayName(service.Name, w.LMClient)
		if err != nil {
			log.Errorf("Failed to find service %q: %v", service.Name, err)
			return
		}

		// Add the service.
		if d == nil {
			newDevice := w.makeDeviceObject(service)
			err = device.Add(newDevice, w.LMClient)
			if err != nil {
				log.Errorf("Failed to add service %q: %v", newDevice.DisplayName, err)
			}
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		oldService := oldObj.(*v1.Service)
		newService := newObj.(*v1.Service)

		// Only add the service if it is of the type 'ClusterIP'.
		if newService.Spec.Type != v1.ServiceTypeClusterIP {
			return
		}

		// If the old service does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		if oldService.Spec.ClusterIP == "" && newService.Spec.ClusterIP != "" {
			d := w.makeDeviceObject(newService)
			err := device.Add(d, w.LMClient)
			if err != nil {
				log.Errorf("Failed to add service %s: %s", d.DisplayName, err)
			}
			log.Infof("Added service %s", d.DisplayName)
			return
		}

		// Covers the case when a service has been terminated (new ip doesn't exist)
		// and if a service needs to be added.
		if oldService.Spec.ClusterIP != "" && newService.Spec.ClusterIP != "" {
			oldDevice, err := device.FindByDisplayName(oldService.Name, w.LMClient)
			if err != nil {
				log.Errorf("Failed to find service %q: %v", oldService.Name, err)
				return
			}

			// Update the service.
			if oldDevice != nil {
				newDevice := w.makeDeviceObject(newService)
				err := device.UpdateAndReplace(newDevice, oldDevice.Id, w.LMClient)
				if err != nil {
					log.Errorf("Failed to update service %s: %v", oldDevice.DisplayName, err)
					return
				}
			}

			log.Infof("Updated service %s with id %d", oldDevice.DisplayName, oldDevice.Id)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		d, err := device.FindByDisplayName(service.Name, w.LMClient)
		if err != nil {
			log.Errorf("Failed to find service %q: %v", service.Name, err)
			return
		}

		if d == nil {
			return
		}

		// Delete the device
		if w.Config.DeleteDevices {
			err = device.Delete(d, w.LMClient)
			if err != nil {
				log.Errorf("Failed to delete service: %v", err)
			}
			log.Infof("Deleted service %s with id %d", d.DisplayName, d.Id)
			return
		}

		// Move the device

		categories := device.BuildSystemCategoriesFromLabels(constants.ServiceDeletedCategory, service.Labels)
		newDevice := &lm.RestDevice{
			CustomProperties: []lm.NameAndValue{
				{
					Name:  "system.categories",
					Value: categories,
				},
			},
		}
		err = device.UpdateAndReplaceField(newDevice, d.Id, constants.CustomPropertiesFieldName, w.LMClient)
		if err != nil {
			log.Errorf("Failed to move service %s: %s", d.DisplayName, err)
			return
		}

		log.Infof("Moved service %s with id %d to deleted group", d.DisplayName, d.Id)
	}
}

func (w Watcher) makeDeviceObject(service *v1.Service) *lm.RestDevice {
	categories := device.BuildSystemCategoriesFromLabels(constants.ServiceCategory, service.Labels)

	fqdn := service.Name + "." + service.Namespace + ".svc.cluster.local"
	d := &lm.RestDevice{
		Name:                 fqdn,
		DisplayName:          fqdn + "-" + string(service.UID),
		DisableAlerting:      true,
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
				Value: service.SelfLink,
			},
			{
				Name:  "auto.name",
				Value: service.Name,
			},
			{
				Name:  "auto.namespace",
				Value: service.Namespace,
			},
			{
				Name:  "auto.uid",
				Value: string(service.UID),
			},
		},
	}

	return d
}
