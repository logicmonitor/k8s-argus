package resourcewatcher

import (
	"runtime/debug"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/eventprocessor"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// DeleteFuncDispatcher delete
func DeleteFuncDispatcher(facade eventprocessor.RunnerFacade, deleteFunc types.DeletePreprocessFunc) types.WatcherDeleteFunc {
	return func(obj interface{}) {
		lctx := lmlog.NewLMContextWith(logrus.WithTime(time.Now()))
		log := lmlog.Logger(lctx)
		if dfs, ok := obj.(cache.DeletedFinalStateUnknown); ok {
			logrus.Warnf("Delete event context is of type: %t", obj)
			// TODO: run partial sync for specified object key in the event: refer cache.DeletedFinalStateUnknown
			//  meanwhile continuing with stale object as its deletion so shouldn't be a problem
			obj = dfs.Obj
		}
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Panic for %s: %s", util.GetCurrentFunctionName(), r)
				log.Errorf("%s", debug.Stack())
			}
		}()

		log.Tracef("%s called for object %v", util.GetCurrentFunctionName(), obj)

		rt, ok := InferResourceType(obj)
		if !ok {
			log.Tracef("Cannot infer object type")
			return
		}

		lctx = getRootContext(lctx, rt, obj, "delete")

		log = lmlog.Logger(lctx)
		log.Debugf("Received delete event")
		meta, _ := rt.ObjectMeta(obj)
		meta.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		RecordDeleteEventLatency(lctx, rt, obj)

		sendToFacade(facade, lctx, rt, "delete", func() {
			deleteFunc(lctx, rt, obj)
		})
	}
}
