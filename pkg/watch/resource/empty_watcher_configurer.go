// Package resource pod provides the logic for mapping a Kubernetes Pod to a
// LogicMonitor w.
package resource

import (
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
)

// Watcher represents a watcher type that watches pods.
type emptyWatcher struct{}

// AddFuncOptions add
func (w *emptyWatcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, error) {
		return []types.DeviceOption{}, nil
	}
}

// UpdateFuncOptions update
func (w *emptyWatcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.DeviceBuilder) ([]types.DeviceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, bool, error) {
		return []types.DeviceOption{}, false, nil
	}
}

// DeleteFuncOptions delete
func (w *emptyWatcher) DeleteFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.DeviceOption {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.DeviceOption {
		return []types.DeviceOption{}
	}
}
