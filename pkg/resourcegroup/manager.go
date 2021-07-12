package resourcegroup

import (
	"fmt"
	"net/http"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resourcecache"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup/dgbuilder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

// Manager implements types.ResourceManager
type Manager struct {
	*dgbuilder.Builder
	*types.LMRequester
	ResourceCache *resourcecache.ResourceCache
}

// CreateResourceGroupTree create resource group and its child tree
// nolint: cyclop
func (m *Manager) CreateResourceGroupTree(lctx *lmctx.LMContext, tree *types.ResourceGroupTree, update bool) error {
	if tree.DontCreate {
		return nil
	}
	conf, err := config.GetConfig()
	if err != nil {
		return aerrors.ErrCacheMiss
	}

	clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
	log := lmlog.Logger(clctx)
	resourceGroup, err := util.BuildResourceGroup(clctx, nil, tree.Options...)
	if err != nil {
		return err
	}
	key := types.ResourceName{
		Name:     *resourceGroup.Name,
		Resource: enums.Namespaces,
	}
	var resourceGroupID int32
	if meta, ok := m.ResourceCache.Exists(lctx, key, fmt.Sprintf("%d", resourceGroup.ParentID), false); ok {
		resourceGroupID = meta.LMID
		if update {
			err2 := m.updateResourceGroup(lctx, tree, update, log, resourceGroupID, meta, clctx, key, resourceGroup)
			if err2 != nil {
				return err2
			}
		}
	} else {
		resourceGroupID, err = m.createResourceGroup(log, clctx, resourceGroup)
		if err != nil {
			return err
		}
	}

	if tree.ChildGroups != nil {
		log.Debugf("Creating Child groups for resource group [%d]", resourceGroupID)
		for _, childTree := range tree.ChildGroups {
			childTree.Options = append(childTree.Options, m.Builder.ParentID(resourceGroupID))
			err := m.CreateResourceGroupTree(lctx, childTree, update)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) createResourceGroup(log *logrus.Entry, clctx *lmctx.LMContext, resourceGroup *models.DeviceGroup) (int32, error) {
	log.Debugf("Creating resource group")
	resp, err := m.addResourceGroup(clctx, resourceGroup)
	if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) != http.StatusBadRequest {
		return 0, fmt.Errorf("failed to create resourceGroup (%s) [parent ID: %d]: %w", *resourceGroup.Name, resourceGroup.ParentID, err)
	}
	if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusBadRequest {
		// TODO cache miss prom metric
		log.Warnf("seems resource group with same name is already present (%s) [parent ID: %d]: %s", *resourceGroup.Name, resourceGroup.ParentID, err)
		resp, err := m.getResourceGroupByName(clctx, resourceGroup.ParentID, *resourceGroup.Name)
		if err != nil {
			return -1, err
		}
		return resp.ID, nil
	}
	createdDeviceGroup := resp.(*lm.AddDeviceGroupOK).Payload
	return createdDeviceGroup.ID, nil
}

func (m *Manager) updateResourceGroup(lctx *lmctx.LMContext, tree *types.ResourceGroupTree, update bool, log *logrus.Entry, resourceGroupID int32, meta types.ResourceMeta, clctx *lmctx.LMContext, key types.ResourceName, resourceGroup *models.DeviceGroup) error {
	log.Infof("Updating existing resource group [%d]", resourceGroupID)
	resp, err := m.getResourceGroupByID(meta.LMID, clctx)
	if err != nil {
		if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
			// if group not found, invalidate cache and call again
			log.Errorf("resource group with ID (%d) does not exist, invalidating cache and performing create resourcegroup", meta.LMID)
			m.ResourceCache.Unset(lctx, key, fmt.Sprintf("%d", resourceGroup.ParentID))
			return m.CreateResourceGroupTree(lctx, tree, update)
		}
		return fmt.Errorf("failed to retrieve resource group for updation %d: %w", meta.LMID, err)
	}
	existingResourceGroup := resp.(*lm.GetDeviceGroupByIDOK).Payload
	existingResourceGroup, err = util.BuildResourceGroup(clctx, existingResourceGroup, tree.Options...)
	if err != nil {
		return fmt.Errorf("failed to modify resource group for updation %d: %w", meta.LMID, err)
	}
	err = m.updateResourceGroupByID(clctx, meta, existingResourceGroup)
	if err != nil {
		return fmt.Errorf("failed to update resource group %d: %w", meta.LMID, err)
	}
	log.Debugf("Updated resource group [%d]", resourceGroupID)
	return nil
}

// DeleteResourceGroup deletes a resource group with the specified resourceGroupID.
func (m *Manager) DeleteResourceGroup(lctx *lmctx.LMContext, rt enums.ResourceType, id int32, deleteIfEmpty bool) error {
	if deleteIfEmpty {
		group, err := m.GetResourceGroupByID(lctx, rt, id)
		if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
			m.UnsetLMIDInCache(lctx, rt, id)
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to retrieve resource group to check its number of resources: %w", err)
		}
		if group.NumOfHosts > 0 {
			return fmt.Errorf("%w: %d", aerrors.ErrResourceGroupIsNotEmpty, group.NumOfHosts)
		}
	}
	params := lm.NewDeleteDeviceGroupByIDParams()
	params.ID = id
	deleteChildren := true
	params.SetDeleteChildren(&deleteChildren)
	deleteHard := true
	params.SetDeleteHard(&deleteHard)
	command := m.DeleteResourceGroupByIDCommand(lctx, params)
	_, err := m.SendReceive(lctx, command)
	if err == nil || util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		m.UnsetLMIDInCache(lctx, rt, id)
	}
	return err
}

// GetResourceGroupByID deletes a resource group with the specified resourceGroupID.
func (m *Manager) GetResourceGroupByID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) (*models.DeviceGroup, error) {
	params := lm.NewGetDeviceGroupByIDParams()
	params.ID = id
	command := m.GetResourceGroupByIDCommand(lctx, params)
	resp, err := m.SendReceive(lctx, command)
	if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		m.UnsetLMIDInCache(lctx, rt, id)
	}
	if err != nil {
		return nil, err
	}
	return resp.(*lm.GetDeviceGroupByIDOK).Payload, err
}
