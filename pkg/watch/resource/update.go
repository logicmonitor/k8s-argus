package resource

import (
	"reflect"
	"runtime/debug"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
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
		objectMeta := *rt.ObjectMeta(newObj)
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
func UpdateFuncDispatcher(
	updateFunc types.UpdatePreprocessFunc,
) types.WatcherUpdateFunc {
	return func(oldObj, newObj interface{}) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Panic for %s: %s", util.GetCurrentFunctionName(), r)
				logrus.Debugf("%s", debug.Stack())
			}
		}()

		logrus.Tracef("%s called for previous object %v and new object %v", util.GetCurrentFunctionName(), oldObj, newObj)

		rt, done := inferResourceType(newObj)
		if done {
			return
		}

		// No need to put old object context, because name, namespace and type are stagnant fields, never editable
		lctx := getRootContext(rt, newObj, "update")
		log := lmlog.Logger(lctx)

		log.Debugf("Received update event")
		rt.ObjectMeta(newObj).ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		updateFunc(lctx, rt, oldObj, newObj)
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
		meta := rt.ObjectMeta(newObj)
		if cacheMeta, ok := resourceCache.Exists(lctx, cache.ResourceName{
			Name:     meta.Name,
			Resource: rt,
		}, meta.Namespace); ok && cacheMeta.UID != meta.UID {
			conf, err := config.GetConfig()
			if err == nil {
				log.Infof("Deleting previous resource (%d) with old UID (%s)", cacheMeta.LMID, cacheMeta.UID)
				options := b.GetDefaultsDeviceOptions(rt, meta, conf)
				options = append(options, b.Auto("uid", string(cacheMeta.UID)))
				err = deleteFun(lctx, rt, newObj, options...)
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
func UpsertBasedOnCache(
	resourceCache types.ResourceCache,
	configurer types.WatcherConfigurer,
	actions types.Actions,
	b types.DeviceBuilder,
) types.UpdateProcessFunc {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj interface{}, newObj interface{}, oldOptions []types.DeviceOption, options []types.DeviceOption) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("Failed to get config: %s", err)
		}
		deviceObj, err2 := util.BuildDevice(lctx, conf, nil, options...)
		if err2 != nil {
			log.Errorf("Failed to build device: %s", err2)

			return
		}
		ce, ok := util.DoesDeviceExistInCacheUtil(lctx, rt, resourceCache, deviceObj)
		if !ok {
			log.Debugf("Missing device, adding it")
			deviceOptions, err := configurer.AddFuncOptions()(lctx, rt, newObj, b)
			if err != nil {
				log.Errorf("add: failed to get device additional options: %s", err)
			}

			options := append(options, deviceOptions...)
			actions.AddFunc()(lctx, rt, newObj, options...) // nolint: errcheck

			return
		}
		updateOptions, needDelete, err := configurer.UpdateFuncOptions()(lctx, rt, oldObj, newObj, b)

		if needDelete {
			log.Infof("Deleting device, if any")
			deleteOptions := configurer.DeleteFuncOptions()(lctx, rt, newObj)
			options = append(options, deleteOptions...)
			actions.DeleteFunc()(lctx, rt, newObj, options...) // nolint: errcheck

			return
		}
		if err != nil {
			log.Errorf("update: failed to get device additional options: %s", err)

			return
		}

		options = append(options, updateOptions...)
		delta := hasDelta(lctx, rt, ce, newObj)
		if delta || len(updateOptions) > 0 {
			log.Infof("Updating device")
			actions.UpdateFunc()(lctx, rt, oldObj, newObj, options...) // nolint: errcheck

			return
		}

		log.Debugf("No change in data, ignoring update")
	}
}

func hasDelta(lctx *lmctx.LMContext, resourceType enums.ResourceType, cacheMeta cache.ResourceMeta, newObj interface{}) bool {
	log := lmlog.Logger(lctx)
	objMeta := resourceType.ObjectMeta(newObj)
	log.Tracef("Existing labels: %v new lables: %v", cacheMeta.Labels, objMeta.Labels)
	if cacheMeta.Labels != nil && objMeta.Labels != nil && len(cacheMeta.Labels) > 0 && len(objMeta.Labels) > 0 {
		return !reflect.DeepEqual(cacheMeta.Labels, objMeta.Labels)
	}

	return false
}
