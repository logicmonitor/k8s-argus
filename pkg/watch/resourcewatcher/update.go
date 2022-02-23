package resourcewatcher

import (
	"errors"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
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

// UpdateFuncWithExclude update
func UpdateFuncWithExclude(
	updateFunc types.UpdatePreprocessFunc,
	deleteFunc types.DeletePreprocessFunc,
) func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}) {
		log := lmlog.Logger(lctx)
		objectMeta, _ := rt.ObjectMeta(newObj)
		exclude, err := EvaluateResourceExclusion(lctx, rt, objectMeta)
		// NOTE: non nil err not considered for returning back to caller, exclude flag will decide it. err can be non nil for subset of rules
		if err != nil {
			log.Debugf("Error occurred while evaluating exclude rules %s", err)
		}
		if exclude {
			log.Debugf("Excluding resource from monitoring, sending delete")
			deleteFunc(lctx, rt, oldObj)
			deleteFunc(lctx, rt, newObj)

			return
		}

		log.Tracef("Resource exclusion evaluated to false")
		updateFunc(lctx, rt, oldObj, newObj)
	}
}

// UpdateFuncDispatcher update
func UpdateFuncDispatcher(facade eventprocessor.RunnerFacade, updateFunc types.UpdatePreprocessFunc) types.WatcherUpdateFunc {
	return func(oldObj, newObj interface{}) {
		lctx := lmlog.NewLMContextWith(logrus.WithTime(time.Now()))
		log := lmlog.Logger(lctx)
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Panic for %s: %s", util.GetCurrentFunctionName(), r)
				log.Errorf("%s", debug.Stack())
			}
		}()

		log.Tracef("%s called for previous object %v and new object %v", util.GetCurrentFunctionName(), oldObj, newObj)

		rt, ok := InferResourceType(newObj)
		if !ok {
			log.Tracef("Cannot infer object type")
			return
		}

		// No need to put old object context, because name, namespace and type are stagnant fields, never editable
		lctx = getRootContext(lctx, rt, newObj, "update")
		log = lmlog.Logger(lctx)

		log.Debugf("Received update event")
		meta, _ := rt.ObjectMeta(newObj)
		meta.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		sendToFacade(facade, lctx, rt, "update", func() {
			updateFunc(lctx, rt, oldObj, newObj)
		})
	}
}

// PreprocessUpdateEventForOldUID deletes previous resource by correlating UID before calling next update function.
func PreprocessUpdateEventForOldUID(
	resourceCache types.ResourceCache,
	deleteFun types.ExecDeleteFunc,
	b *builder.Builder,
	update types.UpdatePreprocessFunc,
) func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj interface{}, newObj interface{}) {
		log := lmlog.Logger(lctx)
		meta, _ := rt.ObjectMeta(newObj)
		container := meta.Namespace
		if !rt.IsNamespaceScopedResource() {
			container = constants.ClusterScopedGroupName
		}
		if cacheMeta, ok := resourceCache.Exists(lctx, types.ResourceName{
			Name:     meta.Name,
			Resource: rt,
		}, container, true); ok && cacheMeta.UID != meta.UID {
			conf, err := config.GetConfig(lctx)
			if err == nil {
				log.Infof("Deleting previous resource (%d) with old UID (%s)", cacheMeta.LMID, cacheMeta.UID)
				options := b.GetDefaultsResourceOptions(rt, meta, conf)
				options = append(options,
					b.Auto("uid", string(cacheMeta.UID)),
					b.Name(cacheMeta.Name),
				)
				delLctx := lmlog.LMContextWithLMResourceID(lctx, cacheMeta.LMID)
				err = deleteFun(delLctx, rt, newObj, options...)
				if err != nil {
					log.Errorf("Failed to delete previous resource: %s", err)

					return
				}
			} else {
				log.Errorf("Cannot delete previous resource: %s", err)

				return
			}
		}
		update(lctx, rt, oldObj, newObj)
	}
}

// UpsertBasedOnCache upsert
func UpsertBasedOnCache( //nolint:cyclop
	resourceCache types.ResourceCache,
	configurer types.WatcherConfigurer,
	actions types.Actions,
	b types.ResourceBuilder,
) types.UpdateProcessFunc {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj interface{}, newObj interface{}, oldOptions []types.ResourceOption, options []types.ResourceOption) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig(lctx)
		if err != nil {
			log.Errorf("Failed to get config: %s", err)
		}
		resourceObj, err2 := util.BuildResource(lctx, conf, nil, options...)
		if err2 != nil {
			log.Errorf("Failed to build resource: %s", err2)

			return
		}
		ce, ok := util.DoesResourceExistInCacheUtil(lctx, rt, resourceCache, resourceObj, false)
		if !ok {
			log.Debugf("Missing resource, adding it")
			resourceOptions, err := configurer.AddFuncOptions()(lctx, rt, newObj, b)
			if err != nil {
				if errors.Is(err, aerrors.ErrPodSucceeded) {
					log.Warnf("add: pod having succeeded status will not be considered for monitoring: %s", err)
					return
				}
				log.Errorf("add: failed to get resource additional options: %s", err)
			}

			options := append(options, resourceOptions...)
			actions.AddFunc()(lctx, rt, newObj, options...) // nolint: errcheck

			return
		}
		upLctx := lmlog.LMContextWithLMResourceID(lctx, ce.LMID)
		log = lmlog.Logger(upLctx)
		updateOptions, needDelete, err := configurer.UpdateFuncOptions()(upLctx, rt, oldObj, newObj, ce, b)

		if needDelete {
			log.Infof("Deleting resource, if any")
			deleteOptions := configurer.DeleteFuncOptions()(upLctx, rt, newObj)
			options = append(options, deleteOptions...)
			actions.DeleteFunc()(upLctx, rt, newObj, options...) // nolint: errcheck

			return
		}

		if err != nil {
			switch {
			case errors.Is(err, aerrors.ErrNoChangeInUpdateOptions):
				log.Warnf("update: no change in update options: %s", err)
			case errors.Is(err, aerrors.ErrPodSucceeded):
				log.Warnf("update: pod having succeeded status will not be considered for monitoring: %s", err)
			default:
				log.Errorf("update: failed to get additional options: %s", err)
			}

			return
		}

		options = append(options, updateOptions...)
		delta := hasDelta(upLctx, rt, ce, newObj)
		if delta || len(updateOptions) > 0 {
			log.Infof("Updating resource")
			actions.UpdateFunc()(upLctx, rt, oldObj, newObj, options...) // nolint: errcheck

			return
		}

		log.Debugf("No change in data, ignoring update")
	}
}

func hasDelta(lctx *lmctx.LMContext, resourceType enums.ResourceType, cacheMeta types.ResourceMeta, newObj interface{}) bool {
	log := lmlog.Logger(lctx)
	objMeta, _ := resourceType.ObjectMeta(newObj)
	log.Tracef("Existing labels: %v new lables: %v", cacheMeta.Labels, objMeta.Labels)
	if cacheMeta.Labels != nil && objMeta.Labels != nil && len(cacheMeta.Labels) > 0 && len(objMeta.Labels) > 0 {
		return !reflect.DeepEqual(cacheMeta.Labels, objMeta.Labels)
	}

	return false
}
