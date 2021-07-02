package lmexec

import (
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
)

func (lmexec *LMExec) AddResourceCommand(lctx *lmctx.LMContext, params *lm.AddDeviceParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.addResource(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/devices",
			Method:     types.HTTPPost,
		},
	}
}

func (lmexec *LMExec) UpdateResourceCommand(lctx *lmctx.LMContext, params *lm.UpdateDeviceParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.updateResource(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/devices/{id}",
			Method:     types.HTTPPut,
		},
	}
}

func (lmexec *LMExec) GetResourceByIDCommand(lctx *lmctx.LMContext, params *lm.GetDeviceByIDParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.getResourceByID(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/devices/{id}",
			Method:     types.HTTPGet,
		},
	}
}

func (lmexec *LMExec) GetResourceListCommand(lctx *lmctx.LMContext, params *lm.GetDeviceListParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.getResourceList(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/devices",
			Method:     types.HTTPGet,
		},
	}
}

func (lmexec *LMExec) PatchResourceCommand(lctx *lmctx.LMContext, params *lm.PatchDeviceParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.patchResource(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/devices/{id}",
			Method:     types.HTTPPatch,
		},
	}
}

func (lmexec *LMExec) DeleteResourceByIDCommand(lctx *lmctx.LMContext, params *lm.DeleteDeviceByIDParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.deleteResourceByID(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/devices/{id}",
			Method:     types.HTTPDelete,
		},
	}
}

func (lmexec *LMExec) UpdateResourcePropertyByNameCommand(lctx *lmctx.LMContext, params *lm.UpdateDevicePropertyByNameParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.updateResourcePropertyByName(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/devices/{resourceId}/properties/{name}",
			Method:     types.HTTPPut,
		},
	}
}

func (lmexec *LMExec) AddResourceGroupCommand(lctx *lmctx.LMContext, params *lm.AddDeviceGroupParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.addResourceGroup(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups",
			Method:     types.HTTPPost,
		},
	}
}

func (lmexec *LMExec) UpdateResourceGroupByIDCommand(lctx *lmctx.LMContext, params *lm.UpdateDeviceGroupByIDParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.updateResourceGroupByID(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups/{id}",
			Method:     types.HTTPPut,
		},
	}
}

func (lmexec *LMExec) AddResourceGroupPropertyCommand(lctx *lmctx.LMContext, params *lm.AddDeviceGroupPropertyParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.addResourceGroupProperty(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups/{gid}/properties",
			Method:     types.HTTPPost,
		},
	}
}

func (lmexec *LMExec) UpdateResourceGroupPropertyByNameCommand(lctx *lmctx.LMContext, params *lm.UpdateDeviceGroupPropertyByNameParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.updateResourceGroupPropertyByName(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups/{gid}/properties/{name}",
			Method:     types.HTTPPut,
		},
	}
}

func (lmexec *LMExec) DeleteResourceGroupPropertyByNameCommand(lctx *lmctx.LMContext, params *lm.DeleteDeviceGroupPropertyByNameParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.deleteResourceGroupPropertyByName(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups/{gid}/properties/{name}",
			Method:     types.HTTPDelete,
		},
	}
}

func (lmexec *LMExec) GetImmediateResourceListByResourceGroupIDCommand(lctx *lmctx.LMContext, params *lm.GetImmediateDeviceListByDeviceGroupIDParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.getImmediateResourceListByResourceGroupID(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups/{id}/devices",
			Method:     types.HTTPGet,
		},
	}
}

func (lmexec *LMExec) GetResourceGroupByIDCommand(lctx *lmctx.LMContext, params *lm.GetDeviceGroupByIDParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.getResourceGroupByID(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups/{id}",
			Method:     types.HTTPGet,
		},
	}
}

func (lmexec *LMExec) GetResourceGroupListCommand(lctx *lmctx.LMContext, params *lm.GetDeviceGroupListParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.getResourceGroupList(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups",
			Method:     types.HTTPGet,
		},
	}
}

func (lmexec *LMExec) DeleteResourceGroupByIDCommand(lctx *lmctx.LMContext, params *lm.DeleteDeviceGroupByIDParams) *types.WorkerCommand {
	return &types.WorkerCommand{
		ExecFunc: lmexec.deleteResourceGroupByID(params),
		Lctx:     lctx,
		APIInfo: types.APIInfo{
			URLPattern: "/device/groups/{id}",
			Method:     types.HTTPDelete,
		},
	}
}
