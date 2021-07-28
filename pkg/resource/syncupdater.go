package resource

import (
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
)

func (m *Manager) UpdateResourceByID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32, options ...types.ResourceOption) (*models.Device, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	log := lmlog.Logger(lctx)
	resourceObj, err := m.FetchResource(lctx, rt, id)
	if err != nil {
		return nil, err
	}
	modifiedResource, err := util.BuildResource(lctx, conf, resourceObj, options...)
	if err != nil {
		log.Errorf("Failed to build modified resource")
		return nil, err
	}
	return m.UpdateAndReplaceResource(lctx, rt, id, modifiedResource)
}

func (m *Manager) DeleteResourceByID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) error {
	return m.DeleteByID(lctx, rt, id)
}
