package resourcewatcher

import (
	"runtime/debug"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/eventprocessor"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resource/builder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddFuncWithExclude add with exclude
func AddFuncWithExclude(
	addFunc types.AddPreprocessFunc,
	deleteFunc types.DeletePreprocessFunc,
) types.AddPreprocessFunc {
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
func AddOrUpdateFunc(
	holders map[enums.ResourceType]*types.ControllerInitSyncStateHolder,
	addFunc types.AddPreprocessFunc,
	updateFunc types.UpdatePreprocessFunc,
) types.AddPreprocessFunc {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		if holders[rt] != nil && !holders[rt].HasSynced() {
			log.Debugf("Initial bulk discovery is in progress")
			updateFunc(lctx, rt, obj, obj)
		}

		addFunc(lctx, rt, obj)
	}
}

// AddFuncDispatcher add dispatcher
func AddFuncDispatcher(facade eventprocessor.RunnerFacade, addFunc types.AddPreprocessFunc) types.WatcherAddFunc {
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

		lctx = getRootContext(lctx, rt, obj, "add")
		log = lmlog.Logger(lctx)

		log.Debugf("Received add event")
		rt.ObjectMeta(obj).ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		sendToFacade(facade, lctx, func() {
			addFunc(lctx, rt, obj)
		})
	}
}

// PreprocessAddEventForOldUID deletes previous resource by correlating UID before calling next add function.
func PreprocessAddEventForOldUID(
	resourceCache types.ResourceCache,
	deleteFun types.ExecDeleteFunc,
	b *builder.Builder,
	add types.AddPreprocessFunc,
) types.AddPreprocessFunc {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		meta := rt.ObjectMeta(obj)
		if cacheMeta, ok := resourceCache.Exists(lctx, types.ResourceName{
			Name:     meta.Name,
			Resource: rt,
		}, meta.Namespace, false); ok && cacheMeta.UID != meta.UID {
			conf, err := config.GetConfig()
			if err == nil {
				log.Infof("Deleting previous resource (%d) with obj UID (%s)", cacheMeta.LMID, cacheMeta.UID)
				options := b.GetDefaultsResourceOptions(rt, meta, conf)
				options = append(options, b.Auto("uid", string(cacheMeta.UID)))
				err = deleteFun(lctx, rt, obj, options...)
				if err != nil {
					log.Errorf("Failed to delete previous resource: %s", err)

					return
				}
			} else {
				log.Errorf("Cannot delete previous resource: %s", err)

				return
			}
		}
		add(lctx, rt, obj)
	}
}
