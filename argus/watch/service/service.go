package service

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/logicmonitor/argus/argus/config"
	"github.com/logicmonitor/argus/constants"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "services"
)

// Watcher represents a watcher type that watches services.
type Watcher struct {
	LMClient  *lmv1.DefaultApi
	K8sClient *kubernetes.Clientset
	Config    *config.Config
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
			log.Infof("Adding service %s", service.Name)
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
			log.Infof("Adding service %s", newService.Name)
			w.addDevice(newService)
		} else if oldService.Spec.ClusterIP != "" && newService.Spec.ClusterIP != "" {
			// Covers the case when a service has been terminated (new ip doesn't exist)
			// and if a service needs to be added.
			filter := fmt.Sprintf("displayName:%s", oldService.Name)
			restResponse, _, err := w.LMClient.GetDeviceList("", -1, 0, filter)
			if err != nil {
				log.Errorf("Failed searching for service %s: %s", oldService.Name, restResponse.Errmsg)
			}
			if restResponse.Data.Total == 1 {
				id := restResponse.Data.Items[0].Id
				log.Infof("Updating service %s with id %d", oldService.Name, id)
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
		restResponse, _, err := w.LMClient.GetDeviceList("", -1, 0, filter)
		if err != nil {
			log.Errorf("Failed searching for service %s: %s", service.Name, restResponse.Errmsg)
		}
		if restResponse.Data.Total == 1 {
			id := restResponse.Data.Items[0].Id
			log.Infof("Deleting service %s with id %d", service.Name, id)
			w.LMClient.DeleteDevice(id)
		}
	}
}

func (w Watcher) addDevice(service *v1.Service) {
	device := w.makeDeviceObject(service)
	restResponse, _, err := w.LMClient.AddDevice(device, false)
	if err != nil {
		log.Error(restResponse.Errmsg)
	}
}

func (w Watcher) updateDevice(service *v1.Service, id int32) {
	device := w.makeDeviceObject(service)
	restResponse, _, err := w.LMClient.UpdateDevice(device, id, "")
	if err != nil {
		log.Error(restResponse.Errmsg)
	}
}

func (w Watcher) moveDevice(service *v1.Service, id int32) {
	device := w.makeDeviceObject(service)
	property := lmv1.NameAndValue{
		Name:  "auto.deleted",
		Value: "",
	}
	device.CustomProperties = append(device.CustomProperties, property)
	restResponse, _, err := w.LMClient.UpdateDevice(device, id, "")
	if err != nil {
		log.Error(restResponse.Errmsg)
	}
}

func (w Watcher) makeDeviceObject(service *v1.Service) (device lmv1.RestDevice) {
	categories := constants.ServiceCategory
	for k, v := range service.Labels {
		categories += "," + k + "=" + v

	}

	fqdn := service.Name + "." + service.Namespace + ".svc.cluster.local"
	device = lmv1.RestDevice{
		Name:                 fqdn,
		DisplayName:          fqdn + "-" + string(service.UID),
		DisableAlerting:      w.Config.DisableAlerting,
		HostGroupIds:         "1",
		PreferredCollectorId: w.Config.PreferredCollector,
		CustomProperties: []lmv1.NameAndValue{
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
