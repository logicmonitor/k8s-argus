package device

import (
	"errors"
	"net/http"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// AddFunc returns func
func (m *Manager) AddFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...types.DeviceOption) (*models.Device, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.DeviceOption) (*models.Device, error) {
		log := lmlog.Logger(lctx)
		d, err := m.Add(lctx, rt, obj, options...)
		if err != nil {
			var _t0 *types.DeviceExists
			if ok := errors.Is(err, _t0); !ok {
				log.Errorf("Failed to add resource: %s", err)

				return nil, err
			}
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
func (m *Manager) UpdateFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ...types.DeviceOption) (*models.Device, error) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}, options ...types.DeviceOption) (*models.Device, error) {
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
func (m *Manager) DeleteFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...types.DeviceOption) error {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.DeviceOption) error {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("Failed to get config")
			return err
		}
		if conf.DeleteDevices {
			err := m.Delete(lctx, rt, obj, options...)
			if err != nil {
				if util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
					log.Infof("Device already does not exist: %s", err)

					return nil
				}
				log.Errorf("Failed to delete resource: %s", err)

				return err
			}
			log.Infof("Deleted device")
		} else {
			err := m.MarkDeleted(lctx, rt, obj, options...)
			if err != nil {
				log.Errorf("Failed to move resource: %s", err)

				return err
			}
			log.Infof("Moved device")
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
	nodeLabels := util.GetLabelsByPrefix(constants.LabelNodeRole, objectMeta.Labels)
	if len(nodeLabels) == 0 {
		log.Debugf("No Role Labels found")

		return
	}
	rn := cache.ResourceName{
		Name:     constants.NodeDeviceGroupName,
		Resource: enums.Namespaces,
	}
	cacheMetaList, _ := m.ResourceCache.Get(lctx, rn)
	parentID := int32(0)
	for _, cacheMeta := range cacheMetaList {
		if cacheMeta.Container == util.ClusterGroupName(conf.ClusterName) {
			parentID = cacheMeta.LMID
		}
	}
	if parentID == 0 {
		log.Errorf("No \"Nodes\" group found in cache to set parent ID for NodeRoleGroups to put under")

		return
	}
	appliesToBuilder := devicegroup.NewAppliesToBuilder().Auto("clustername").Equals(conf.ClusterName)
	devicegroup.NewAppliesToBuilder().And().Auto("clustername").Equals(conf.ClusterName)
	for k := range nodeLabels {
		role := strings.ReplaceAll(k, constants.LabelNodeRole, "")
		_, ok := m.ResourceCache.Exists(lctx, cache.ResourceName{
			Name:     role,
			Resource: enums.Namespaces,
		}, "Nodes")
		if ok {
			log.Infof("Devicegroup for role %s, already exists", role)

			continue
		}
		appliesToBuilder2 := *appliesToBuilder.(*devicegroup.AppliesToBuilderImpl)
		deleteAppliesToBuilder2 := *appliesToBuilder.(*devicegroup.AppliesToBuilderImpl)
		// nolint: exhaustivestruct
		opts := &devicegroup.Options{
			// TODO: Take parent id from cache
			ParentID:              parentID,
			Name:                  role,
			DisableAlerting:       conf.DisableAlerting,
			AppliesTo:             appliesToBuilder2.And().Exists(constants.LabelCustomPropertyPrefix + k).And().HasCategory(rt.GetCategory()),
			Client:                m.LMClient,
			DeleteDevices:         conf.DeleteDevices,
			AppliesToDeletedGroup: deleteAppliesToBuilder2.And().Exists(constants.LabelCustomPropertyPrefix + k).And().HasCategory(rt.GetDeletedCategory()),
			CustomProperties:      devicegroup.NewPropertyBuilder(),
		}

		log.Debugf("node role group options : %v", opts)

		_, err := devicegroup.Create(lctx, opts, m.ResourceCache)
		if err != nil {
			log.Errorf("Failed to add device group for node role to %q: %v", role, err)
		}
		log.Infof("Added node role group: %s", role)
		// TODO: Update Cache with device group data
	}
}
