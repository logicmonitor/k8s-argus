package types

import (
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
)

type ResourceExecutor interface {
	AddResourceCommand(*lmctx.LMContext, *lm.AddDeviceParams) *WorkerCommand

	UpdateResourceCommand(*lmctx.LMContext, *lm.UpdateDeviceParams) *WorkerCommand

	GetResourceByIDCommand(*lmctx.LMContext, *lm.GetDeviceByIDParams) *WorkerCommand

	UpdateResourcePropertyByNameCommand(*lmctx.LMContext, *lm.UpdateDevicePropertyByNameParams) *WorkerCommand

	GetResourceListCommand(*lmctx.LMContext, *lm.GetDeviceListParams) *WorkerCommand

	PatchResourceCommand(*lmctx.LMContext, *lm.PatchDeviceParams) *WorkerCommand

	DeleteResourceByIDCommand(*lmctx.LMContext, *lm.DeleteDeviceByIDParams) *WorkerCommand
}

type ResourceGroupExecutor interface {
	GetImmediateResourceListByResourceGroupIDCommand(*lmctx.LMContext, *lm.GetImmediateDeviceListByDeviceGroupIDParams) *WorkerCommand

	AddResourceGroupCommand(*lmctx.LMContext, *lm.AddDeviceGroupParams) *WorkerCommand

	UpdateResourceGroupByIDCommand(*lmctx.LMContext, *lm.UpdateDeviceGroupByIDParams) *WorkerCommand

	AddResourceGroupPropertyCommand(*lmctx.LMContext, *lm.AddDeviceGroupPropertyParams) *WorkerCommand

	UpdateResourceGroupPropertyByNameCommand(*lmctx.LMContext, *lm.UpdateDeviceGroupPropertyByNameParams) *WorkerCommand

	DeleteResourceGroupPropertyByNameCommand(*lmctx.LMContext, *lm.DeleteDeviceGroupPropertyByNameParams) *WorkerCommand

	GetResourceGroupByIDCommand(*lmctx.LMContext, *lm.GetDeviceGroupByIDParams) *WorkerCommand

	GetResourceGroupListCommand(*lmctx.LMContext, *lm.GetDeviceGroupListParams) *WorkerCommand

	DeleteResourceGroupByIDCommand(*lmctx.LMContext, *lm.DeleteDeviceGroupByIDParams) *WorkerCommand
}

// LMExecutor all lm rest apis used
type LMExecutor interface {
	ResourceExecutor
	ResourceGroupExecutor
}
