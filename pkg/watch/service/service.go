package service

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
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
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)

		// Only add the service if it is of the type 'ClusterIP'.
		if service.Spec.Type != v1.ServiceTypeClusterIP {
			return
		}

		if service.Spec.ClusterIP != "" {
			w.addDevice(service)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
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
			w.addDevice(newService)
		} else if oldService.Spec.ClusterIP != "" && newService.Spec.ClusterIP != "" {
			// Covers the case when a service has been terminated (new ip doesn't exist)
			// and if a service needs to be added.
			filter := fmt.Sprintf("displayName:%s", oldService.Name)
			restResponse, apiResponse, err := w.LMClient.GetDeviceList("", -1, 0, filter)
			if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
				log.Errorf("Failed to find service %s: %s", oldService.Name, _err)

				return
			}
			if restResponse.Data.Total == 1 {
				id := restResponse.Data.Items[0].Id
				w.updateDevice(newService, id)
			}
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		filter := fmt.Sprintf("displayName:%s", service.Name)
		restResponse, apiResponse, err := w.LMClient.GetDeviceList("", -1, 0, filter)
		if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
			log.Errorf("Failed to find service %s: %s", service.Name, _err)

			return
		}
		if restResponse.Data.Total == 1 {
			id := restResponse.Data.Items[0].Id
			if w.Config.DeleteDevices {
				restResponse, apiResponse, err := w.LMClient.DeleteDevice(id)
				if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
					log.Printf("Failed to delete device with id %q: %s", id, _err)

					return
				}
				log.Infof("Deleted service %s with id %d", service.Name, id)
			} else {
				categories := constants.ServiceDeletedCategory
				for k, v := range service.Labels {
					categories += "," + k + "=" + v

				}
				device := lm.RestDevice{
					CustomProperties: []lm.NameAndValue{
						{
							Name:  "system.categories",
							Value: categories,
						},
					},
				}
				restResponse, apiResponse, err := w.LMClient.PatchDeviceById(device, id, "replace", "customProperties")
				if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
					log.Errorf("Failed to move service %s: %s", service.Name, _err)

					return
				}
				log.Infof("Moved service %s with id %d to deleted group", service.Name, id)
			}
		}
	}
}

func (w Watcher) addDevice(service *v1.Service) {
	device := w.makeDeviceObject(service)
	restResponse, apiResponse, err := w.LMClient.AddDevice(device, false)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to add service %s: %s", service.Name, _err)

		return
	}
	log.Infof("Added service %s", service.Name)
}

func (w Watcher) updateDevice(service *v1.Service, id int32) {
	device := w.makeDeviceObject(service)
	restResponse, apiResponse, err := w.LMClient.UpdateDevice(device, id, "replace")
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to update service %s: %s", service.Name, _err)

		return
	}
	log.Infof("Updated service %s with id %d", service.Name, id)
}

func (w Watcher) makeDeviceObject(service *v1.Service) (device lm.RestDevice) {
	categories := constants.ServiceCategory
	for k, v := range service.Labels {
		categories += "," + k + "=" + v

	}

	fqdn := service.Name + "." + service.Namespace + ".svc.cluster.local"
	device = lm.RestDevice{
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

	return
}
