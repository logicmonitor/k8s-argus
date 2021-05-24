package resource

import (
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddFuncWithExclude add with exclude
func AddFuncWithExclude(
	addFunc func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}),
	deleteFunc func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{})) func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		objectMeta := *rt.ObjectMeta(obj)
		exclude, err := EvaluateResourceExclusion(lctx, rt, objectMeta)
		// NOTE: non nil err not considered for returning back to caller, exclude flag will decide it. err can be non nil for subset of rules
		if err != nil {
			log.Debugf("Error occurred while evaluating exclude rules %s", err)
		}
		if exclude {
			log.Debugf("Excluding resource from monitoring, sending delete")
			deleteFunc(lctx, rt, obj)

			return
		}

		log.Tracef("Resource exclusion evaluated to false")
		addFunc(lctx, rt, obj)
	}
}

// AddOrUpdateFunc add or update func only when resources are bulk discovered at the time of application restart
func AddOrUpdateFunc(holders map[enums.ResourceType]*types.ControllerInitSyncStateHolder,
	addFunc func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}),
	updateFunc func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{})) func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		if holders[rt] != nil && !holders[rt].HasSynced() {
			log.Infof("Initial bulk discovery is in progress")
			updateFunc(lctx, rt, obj, obj)
		}

		addFunc(lctx, rt, obj)
	}
}

// AddFuncDispatcher add dispatcher
func AddFuncDispatcher(addFunc func(
	lctx *lmctx.LMContext,
	rt enums.ResourceType,
	obj interface{},
)) func(obj interface{}) {
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

		lctx := getRootContext(rt, obj, "add")

		log := lmlog.Logger(lctx)
		log.Debugf("Received add event")
		rt.ObjectMeta(obj).ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		log.Tracef("Add event context: %s", spew.Sdump(obj))
		addFunc(lctx, rt, obj)
	}
}
