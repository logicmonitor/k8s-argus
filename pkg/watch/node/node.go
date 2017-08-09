package node

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
// nolint: dupl
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)

		// We need an IP address.
		if GetInternalAddress(node.Status.Addresses) == nil {
			return
		}

		// Check if the node has already been added.
		d, err := device.FindByDisplayName(node.Name, w.LMClient)
		if err != nil {
			log.Errorf("Failed to find node %q: %v", node.Name, err)
			return
		}

		// Add the node.
		if d == nil {
			newDevice := w.makeDeviceObject(node)
			err = device.Add(newDevice, w.LMClient)
			if err != nil {
				log.Errorf("Failed to add node %q: %v", newDevice.DisplayName, err)
			}
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		oldNode := oldObj.(*v1.Node)
		newNode := newObj.(*v1.Node)

		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		oldInternalAddress := GetInternalAddress(oldNode.Status.Addresses)
		newInternalAddress := GetInternalAddress(oldNode.Status.Addresses)
		if oldInternalAddress == nil && newInternalAddress != nil {
			d := w.makeDeviceObject(newNode)
			err := device.Add(d, w.LMClient)
			if err != nil {
				log.Errorf("Failed to add node %s: %s", d.DisplayName, err)
			}
			log.Infof("Added node %s", d.DisplayName)
			return
		}

		// Covers the case when a node has been terminated (new ip doesn't exist)
		// and if a node needs to be added.
		if oldInternalAddress.Address != newInternalAddress.Address {
			oldDevice, err := device.FindByDisplayName(oldNode.Name, w.LMClient)
			if err != nil {
				log.Errorf("Failed to find node %q: %v", oldNode.Name, err)
				return
			}

			// Update the node.
			if oldDevice != nil {
				newDevice := w.makeDeviceObject(newNode)
				err := device.UpdateAndReplace(newDevice, oldDevice.Id, w.LMClient)
				if err != nil {
					log.Errorf("Failed to update node %s: %v", oldDevice.DisplayName, err)
					return
				}
			}

			log.Infof("Updated node %s with id %d", oldDevice.DisplayName, oldDevice.Id)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		d, err := device.FindByDisplayName(node.Name, w.LMClient)
		if err != nil {
			log.Errorf("Failed to find node %q: %v", node.Name, err)
			return
		}

		if d == nil {
			return
		}

		// Delete the device
		if w.Config.DeleteDevices {
			err = device.Delete(d, w.LMClient)
			if err != nil {
				log.Errorf("Failed to delete node: %v", err)
			}
			log.Infof("Deleted node %s with id %d", d.DisplayName, d.Id)
			return
		}

		// Move the device

		categories := device.BuildSystemCategoriesFromLabels(constants.NodeDeletedCategory, node.Labels)
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
			log.Errorf("Failed to move node %s: %s", d.DisplayName, err)
			return
		}

		log.Infof("Moved node %s with id %d to deleted group", d.DisplayName, d.Id)
	}
}

func (w Watcher) makeDeviceObject(node *v1.Node) *lm.RestDevice {
	categories := device.BuildSystemCategoriesFromLabels(constants.NodeCategory, node.Labels)

	d := &lm.RestDevice{
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

	return d
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
