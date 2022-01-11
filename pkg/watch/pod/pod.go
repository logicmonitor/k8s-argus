// Package pod provides the logic for mapping a Kubernetes Pod to a
// LogicMonitor w.
package pod

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

// Watcher represents a watcher type that watches pods.
type Watcher struct{}

// AddFuncOptions addfunc options
func (w *Watcher) AddFuncOptions() func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, b types.ResourceBuilder) ([]types.ResourceOption, error) {
		if rt != enums.Pods {
			return []types.ResourceOption{}, fmt.Errorf("resourceType is not of type pods")
		}
		p := obj.(*corev1.Pod) // nolint: forcetypeassert
		if p.Status.PodIP == "" {
			return []types.ResourceOption{}, fmt.Errorf("empty Status.PodIP")
		}
		// If pod is in succeeded state, means it completed it execution
		// perhaps pods created for jobs, goes in succeeded state
		if p.Status.Phase == corev1.PodSucceeded {
			return []types.ResourceOption{}, aerrors.ErrPodSucceeded
		}
		options := []types.ResourceOption{
			b.Name(getPodDNSName(p)),
			b.System("ips", p.Status.PodIP),
		}

		// Pod running on fargate doesn't support HostNetwork so check fargate profile label, if label exists then mark hostNetwork as true
		if p.Spec.HostNetwork || p.Labels[constants.LabelFargateProfile] != "" {
			options = append(options, b.Custom("kubernetes.pod.hostNetwork", "true"))
		}

		options = append(options, b.Auto("nodename", p.Spec.NodeName))

		return options, nil
	}
}

// UpdateFuncOptions update
func (w *Watcher) UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, types.ResourceMeta, types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, cacheMeta types.ResourceMeta, b types.ResourceBuilder) ([]types.ResourceOption, bool, error) {
		oldPod := oldObj.(*corev1.Pod) // nolint: forcetypeassert
		p := newObj.(*corev1.Pod)      // nolint: forcetypeassert
		options := make([]types.ResourceOption, 0)
		// If pod is in succeeded state, means it completed it execution
		// perhaps pods created for jobs, goes in succeeded state
		if p.Status.Phase == corev1.PodSucceeded {
			return options, true, aerrors.ErrPodSucceeded
		}
		if p.Status.PodIP == "" {
			return options, false, fmt.Errorf("empty Status.PodIP")
		}
		if oldPod.Status.PodIP != p.Status.PodIP || (!p.Spec.HostNetwork && cacheMeta.Name != getPodDNSName(p)) {
			options = append(options, []types.ResourceOption{
				b.Name(getPodDNSName(p)),
				b.System("ips", p.Status.PodIP),
			}...)
		}

		if oldPod.Spec.HostNetwork != p.Spec.HostNetwork {
			options = append(options, b.Custom("kubernetes.pod.hostNetwork", "true"))
		}

		_, oldOk := oldPod.Labels[constants.LabelFargateProfile]
		_, newOk := p.Labels[constants.LabelFargateProfile]
		if !oldOk && newOk {
			// Pod running on fargate doesn't support HostNetwork so check fargate profile label, if label exists then mark hostNetwork as true
			options = append(options, b.Custom("kubernetes.pod.hostNetwork", "true"))
		}

		// TODO: Just kept in options, but auto properties are not allowed to update by santaba once updated.
		if oldPod.Spec.NodeName != p.Spec.NodeName {
			options = append(options, b.Auto("nodename", p.Spec.NodeName))
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

func getPodDNSName(pod *corev1.Pod) string {
	// if the pod is configured as "hostnetwork=true" or running on fargate, we will use the pod name as the IP/DNS name of the pod resource
	if pod.Spec.HostNetwork || pod.Labels[constants.LabelFargateProfile] != "" {
		return pod.Name
	}

	return pod.Status.PodIP
}
