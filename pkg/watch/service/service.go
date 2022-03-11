// Package service provides the logic for mapping a Kubernetes Service to a
// LogicMonitor w.
package service

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
)

// Watcher represents a watcher type that watches services.
type Watcher struct{}

// AddFuncOptions addfunc options
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
		if rt != enums.Services {
			return []types.ResourceOption{}, fmt.Errorf("resourceType is not of type services")
		}
		svc := obj.(*corev1.Service) // nolint: forcetypeassert
		if svc.Spec.ClusterIP == "" {
			return []types.ResourceOption{}, fmt.Errorf("empty Spec.ClusterIP")
		}

		options := []types.ResourceOption{
			b.Custom(constants.SelectorCustomProperty, utilities.GenerateSelectorExpression(svc.Spec.Selector)),
			b.Custom(constants.SelectorCustomProperty+constants.AppliesToPropSuffix, utilities.GenerateSelectorAppliesTo(svc.Spec.Selector)),
		}

		// headless services set clusterip to None: https://kubernetes.io/docs/concepts/services-networking/service/#headless-services
		// do not replace Name property, keep it as default name-svc-namespace
		if svc.Spec.ClusterIP != constants.HeadlessServiceIPNone {
			options = append(options, b.Name(svc.Spec.ClusterIP))
		} else {
			options = append(options, b.Name(rt.LMName(meta.AsPartialObjectMetadata(&svc.ObjectMeta))))
		}

		return options, nil
	}
}

// UpdateFuncOptions update options
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.ResourceMeta, types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, cacheMeta types.ResourceMeta, b types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
		oldService := oldObj.(*corev1.Service) // nolint: forcetypeassert
		svc := newObj.(*corev1.Service)        // nolint: forcetypeassert
		var options []types.ResourceOption
		if svc.Spec.ClusterIP != constants.HeadlessServiceIPNone && cacheMeta.Name != svc.Spec.ClusterIP {
			options = append(options, b.Name(svc.Spec.ClusterIP))
		} else if svc.Spec.ClusterIP == constants.HeadlessServiceIPNone && cacheMeta.Name != rt.LMName(meta.AsPartialObjectMetadata(&svc.ObjectMeta)) {
			options = append(options, b.Name(rt.LMName(meta.AsPartialObjectMetadata(&svc.ObjectMeta))))
		}

		// If MatchLabels of new & old daemonsets are different, append in options
		oldSelectorExpr := utilities.GenerateSelectorExpression(oldService.Spec.Selector)
		newSelectorExpr := utilities.GenerateSelectorExpression(svc.Spec.Selector)
		if oldSelectorExpr != newSelectorExpr {
			options = append(options, b.Custom(constants.SelectorCustomProperty, newSelectorExpr))
		}
		oldSelectorAppliesTo := utilities.GenerateSelectorAppliesTo(oldService.Spec.Selector)
		newSelectorAppliesTo := utilities.GenerateSelectorAppliesTo(svc.Spec.Selector)
		if oldSelectorAppliesTo != newSelectorAppliesTo {
			options = append(options, b.Custom(constants.SelectorCustomProperty, newSelectorAppliesTo))
		}
		return options, false, nil
	}
}

// DeleteFuncOptions delete options
func (w *Watcher) DeleteFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
		return []types.ResourceOption{}
	}
}
