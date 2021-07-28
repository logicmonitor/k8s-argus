package resourcegroup

import (
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
)

func (m *Manager) setResourceGroupCache(lctx *lmctx.LMContext, resourceGroup *models.DeviceGroup) {
	key := types.ResourceName{
		Name:     *resourceGroup.Name,
		Resource: enums.Namespaces,
	}
	meta, _ := util.GetResourceMetaFromDeviceGroup(resourceGroup)
	m.ResourceCache.Set(lctx, key, meta)
}

// UnsetLMIDInCache unset , costly and performance impacting operation, use this cautiously - where it will get called very rarely.
func (m *Manager) UnsetLMIDInCache(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) bool {
	return m.ResourceCache.UnsetLMID(lctx, rt, id)
}
