package replicaset

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	appsv1 "k8s.io/api/apps/v1"
)

// Watcher represents a watcher type that watches replicasets.
type Watcher struct{}

// AddFuncOptions addfunc options
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
		if rt != enums.ReplicaSets {
			return []types.ResourceOption{}, fmt.Errorf("resourceType is not of type replicasets")
		}
		rs := obj.(*appsv1.ReplicaSet)

		options := []types.ResourceOption{
			b.Custom(constants.SelectorCustomProperty, utilities.GenerateSelectorExpression(rs.Spec.Selector)),
			b.Custom(constants.SelectorCustomProperty+constants.AppliesToPropSuffix, utilities.GenerateSelectorAppliesTo(rs.Spec.Selector)),
		}

		return options, nil
	}
}

// UpdateFuncOptions update
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.ResourceMeta, types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, cacheMeta types.ResourceMeta, b types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
		oldReplicaset := oldObj.(*appsv1.ReplicaSet) // nolint: forcetypeassert
		rs := newObj.(*appsv1.ReplicaSet)            // nolint: forcetypeassert
		options := make([]types.ResourceOption, 0)

		// If MatchLabels of new & old daemonsets are different, append in options
		oldSelectorExpr := utilities.GenerateSelectorExpression(oldReplicaset.Spec.Selector)
		newSelectorExpr := utilities.GenerateSelectorExpression(rs.Spec.Selector)
		if oldSelectorExpr != newSelectorExpr {
			options = append(options, b.Custom(constants.SelectorCustomProperty, newSelectorExpr))
		}
		oldSelectorAppliesTo := utilities.GenerateSelectorAppliesTo(oldReplicaset.Spec.Selector)
		newSelectorAppliesTo := utilities.GenerateSelectorAppliesTo(rs.Spec.Selector)
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
