// Package node provides the logic for mapping a Kubernetes Node to a
// LogicMonitor device.
package node

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

// Watcher represents a watcher type that watches nodes.
type Watcher struct{}

// AddFuncOptions add
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, error) {
		if rt != enums.Nodes {
			return []types.DeviceOption{}, fmt.Errorf("resourceType is not of type nodes")
		}
		node := obj.(*corev1.Node) // nolint: forcetypeassert
		internalAddress := getInternalAddress(node.Status.Addresses)
		if internalAddress == nil {
			return []types.DeviceOption{}, fmt.Errorf("no internal ip address present")
		}

		options := []types.DeviceOption{
			b.Name(internalAddress.Address),
		}

		return options, nil
	}
}

// UpdateFuncOptions update
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.DeviceBuilder) ([]types.DeviceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, bool, error) {
		if rt != enums.Nodes {
			return []types.DeviceOption{}, false, fmt.Errorf("resourceType is not of type nodes")
		}
		oldNode := oldObj.(*corev1.Node) // nolint: forcetypeassert
		node := newObj.(*corev1.Node)    // nolint: forcetypeassert
		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		oldInternalAddress := getInternalAddress(oldNode.Status.Addresses)
		internalAddress := getInternalAddress(node.Status.Addresses)

		var err error
		var options []types.DeviceOption
		if internalAddress == nil {
			err = fmt.Errorf("no internal ip address present")
		} else if oldInternalAddress.Address != internalAddress.Address {
			options = append(options, b.Name(internalAddress.Address))
		}

		if err != nil {
			return options, false, err
		}

		if !cmp.Equal(oldNode.Labels, node.Labels) {
			options = append(options, b.ResourceLabels(node.Labels))
		}

		if len(options) > 0 {
			return options, false, err
		}

		return options, false, fmt.Errorf("no change in additional options")
	}
}

// DeleteFuncOptions delete
func (w *Watcher) DeleteFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.DeviceOption {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.DeviceOption {
		return []types.DeviceOption{}
	}
}

// getInternalAddress finds the node's internal address.
func getInternalAddress(addresses []corev1.NodeAddress) *corev1.NodeAddress {
	var hostname *corev1.NodeAddress

	for i := range addresses {
		address := addresses[i]
		if address.Type == corev1.NodeInternalIP {
			return &address
		}
		if address.Type == corev1.NodeHostName {
			// if there is no internal IP for this node, the host name will be used
			hostname = &address
		}
	}

	return hostname
}
