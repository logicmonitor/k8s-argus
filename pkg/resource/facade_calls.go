package resource

import (
	"fmt"
	"net/http"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

func (m *Manager) addResource(lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device) (*models.Device, error) {
	// Adding auto prop on resource created by argus
	key := constants.AutoPropCreatedBy
	val := fmt.Sprintf("%s%s", constants.CreatedByPrefix, constants.Version)
	resource.CustomProperties = append(resource.CustomProperties,
		&models.NameAndValue{
			Name:  &key,
			Value: &val,
		},
	)
	params := lm.NewAddDeviceParams()
	addFromWizard := false
	params.SetAddFromWizard(&addFromWizard)
	params.SetBody(resource)
	command := m.AddResourceCommand(lctx, params)
	restResponse, err := m.LMFacade.SendReceive(lctx, command)
	if err == nil {
		m.SetResourceInCache(lctx, rt, restResponse.(*lm.AddDeviceOK).Payload)

		return restResponse.(*lm.AddDeviceOK).Payload, nil
	}

	return nil, err
}

func (m *Manager) UpdateAndReplaceResource(lctx *lmctx.LMContext, rt enums.ResourceType, id int32, resource *models.Device) (*models.Device, error) {
	opType := "replace"
	// opType := "refresh"
	params := lm.NewUpdateDeviceParams()
	params.SetID(id)
	params.SetBody(resource)
	params.SetOpType(&opType)
	command := m.UpdateResourceCommand(lctx, params)
	restResponse, err := m.LMFacade.SendReceive(lctx, command)

	if err == nil {
		m.SetResourceInCache(lctx, rt, restResponse.(*lm.UpdateDeviceOK).Payload)

		return restResponse.(*lm.UpdateDeviceOK).Payload, nil
	}

	if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		m.UnsetResourceInCache(lctx, rt, resource)
		return nil, fmt.Errorf("resource being updated does not exist: %s: %w", err, aerrors.ErrInvalidCache)
	}

	return nil, err
}

// PatchResource implements types.DeviceManager.
func (m *Manager) PatchResource(lctx *lmctx.LMContext, rt enums.ResourceType, device *models.Device, fields string) (*models.Device, error) {
	params := lm.NewPatchDeviceParams()
	params.SetID(device.ID)
	params.SetPatchFields(&fields)
	params.SetBody(device)
	opType := "replace"
	params.SetOpType(&opType)

	cmd := m.PatchResourceCommand(lctx, params)

	restResponse, err := m.LMFacade.SendReceive(lctx, cmd)
	if err == nil {
		m.SetResourceInCache(lctx, rt, restResponse.(*lm.PatchDeviceOK).Payload)

		return restResponse.(*lm.PatchDeviceOK).Payload, nil
	}
	if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		m.UnsetResourceInCache(lctx, rt, device)
		return nil, fmt.Errorf("resource being patched does not exist: %s: %w", err, aerrors.ErrInvalidCache)
	}

	return nil, err
}

// deleteResource implements types.ResourceManager.
func (m *Manager) deleteResource(lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device) error {
	params := lm.NewDeleteDeviceByIDParams()
	params.SetID(resource.ID)
	command := m.DeleteResourceByIDCommand(lctx, params)
	_, err := m.LMFacade.SendReceive(lctx, command)
	if err == nil {
		m.UnsetResourceInCache(lctx, rt, resource)
	}
	if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		m.UnsetResourceInCache(lctx, rt, resource)
	}

	return err
}

// FetchResource get resource with lmid
func (m *Manager) FetchResource(lctx *lmctx.LMContext, rt enums.ResourceType, lmid int32) (*models.Device, error) {
	params := lm.NewGetDeviceByIDParams()
	params.SetID(lmid)
	command := m.GetResourceByIDCommand(lctx, params)
	restResponse, err := m.LMFacade.SendReceive(lctx, command)
	if err == nil {
		resource := restResponse.(*lm.GetDeviceByIDOK).Payload

		return resource, nil
	}
	if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		// calling UnsetLMID cache only on 404 to reduce ops
		m.UnsetLMIDInCache(lctx, rt, lmid)
		return nil, fmt.Errorf("failed to fetch resource: %s: %w", err, aerrors.ErrInvalidCache)
	}

	return nil, err
}

// DeleteByID implements types.ResourceManager.
func (m *Manager) DeleteByID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) error {
	params := lm.NewDeleteDeviceByIDParams()
	params.SetID(id)
	command := m.DeleteResourceByIDCommand(lctx, params)
	_, err := m.LMFacade.SendReceive(lctx, command)
	if err == nil || util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		// calling UnsetLMID cache only on 404 to reduce ops
		m.UnsetLMIDInCache(lctx, rt, id)
	}

	return err
}
