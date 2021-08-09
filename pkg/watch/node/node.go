// Package node provides the logic for mapping a Kubernetes Node to a
// LogicMonitor resource.
package node

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

// Watcher represents a watcher type that watches nodes.
type Watcher struct{}

// AddFuncOptions add
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
		if rt != enums.Nodes {
			return []types.ResourceOption{}, fmt.Errorf("resourceType is not of type nodes")
		}
		node := obj.(*corev1.Node) // nolint: forcetypeassert
		internalAddress := getInternalAddress(node.Status.Addresses)
		if internalAddress == nil {
			return []types.ResourceOption{}, fmt.Errorf("no internal ip address present")
		}

		options := []types.ResourceOption{
			b.Name(internalAddress.Address),
		}

		return options, nil
	}
}

// UpdateFuncOptions update
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
		if rt != enums.Nodes {
			return []types.ResourceOption{}, false, fmt.Errorf("resourceType is not of type nodes")
		}
		oldNode := oldObj.(*corev1.Node) // nolint: forcetypeassert
		node := newObj.(*corev1.Node)    // nolint: forcetypeassert
		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new resource.
		oldInternalAddress := getInternalAddress(oldNode.Status.Addresses)
		internalAddress := getInternalAddress(node.Status.Addresses)

		var err error
		var options []types.ResourceOption
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

		if !cmp.Equal(oldNode.Annotations, node.Annotations) {
			options = append(options, b.ResourceAnnotations(node.Annotations))
		}

		if len(options) > 0 {
			return options, false, err
		}

		return options, false, aerrors.ErrNoChangeInUpdateOptions
	}
}

// DeleteFuncOptions delete
func (w *Watcher) DeleteFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
		return []types.ResourceOption{}
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
