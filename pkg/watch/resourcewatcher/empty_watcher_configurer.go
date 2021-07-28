// Package resourcewatcher provides the logic for mapping a Kubernetes resource to a
// LogicMonitor resource.
package resourcewatcher

import (
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
)

// Watcher represents a watcher type that watches pods.
type emptyWatcher struct{}

// AddFuncOptions add
func (w *emptyWatcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
		return []types.ResourceOption{}, nil
	}
}

// UpdateFuncOptions update
func (w *emptyWatcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
		return []types.ResourceOption{}, false, nil
	}
}

// DeleteFuncOptions delete
func (w *emptyWatcher) DeleteFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
		return []types.ResourceOption{}
	}
}
