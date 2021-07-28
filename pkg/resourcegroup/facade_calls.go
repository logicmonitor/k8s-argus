package resourcegroup

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

func (m *Manager) addResourceGroup(lctx *lmctx.LMContext, resourceGroup *models.DeviceGroup) (interface{}, error) {
	params := lm.NewAddDeviceGroupParams()
	params.SetBody(resourceGroup)
	command := m.LMExecutor.AddResourceGroupCommand(lctx, params)
	resp, err := m.LMFacade.SendReceive(lctx, command)
	if err == nil {
		m.setResourceGroupCache(lctx, resp.(*lm.AddDeviceGroupOK).Payload)
	}
	return resp, err
}

func (m *Manager) updateResourceGroupByID(lctx *lmctx.LMContext, meta types.ResourceMeta, existingDeviceGroup *models.DeviceGroup) error {
	updateParams := lm.NewUpdateDeviceGroupByIDParams()
	updateParams.SetID(meta.LMID)
	updateParams.SetBody(existingDeviceGroup)
	updateCommand := m.UpdateResourceGroupByIDCommand(lctx, updateParams)
	resp, err := m.SendReceive(lctx, updateCommand)
	if err == nil {
		m.setResourceGroupCache(lctx, resp.(*lm.UpdateDeviceGroupByIDOK).Payload)
	}
	return err
}

func (m *Manager) getResourceGroupByID(id int32, clctx *lmctx.LMContext) (interface{}, error) {
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(id)
	command := m.GetResourceGroupByIDCommand(clctx, params)
	resp, err := m.SendReceive(clctx, command)
	return resp, err
}

func (m *Manager) getResourceGroupByName(clctx *lmctx.LMContext, parentID int32, name string) (*models.DeviceGroup, error) {
	params := lm.NewGetDeviceGroupListParams()
	filter := fmt.Sprintf("parentId:\"%v\",name:\"%s\"", parentID, name)
	params.SetFilter(&filter)
	command := m.GetResourceGroupListCommand(clctx, params)
	restResponse, err := m.SendReceive(clctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to get device group list when searching for %q: %w", name, err)
	}

	var deviceGroup *models.DeviceGroup
	if len(restResponse.(*lm.GetDeviceGroupListOK).Payload.Items) == 0 {
		return deviceGroup, fmt.Errorf("could not find device group %q with parentId %v", name, parentID)
	}
	deviceGroup = restResponse.(*lm.GetDeviceGroupListOK).Payload.Items[0]

	m.setResourceGroupCache(clctx, deviceGroup)

	return deviceGroup, nil
}
