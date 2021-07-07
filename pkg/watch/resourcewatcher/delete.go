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
)

// DeleteFuncDispatcher delete
func DeleteFuncDispatcher(facade eventprocessor.RunnerFacade, deleteFunc types.DeletePreprocessFunc) types.WatcherDeleteFunc {
	return func(obj interface{}) {
		lctx := lmlog.NewLMContextWith(logrus.WithTime(time.Now()))
		log := lmlog.Logger(lctx)
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
		rt.ObjectMeta(obj).ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		RecordDeleteEventLatency(lctx, rt, obj)

		sendToFacade(facade, lctx, func() {
			deleteFunc(lctx, rt, obj)
		})
	}
}
