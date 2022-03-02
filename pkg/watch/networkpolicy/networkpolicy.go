package networkpolicy

import (
	"fmt"
	"reflect"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	networkingv1 "k8s.io/api/networking/v1"
)

// Watcher represents a watcher type that watches networkpolicies.
type Watcher struct{}

// AddFuncOptions addfunc options
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
		if rt != enums.NetworkPolicies {
			return []types.ResourceOption{}, fmt.Errorf("resourceType is not of type networkpolicies")
		}
		netpol := obj.(*networkingv1.NetworkPolicy)

		options := []types.ResourceOption{
			b.Custom(
				constants.SelectorCustomPropertyPrefix+constants.PodSelectorKey,
				utilities.CoalesceMatchLabels(netpol.Spec.PodSelector.MatchLabels),
			),
		}

		return options, nil
	}
}

// UpdateFuncOptions update
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.ResourceMeta, types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, cacheMeta types.ResourceMeta, b types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
		oldNetworkPolicy := oldObj.(*networkingv1.NetworkPolicy) // nolint: forcetypeassert
		netpol := newObj.(*networkingv1.NetworkPolicy)           // nolint: forcetypeassert
		options := make([]types.ResourceOption, 0)

		// If MatchLabels of new & old daemonsets are different, add in append
		if !reflect.DeepEqual(oldNetworkPolicy.Spec.PodSelector.MatchLabels, netpol.Spec.PodSelector.MatchLabels) {
			options = append(options, b.Custom(
				constants.SelectorCustomPropertyPrefix+constants.PodSelectorKey,
				utilities.CoalesceMatchLabels(netpol.Spec.PodSelector.MatchLabels),
			))
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