package resource

import (
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// SetResourceInCache set
func (m *Manager) SetResourceInCache(lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device) bool {
	log := lmlog.Logger(lctx)
	if resource == nil {
		return false
	}
	resourceName, err := util.GetResourceNameFromResource(rt, resource)
	if err != nil {
		return false
	}
	resourceMeta, err := util.GetResourceMetaFromResource(resource)
	if err != nil {
		log.Errorf("Failed to get meta data from resource object: %s", err)

		return false
	}

	return m.ResourceCache.Set(lctx, resourceName, resourceMeta)
}

// UnsetResourceInCache unset
func (m *Manager) UnsetResourceInCache(lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device) bool {
	if resource == nil {
		return false
	}
	resourceName, err := util.GetResourceNameFromResource(rt, resource)
	if err != nil {
		return false
	}

	return m.ResourceCache.Unset(lctx, resourceName, util.ResourceCacheContainerValue(resource))
}

// UnsetLMIDInCache unset , costly and performance impacting operation, use this cautiously - where it will get called very rarely.
func (m *Manager) UnsetLMIDInCache(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) bool {
	log := lmlog.Logger(lctx)
	log.Infof("Deleting cache entry using LMID: %d", id)
	return m.ResourceCache.UnsetLMID(lctx, rt, id)
}

// DoesResourceExistInCache exists
func (m *Manager) DoesResourceExistInCache(lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device, softRefresh bool) (types.ResourceMeta, bool) {
	return util.DoesResourceExistInCacheUtil(lctx, rt, m.ResourceCache, resource, softRefresh)
}

// DoesResourceConflictInCluster conflicts
func (m *Manager) DoesResourceConflictInCluster(lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device) ([]types.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	resourceName, err := util.GetResourceNameFromResource(rt, resource)
	if err != nil {
		log.Errorf("Failed to get resource key from resource: %s", err)

		return []types.ResourceMeta{}, false
	}
	log.Debugf("ResourceName: %s", resourceName)
	list, ok := m.ResourceCache.Get(lctx, resourceName)
	if !ok {
		log.Debugf("No entry found")

		return []types.ResourceMeta{}, false
	}
	namespace := util.ResourceCacheContainerValue(resource)
	log.Debugf("List of meta for %s: %v", resourceName, list)
	for idx, v := range list {
		if v.Container == namespace {
			list[idx] = list[len(list)-1]
			list = list[:len(list)-1]

			break
		}
	}
	if len(list) > 0 {
		return list, true
	}

	return list, false
}
