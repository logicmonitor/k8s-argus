// Package service provides the logic for mapping a Kubernetes Service to a
// LogicMonitor w.
package service

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

// Watcher represents a watcher type that watches services.
type Watcher struct{}

// AddFuncOptions addfunc options
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, error) {
		if rt != enums.Services {
			return []types.DeviceOption{}, fmt.Errorf("resourceType is not of type services")
		}
		service := obj.(*corev1.Service) // nolint: forcetypeassert
		if service.Spec.ClusterIP == "" {
			return []types.DeviceOption{}, fmt.Errorf("empty Spec.ClusterIP")
		}

		var options []types.DeviceOption

		// headless services set clusterip to None: https://kubernetes.io/docs/concepts/services-networking/service/#headless-services
		// do not replace Name property, keep it as default name-svc-namespace
		if service.Spec.ClusterIP != "None" {
			options = append(options, b.Name(service.Spec.ClusterIP))
		}

		return options, nil
	}
}

// UpdateFuncOptions update options
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.DeviceBuilder) ([]types.DeviceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, b types.DeviceBuilder) ([]types.DeviceOption, bool, error) {
		return []types.DeviceOption{}, false, nil
	}
}

// DeleteFuncOptions delete options
func (w *Watcher) DeleteFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.DeviceOption {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.DeviceOption {
		return []types.DeviceOption{}
	}
}
