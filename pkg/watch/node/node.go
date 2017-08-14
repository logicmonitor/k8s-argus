// Package node provides the logic for mapping a Kubernetes Node to a
// LogicMonitor device.
package node

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "nodes"
)

// Watcher represents a watcher type that watches nodes.
type Watcher struct {
	types.DeviceManager
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &v1.Node{}
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		// Require an IP address.
		if getInternalAddress(node.Status.Addresses) == nil {
			return
		}
		w.add(node)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*v1.Node)
		new := newObj.(*v1.Node)
		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		oldInternalAddress := getInternalAddress(old.Status.Addresses)
		newInternalAddress := getInternalAddress(new.Status.Addresses)
		if oldInternalAddress == nil && newInternalAddress != nil {
			w.add(new)
		}
		// Covers the case when the old node is in the process of terminating
		// and the new node is coming up to replace it.
		if oldInternalAddress.Address != newInternalAddress.Address {
			w.update(old, new)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)

		// Delete the node.
		if w.Config().DeleteDevices {
			if err := w.DeleteByName(node.Name); err != nil {
				log.Errorf("Failed to delete node: %v", err)
				return
			}
			log.Infof("Deleted node %s", node.Name)
			return
		}

		// Move the node.
		w.move(node)
	}
}

// nolint: dupl
func (w *Watcher) add(node *v1.Node) {
	if _, err := w.Add(w.args(node, constants.NodeCategory)...); err != nil {
		log.Errorf("Failed to add node %q: %v", node.Name, err)
		return
	}
	log.Infof("Added node %q", node.Name)
}

func (w *Watcher) update(old, new *v1.Node) {
	if _, err := w.UpdateAndReplaceByName(old.Name, w.args(new, constants.NodeCategory)...); err != nil {
		log.Errorf("Failed to update node %q: %v", new.Name, err)
		return
	}
	log.Infof("Updated node %q", old.Name)
}

func (w *Watcher) move(node *v1.Node) {
	if _, err := w.UpdateAndReplaceFieldByName(node.Name, constants.CustomPropertiesFieldName, w.args(node, constants.NodeDeletedCategory)...); err != nil {
		log.Errorf("Failed to move node %q: %v", node.Name, err)
		return
	}
	log.Infof("Moved node %q", node.Name)
}

func (w *Watcher) args(node *v1.Node, category string) []types.DeviceOption {
	categories := utilities.BuildSystemCategoriesFromLabels(category, node.Labels)
	return []types.DeviceOption{
		w.Name(getInternalAddress(node.Status.Addresses).Address),
		w.DisplayName(node.Name),
		w.SystemCategories(categories),
		w.Auto("name", node.Name),
		w.Auto("selflink", node.SelfLink),
		w.Auto("uid", string(node.UID)),
	}
}

// getInternalAddress finds the node's internal address.
func getInternalAddress(addresses []v1.NodeAddress) *v1.NodeAddress {
	for _, address := range addresses {
		if address.Type == v1.NodeInternalIP {
			return &address
		}
	}

	return nil
}
