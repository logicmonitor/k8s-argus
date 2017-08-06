package node

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
	resource = "nodes"
)

// Watcher represents a watcher type that watches nodes.
type Watcher struct {
	*types.Base
}

// Resource is a function that implements the Watcher interface.
func (w Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w Watcher) ObjType() runtime.Object {
	return &v1.Node{}
}

// AddFunc is a function that implements the Watcher interface.
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		if GetInternalAddress(node.Status.Addresses) != nil {
			w.addDevice(node)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		oldNode := oldObj.(*v1.Node)
		newNode := newObj.(*v1.Node)
		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		oldInternalAddress := GetInternalAddress(oldNode.Status.Addresses)
		newInternalAddress := GetInternalAddress(oldNode.Status.Addresses)
		if oldInternalAddress == nil && newInternalAddress != nil {
			w.addDevice(newNode)
		} else if oldInternalAddress.Address != newInternalAddress.Address {
			// Covers the case when a node has been terminated (new ip doesn't exist)
			// and if a node needs to be added.
			filter := fmt.Sprintf("displayName:%s", oldNode.Name)
			restResponse, apiResponse, err := w.LMClient.GetDeviceList("", -1, 0, filter)
			if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
				log.Errorf("Failed to find node %s: %s", oldNode.Name, _err)

				return
			}

			if restResponse.Data.Total == 1 {
				id := restResponse.Data.Items[0].Id
				w.updateDevice(newNode, id)
			}
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		filter := fmt.Sprintf("displayName:%s", node.Name)
		restResponse, apiResponse, err := w.LMClient.GetDeviceList("", -1, 0, filter)
		if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
			log.Errorf("Failed to find node %s: %s", node.Name, _err)

			return
		}
		if restResponse.Data.Total == 1 {
			id := restResponse.Data.Items[0].Id
			if w.Config.DeleteDevices {
				restResponse, apiResponse, err := w.LMClient.DeleteDevice(id)
				if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
					log.Errorf("Failed to delete device with id %q: %s", id, _err)
				}
				log.Infof("Deleted node %s with id %d", node.Name, id)

			} else {
				categories := constants.NodeDeletedCategory
				for k, v := range node.Labels {
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
					log.Errorf("Failed to move node %s: %s", node.Name, _err)

					return
				}
				log.Infof("Moved node %s with id %d to deleted group", node.Name, id)
			}
		}
	}
}

func (w Watcher) addDevice(node *v1.Node) {
	device := w.makeDeviceObject(node)
	restResponse, apiResponse, err := w.LMClient.AddDevice(device, false)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to add node %s: %s", node.Name, _err)

		return
	}
	log.Infof("Added node %s", node.Name)
}

func (w Watcher) updateDevice(node *v1.Node, id int32) {
	device := w.makeDeviceObject(node)
	restResponse, apiResponse, err := w.LMClient.UpdateDevice(device, id, "replace")
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to update node %s: %s", node.Name, _err)

		return
	}
	log.Infof("Updated node %s with id %d", node.Name, id)
}

func (w Watcher) makeDeviceObject(node *v1.Node) (device lm.RestDevice) {
	categories := constants.NodeCategory
	for k, v := range node.Labels {
		categories += "," + k + "=" + v

	}

	device = lm.RestDevice{
		Name:                 GetInternalAddress(node.Status.Addresses).Address,
		DisplayName:          node.Name,
		DisableAlerting:      w.Config.DisableAlerting,
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
				Value: node.SelfLink,
			},
			{
				Name:  "auto.name",
				Value: node.Name,
			},
			{
				Name:  "auto.uid",
				Value: string(node.UID),
			},
		},
	}

	return
}

// GetInternalAddress finds the node's internal address.
func GetInternalAddress(addresses []v1.NodeAddress) *v1.NodeAddress {
	for _, address := range addresses {
		if address.Type == v1.NodeInternalIP {
			return &address
		}
	}

	return nil
}
