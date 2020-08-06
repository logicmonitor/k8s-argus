package lmexec

import (
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
)

// LMExec Provides utility function for SDK calls using Base object
type LMExec struct {
	*types.Base
}

// AddDevice Add new device
func (lmexec *LMExec) AddDevice(params *lm.AddDeviceParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.AddDevice(params)
	}
}

// UpdateDevice Add new device
func (lmexec *LMExec) UpdateDevice(params *lm.UpdateDeviceParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.UpdateDevice(params)
	}
}

// GetDeviceList Add new device
func (lmexec *LMExec) GetDeviceList(params *lm.GetDeviceListParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.GetDeviceList(params)
	}
}

// PatchDevice Add new device
func (lmexec *LMExec) PatchDevice(params *lm.PatchDeviceParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.PatchDevice(params)
	}
}

// DeleteDeviceByID Add new device
func (lmexec *LMExec) DeleteDeviceByID(params *lm.DeleteDeviceByIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.DeleteDeviceByID(params)
	}
}

// GetImmediateDeviceListByDeviceGroupID Add new device
func (lmexec *LMExec) GetImmediateDeviceListByDeviceGroupID(params *lm.GetImmediateDeviceListByDeviceGroupIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.GetImmediateDeviceListByDeviceGroupID(params)
	}
}
