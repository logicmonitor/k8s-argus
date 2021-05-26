package device

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

func (m *Manager) addDevice(lctx *lmctx.LMContext, rt enums.ResourceType, device *models.Device) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	log.Tracef("Adding device: %s", spew.Sdump(device))
	params := lm.NewAddDeviceParams()
	addFromWizard := false
	params.SetAddFromWizard(&addFromWizard)
	params.SetBody(device)
	cmd := &types.HTTPCommand{
		IsGlobal: false,
		Command: &types.Command{ // nolint: exhaustivestruct
			ExecFun: m.AddDevice(params),
			LMCtx:   lctx,
		},
		Method:   "POST",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.AddDeviceErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, rt, cmd)
	if err == nil {
		m.SetDeviceInCache(lctx, rt, restResponse.(*lm.AddDeviceOK).Payload)

		return restResponse.(*lm.AddDeviceOK).Payload, nil
	}

	return nil, err
}

func (m *Manager) UpdateAndReplaceResource(lctx *lmctx.LMContext, rt enums.ResourceType, id int32, device *models.Device) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	log.Tracef("Updating device: %s", spew.Sdump(device))
	opType := "replace"
	// opType := "refresh"
	params := lm.NewUpdateDeviceParams()
	params.SetID(id)
	params.SetBody(device)
	params.SetOpType(&opType)
	cmd := &types.HTTPCommand{
		IsGlobal: false,
		Command: &types.Command{ // nolint: exhaustivestruct
			ExecFun: m.UpdateDevice(params),
			LMCtx:   lctx,
		},
		Method:   "PUT",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.UpdateDeviceErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, rt, cmd)

	if err == nil {
		m.SetDeviceInCache(lctx, rt, restResponse.(*lm.UpdateDeviceOK).Payload)

		return restResponse.(*lm.UpdateDeviceOK).Payload, nil
	}

	return nil, err
}

// deleteDevice implements types.DeviceManager.
func (m *Manager) deleteDevice(lctx *lmctx.LMContext, rt enums.ResourceType, device *models.Device) error {
	log := lmlog.Logger(lctx)
	log.Tracef("Deleting device: %s", spew.Sdump(device))
	params := lm.NewDeleteDeviceByIDParams()
	params.SetID(device.ID)
	cmd := &types.HTTPCommand{
		IsGlobal: false,
		Command: &types.Command{ // nolint: exhaustivestruct
			ExecFun: m.DeleteDeviceByID(params),
			LMCtx:   lctx,
		},
		Method:   "DELETE",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.DeleteDeviceByIDErrResp,
		},
	}
	_, err := m.LMFacade.SendReceive(lctx, rt, cmd)
	if err == nil {
		m.UnsetDeviceInCache(lctx, rt, device)
	}
	if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		m.UnsetDeviceInCache(lctx, rt, device)
	}

	return err
}

// FetchDevice get device with lmid
func (m *Manager) FetchDevice(lctx *lmctx.LMContext, rt enums.ResourceType, lmid int32) (*models.Device, error) {
	params := lm.NewGetDeviceByIDParams()
	params.SetID(lmid)
	cmd := &types.HTTPCommand{
		IsGlobal: false,
		Command: &types.Command{ // nolint: exhaustivestruct
			ExecFun: m.GetDeviceByID(params),
			LMCtx:   lctx,
		},
		Method:   "GET",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.GetDeviceByIDErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, rt, cmd)
	if err == nil {
		device := restResponse.(*lm.GetDeviceByIDOK).Payload
		m.UnsetDeviceInCache(lctx, rt, device)

		return device, nil
	}

	return nil, err
}

// DeleteByID implements types.DeviceManager.
func (m *Manager) DeleteByID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) error {
	params := lm.NewDeleteDeviceByIDParams()
	params.SetID(id)
	cmd := &types.HTTPCommand{
		IsGlobal: false,
		Command: &types.Command{ // nolint: exhaustivestruct
			ExecFun: m.DeleteDeviceByID(params),
			LMCtx:   lctx,
		},
		Method:   "DELETE",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.DeleteDeviceByIDErrResp,
		},
	}
	_, err := m.LMFacade.SendReceive(lctx, rt, cmd)

	return err
}
