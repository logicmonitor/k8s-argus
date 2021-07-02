package lmexec

import (
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
)

// LMExec Provides utility function for SDK calls using Base object
// LMExec is holding resource related api calls at the moment to mitigate rate limit handling
type LMExec struct {
	LMClient *client.LMSdkGo
}

// addResource Add new resource
func (lmexec *LMExec) addResource(params *lm.AddDeviceParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.AddDevice(params)
	}
}

// updateResource Add new resource
func (lmexec *LMExec) updateResource(params *lm.UpdateDeviceParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.UpdateDevice(params)
	}
}

// getResourceByID get resource by id
func (lmexec *LMExec) getResourceByID(params *lm.GetDeviceByIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.GetDeviceByID(params)
	}
}

// getResourceList Add new resource
func (lmexec *LMExec) getResourceList(params *lm.GetDeviceListParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.GetDeviceList(params)
	}
}

// patchResource Add new resource
func (lmexec *LMExec) patchResource(params *lm.PatchDeviceParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.PatchDevice(params)
	}
}

// deleteResourceByID Add new resource
func (lmexec *LMExec) deleteResourceByID(params *lm.DeleteDeviceByIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.DeleteDeviceByID(params)
	}
}

// getImmediateResourceListByResourceGroupID Add new resource
func (lmexec *LMExec) getImmediateResourceListByResourceGroupID(params *lm.GetImmediateDeviceListByDeviceGroupIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.GetImmediateDeviceListByDeviceGroupID(params)
	}
}

// updateResourcePropertyByName updates specified resource property.
func (lmexec *LMExec) updateResourcePropertyByName(params *lm.UpdateDevicePropertyByNameParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.UpdateDevicePropertyByName(params)
	}
}

// updateResourcePropertyByName updates specified resource property.
func (lmexec *LMExec) addResourceGroup(params *lm.AddDeviceGroupParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.AddDeviceGroup(params)
	}
}

// updateDeviceGroupByIDGroup updates specified resource property.
func (lmexec *LMExec) updateResourceGroupByID(params *lm.UpdateDeviceGroupByIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.UpdateDeviceGroupByID(params)
	}
}

// addResourceGroupProperty updates specified resource property.
func (lmexec *LMExec) addResourceGroupProperty(params *lm.AddDeviceGroupPropertyParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.AddDeviceGroupProperty(params)
	}
}

// updateResourceGroupPropertyByName updates specified resource property.
func (lmexec *LMExec) updateResourceGroupPropertyByName(params *lm.UpdateDeviceGroupPropertyByNameParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.UpdateDeviceGroupPropertyByName(params)
	}
}

// deleteResourceGroupPropertyByName updates specified resource property.
func (lmexec *LMExec) deleteResourceGroupPropertyByName(params *lm.DeleteDeviceGroupPropertyByNameParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.DeleteDeviceGroupPropertyByName(params)
	}
}

// getResourceGroupByID updates specified resource property.
func (lmexec *LMExec) getResourceGroupByID(params *lm.GetDeviceGroupByIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.GetDeviceGroupByID(params)
	}
}

// deleteResourceGroupByID updates specified resource property.
func (lmexec *LMExec) deleteResourceGroupByID(params *lm.DeleteDeviceGroupByIDParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.DeleteDeviceGroupByID(params)
	}
}

// getResourceGroupList updates specified resource property.
func (lmexec *LMExec) getResourceGroupList(params *lm.GetDeviceGroupListParams) types.ExecRequest {
	return func() (interface{}, error) {
		return lmexec.LMClient.LM.GetDeviceGroupList(params)
	}
}
