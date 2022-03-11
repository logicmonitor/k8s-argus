package statefulset

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	appsv1 "k8s.io/api/apps/v1"
)

// Watcher represents a watcher type that watches statefulsets.
type Watcher struct{}

// AddFuncOptions addfunc options
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
		if rt != enums.StatefulSets {
			return []types.ResourceOption{}, fmt.Errorf("resourceType is not of type statefulsets")
		}
		sts := obj.(*appsv1.StatefulSet)

		options := []types.ResourceOption{
			b.Custom(constants.SelectorCustomProperty, utilities.GenerateSelectorExpression(sts.Spec.Selector)),
			b.Custom(constants.SelectorCustomProperty+constants.AppliesToPropSuffix, utilities.GenerateSelectorAppliesTo(sts.Spec.Selector)),
		}

		return options, nil
	}
}

// UpdateFuncOptions update
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.ResourceMeta, types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, cacheMeta types.ResourceMeta, b types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
		oldStatefulset := oldObj.(*appsv1.StatefulSet) // nolint: forcetypeassert
		sts := newObj.(*appsv1.StatefulSet)            // nolint: forcetypeassert
		options := make([]types.ResourceOption, 0)

		// If MatchLabels of new & old daemonsets are different, append in options
		oldSelectorExpr := utilities.GenerateSelectorExpression(oldStatefulset.Spec.Selector)
		newSelectorExpr := utilities.GenerateSelectorExpression(sts.Spec.Selector)
		if oldSelectorExpr != newSelectorExpr {
			options = append(options, b.Custom(constants.SelectorCustomProperty, newSelectorExpr))
		}
		oldSelectorAppliesTo := utilities.GenerateSelectorAppliesTo(oldStatefulset.Spec.Selector)
		newSelectorAppliesTo := utilities.GenerateSelectorAppliesTo(sts.Spec.Selector)
		if oldSelectorAppliesTo != newSelectorAppliesTo {
			options = append(options, b.Custom(constants.SelectorCustomProperty+constants.AppliesToPropSuffix, newSelectorAppliesTo))
		}
		return options, false, nil
	}
}

// DeleteFuncOptions delete
func (w *Watcher) DeleteFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) []types.ResourceOption {
		return []types.ResourceOption{}
	}
}
