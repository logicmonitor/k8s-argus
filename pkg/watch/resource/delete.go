package resource

import (
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteFuncDispatcher delete
func DeleteFuncDispatcher(
	deleteFunc types.DeletePreprocessFunc,
) types.WatcherDeleteFunc {
	return func(obj interface{}) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Panic for %s: %s", util.GetCurrentFunctionName(), r)
				logrus.Debugf("%s", debug.Stack())
			}
		}()

		logrus.Tracef("%s called for object %v", util.GetCurrentFunctionName(), obj)

		rt, done := inferResourceType(obj)
		if done {
			return
		}

		lctx := getRootContext(rt, obj, "delete")

		log := lmlog.Logger(lctx)
		log.Debugf("Received delete event")
		rt.ObjectMeta(obj).ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		log.Tracef("Delete event context: %s", spew.Sdump(obj))
		RecordDeleteEventLatency(lctx, rt, obj)
		deleteFunc(lctx, rt, obj)
	}
}
