package resource

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup/dgbuilder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// AddFunc returns func
func (m *Manager) AddFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...types.ResourceOption) (*models.Device, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.ResourceOption) (*models.Device, error) {
		log := lmlog.Logger(lctx)
		d, err := m.Add(lctx, rt, obj, options...)
		if err != nil {
			if errors.Is(err, aerrors.ErrResourceExists) {
				log.Warnf("%s", err)

				return nil, err
			}
			return nil, err
		}

		log.Infof("Added resource")
		// Special handling just for Nodes, do not add anymore if for other resources
		if rt == enums.Nodes {
			m.createNodeRoleGroups(lctx, rt, obj)
		}
		return d, err
	}
}

// UpdateFunc returns func
func (m *Manager) UpdateFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ...types.ResourceOption) (*models.Device, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, options ...types.ResourceOption) (*models.Device, error) {
		log := lmlog.Logger(lctx)
		d, err := m.Update(lctx, rt, oldObj, newObj, options...)
		if err != nil {
			log.Errorf("Failed to update resource: %s", err)
			return nil, err
		}

		log.Infof("Updated resource")
		// Special handling just for Nodes, do not add anymore if for other resources
		if rt == enums.Nodes {
			m.createNodeRoleGroups(lctx, rt, newObj)
		}
		return d, err
	}
}

// DeleteFunc returns function
func (m *Manager) DeleteFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...types.ResourceOption) error {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.ResourceOption) error {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("Failed to get config")
			return err
		}
		resource, err := util.BuildResource(lctx, conf, nil, options...)
		if err != nil {
			return err
		}
		if conf.DeleteResources &&
			!util.IsArgusPod(lctx, rt, resource) {
			err := m.Delete(lctx, rt, obj, options...)
			if err != nil {
				if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
					log.Infof("Resource already does not exist: %s", err)

					return nil
				}
				log.Errorf("Failed to delete resource: %s", err)

				return err
			}
			log.Infof("Deleted resource")
		} else {
			err := m.MarkDeleted(lctx, rt, obj, options...)
			if err != nil {
				log.Errorf("Failed to move resource: %s", err)

				return err
			}
			log.Infof("Moved resource")
		}
		return nil
	}
}

func (m *Manager) createNodeRoleGroups(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	log := lmlog.Logger(lctx)

	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return
	}

	objectMeta := rt.ObjectMeta(obj)
	nodeRoleLabels := util.GetLabelsByPrefix(constants.LabelNodeRole, objectMeta.Labels)
	if len(nodeRoleLabels) == 0 {
		log.Debugf("No Role Labels found")

		return
	}
	rn := types.ResourceName{
		Name:     constants.NodeResourceGroupName,
		Resource: enums.Namespaces,
	}
	cacheMetaList, _ := m.ResourceCache.Get(lctx, rn)
	parentID := int32(0)
	for _, cacheMeta := range cacheMetaList {
		if cacheMeta.Container == fmt.Sprintf("%d", util.GetClusterGroupID(lctx, m.LMRequester)) {
			parentID = cacheMeta.LMID
		}
	}
	if parentID == 0 {
		log.Errorf("No \"Nodes\" group found in cache to set parent ID for NodeRoleGroups to put under")

		return
	}
	for k := range nodeRoleLabels {
		role := strings.ReplaceAll(k, constants.LabelNodeRole, "")

		resourceTree := &types.ResourceGroupTree{
			Options: []types.ResourceGroupOption{
				m.ResourceGroupManager.GroupName(role),
				m.ResourceGroupManager.ParentID(parentID),
				m.ResourceGroupManager.DisableAlerting(conf.DisableAlerting),
				m.ResourceGroupManager.AppliesTo(
					dgbuilder.NewAppliesToBuilder().
						Auto("clustername").Equals(conf.ClusterName).And().
						Exists(constants.LabelCustomPropertyPrefix + k).And().
						HasCategory(rt.GetCategory()),
				),
			},
		}
		err := m.CreateResourceGroupTree(lctx, resourceTree, false)
		if err != nil {
			log.Errorf("Failed to add resource group for node role to %q: %v", role, err)

			return
		}
		log.Infof("Added node role group: %s", role)
	}
}
