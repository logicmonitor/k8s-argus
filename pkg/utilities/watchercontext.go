package utilities

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
)

// WatcherContext adds metadata in LMContext
func WatcherContext(lctx *lmctx.LMContext, w types.Watcher) *lmctx.LMContext {
	lctx.Set(constants.IsPingDevice, isResourcePingEnabled(w.Resource()))
	return lctx
}

// isResourcePingEnabled returns true if resource can be pinged.
// While adding new resource, add resource condition here.
func isResourcePingEnabled(resource string) bool {
	switch resource {
	case constants.Pods:
		return true
	case constants.Deployments, constants.Services, constants.Nodes, constants.HorizontalPodAutoScalers:
		return false
	default:
		return true
	}
}
