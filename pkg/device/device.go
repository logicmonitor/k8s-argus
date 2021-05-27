package device

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Manager implements types.DeviceManager
type Manager struct {
	*types.Base
	*builder.Builder
	types.LMExecutor
	types.LMFacade
	ResourceCache *devicecache.ResourceCache
}

// GetResourceCache get cache
func (m *Manager) GetResourceCache() types.ResourceCache {
	return m.ResourceCache
}

// FindByDisplayName implements types.DeviceManager.
func (m *Manager) FindByDisplayName(lctx *lmctx.LMContext, resource enums.ResourceType, name string) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	filter := fmt.Sprintf("displayName:\"%s\"", name)
	params := lm.NewGetDeviceListParams()
	params.SetFilter(&filter)
	cmd := &types.HTTPCommand{
		IsGlobal: false,
		Command: &types.Command{ // nolint: exhaustivestruct
			ExecFun: m.GetDeviceList(params),
			LMCtx:   lctx,
		},
		Method:   "GET",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.GetDeviceListErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	// restResponse, err := m.LMClient.LM.GetDeviceList(params)
	if err != nil {
		return nil, fmt.Errorf("get device list api failed: %w", err)
	}
	resp := restResponse.(*lm.GetDeviceListOK)
	log.Debugf("%#v", resp)
	if resp.Payload.Total == 1 {
		return resp.Payload.Items[0], nil
	}

	return nil, nil
}

// GetAndStoreAll get
// nolint: gocyclo
func (m *Manager) GetAndStoreAll(rt enums.ResourceType) []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)

	switch rt {
	case enums.Pods:
		result = m.getPodList()
	case enums.Deployments:
		result = m.getDeploymentList()
	case enums.Services:
		result = m.getServiceList()
	case enums.Nodes:
		result = m.getNodeList()
	case enums.Hpas:
		result = m.getHPAList()
	case enums.ETCD, enums.Namespaces, enums.Unknown:
		return result
	}

	return result
}

func (m *Manager) getPodList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	podList, err := m.Base.K8sClient.CoreV1().Pods("").List(constants.DefaultListOptions)

	if err != nil || podList == nil {
		return result
	}
	for _, i := range podList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func (m *Manager) getServiceList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	serviceList, err := m.Base.K8sClient.CoreV1().Services("").List(constants.DefaultListOptions)

	if err != nil || serviceList == nil {
		return result
	}
	for _, i := range serviceList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func (m *Manager) getDeploymentList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	deploymentList, err := m.Base.K8sClient.AppsV1().Deployments("").List(constants.DefaultListOptions)

	if err != nil || deploymentList == nil {
		return result
	}
	for _, i := range deploymentList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func (m *Manager) getNodeList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	nodeList, err := m.Base.K8sClient.CoreV1().Nodes().List(constants.DefaultListOptions)

	if err != nil || nodeList == nil {
		return result
	}
	for _, i := range nodeList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func (m *Manager) getHPAList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	horizontalPodAutoscalerList, err := m.Base.K8sClient.AutoscalingV2beta2().HorizontalPodAutoscalers("").List(constants.DefaultListOptions)

	if err != nil || horizontalPodAutoscalerList == nil {
		return result
	}
	for _, i := range horizontalPodAutoscalerList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

// GetAllK8SResources get all k8s resources present in cluster
func (m *Manager) GetAllK8SResources() *devicecache.Store {
	tmpStore := devicecache.NewStore()
	conf, err := config.GetConfig()
	if err != nil {
		logrus.Errorf("Failed to read config")

		return nil
	}
	for _, rt := range enums.ALLResourceTypes {
		for _, meta := range m.GetAndStoreAll(rt) {
			displayName := util.GetDisplayNameNew(rt, &meta, conf) //nolint:gosec
			tmpStore.Set(cache.ResourceName{
				Name:     meta.Name,
				Resource: rt,
			}, cache.ResourceMeta{ // nolint: exhaustivestruct
				Container:   meta.Namespace,
				Labels:      meta.Labels,
				DisplayName: displayName,
			})
		}
	}

	return tmpStore
}

// ModifyToUnique modify
func (m *Manager) ModifyToUnique(lctx *lmctx.LMContext, resource enums.ResourceType, device *models.Device, obj interface{}) []types.DeviceOption {
	options := make([]types.DeviceOption, 0)
	conf, err := config.GetConfig()
	if err != nil {
		return options
	}
	log := lmlog.Logger(lctx)
	fullDisplayName := util.GetFullDisplayName(device, resource, conf.ClusterName)
	if fullDisplayName == *device.DisplayName {
		log.Info("Removing conflicts category from resource")
		options = append(options, m.SystemCategory(resource.GetConflictsCategory(), enums.Delete))
	} else {
		log.Info("Adding conflicts category on resource and renaming to full name")
		options = append(options, []types.DeviceOption{
			m.SystemCategory(resource.GetConflictsCategory(), enums.Add),
			m.DisplayName(fullDisplayName),
		}...)
	}

	return options
}
