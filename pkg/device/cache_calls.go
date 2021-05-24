package device

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// SetDeviceInCache set
func (m *Manager) SetDeviceInCache(lctx *lmctx.LMContext, rt enums.ResourceType, device *models.Device) bool {
	log := lmlog.Logger(lctx)
	if device == nil {
		return false
	}
	resourceName, err := util.GetResourceNameFromDevice(rt, device)
	if err != nil {
		return false
	}
	resourceMeta, err := util.GetResourceMetaFromDevice(device)
	if err != nil {
		log.Errorf("Failed to get meta data from device object: %s", err)

		return false
	}

	return m.ResourceCache.Set(resourceName, resourceMeta)
}

// UnsetDeviceInCache unset
func (m *Manager) UnsetDeviceInCache(lctx *lmctx.LMContext, rt enums.ResourceType, device *models.Device) bool {
	if device == nil {
		return false
	}
	resourceName, err := util.GetResourceNameFromDevice(rt, device)
	if err != nil {
		return false
	}

	return m.ResourceCache.Unset(resourceName, util.GetPropertyValue(device, constants.K8sDeviceNamespacePropertyKey))
}

// DoesDeviceExistInCache exists
func (m *Manager) DoesDeviceExistInCache(lctx *lmctx.LMContext, rt enums.ResourceType, device *models.Device) (cache.ResourceMeta, bool) {
	return util.DoesDeviceExistInCacheUtil(lctx, rt, m.ResourceCache, device)
}

// DoesDeviceConflictInCluster conflicts
func (m *Manager) DoesDeviceConflictInCluster(lctx *lmctx.LMContext, resource enums.ResourceType, device *models.Device) ([]cache.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	resourceName, err := util.GetResourceNameFromDevice(resource, device)
	if err != nil {
		log.Errorf("Failed to get resource key from device: %s", err)

		return []cache.ResourceMeta{}, false
	}
	log.Debugf("ResourceName: %s", resourceName)
	list, ok := m.ResourceCache.Get(lctx, resourceName)
	if !ok {
		log.Debugf("No entry found")

		return []cache.ResourceMeta{}, false
	}
	namespace := util.GetPropertyValue(device, constants.K8sDeviceNamespacePropertyKey)
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
