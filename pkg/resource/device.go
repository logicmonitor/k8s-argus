package resource

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resource/builder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// Manager implements types.ResourceManager
type Manager struct {
	*builder.Builder
	types.ResourceGroupManager
	*types.LMRequester
	types.ResourceCache
}

// GetResourceCache get cache
func (m *Manager) GetResourceCache() types.ResourceCache {
	return m.ResourceCache
}

// FindByDisplayName implements types.ResourceManager.
func (m *Manager) FindByDisplayName(lctx *lmctx.LMContext, resource enums.ResourceType, name string) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	filter := fmt.Sprintf("displayName:\"%s\"", name)
	params := lm.NewGetDeviceListParams()
	params.SetFilter(&filter)
	command := m.GetResourceListCommand(lctx, params)
	restResponse, err := m.LMFacade.SendReceive(lctx, command)
	if err != nil {
		return nil, fmt.Errorf("get resource list api failed: %w", err)
	}
	resp := restResponse.(*lm.GetDeviceListOK)
	log.Debugf("%#v", resp)
	if resp.Payload.Total == 1 {
		return resp.Payload.Items[0], nil
	}

	return nil, nil
}
