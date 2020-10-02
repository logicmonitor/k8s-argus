package utilities

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
)

// WatcherContext adds metadata in LMContext
func WatcherContext(lctx *lmctx.LMContext, w types.Watcher) *lmctx.LMContext {
	lctx.Set(constants.IsPingEnabled, isResourcePingEnabled(w.Resource()))
	return lctx
}

// isResourcePingEnabled is a function to check whether resource can be pinged
func isResourcePingEnabled(resource string) bool {
	switch resource {
	case constants.Pods:
		return true
	default:
		return false
	}
}
