package sync

import (
	"strings"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/watch/deployment"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/sirupsen/logrus"
	"github.com/vkumbhar94/lm-sdk-go/models"
)

// InitSyncer implements the initial sync through logicmonitor API
type InitSyncer struct {
	DeviceManager *device.Manager
}

// InitSync implements the initial sync through logicmonitor API
func (i *InitSyncer) InitSync(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	log.Infof("Start to sync the resource devices")
	clusterName := i.DeviceManager.Base.Config.ClusterName
	// get the cluster info
	parentGroupID := i.DeviceManager.Config().ClusterGroupID
	groupName := constants.ClusterDeviceGroupPrefix + clusterName
	rest, err := devicegroup.Find(parentGroupID, groupName, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Infof("Failed to get the cluster group: %v, parentID: %v", groupName, parentGroupID)
	}
	if rest == nil {
		return
	}
	// get the node, pod, service, deployment info
	if rest.SubGroups != nil && len(rest.SubGroups) != 0 {
		i.runSync(lctx, rest)
	}
	log.Infof("Finished syncing the resource devices")
}

func (i *InitSyncer) runSync(lctx *lmctx.LMContext, rest *models.DeviceGroup) {
	log := lmlog.Logger(lctx)
	wg := sync.WaitGroup{}
	wg.Add(len(rest.SubGroups))
	for _, subgroup := range rest.SubGroups {
		switch subgroup.Name {
		case constants.NodeDeviceGroupName:
			go func() {
				defer wg.Done()
				lctxNodes := lmlog.NewLMContextWith(log.WithFields(logrus.Fields{"res": "init-sync-nodes"}))
				i.initSyncNodes(lctxNodes, rest.ID)
				log.Infof("Finish syncing %v", constants.NodeDeviceGroupName)
			}()
		// nolint: dupl
		case constants.PodDeviceGroupName:
			go func() {
				defer wg.Done()
				lctxPods := lmlog.NewLMContextWith(log.WithFields(logrus.Fields{"res": "init-sync-pods"}))
				i.initSyncNamespacedResource(lctxPods, constants.PodDeviceGroupName, rest.ID)
				log.Infof("Finish syncing %v", constants.PodDeviceGroupName)
			}()
		case constants.ServiceDeviceGroupName:
			go func() {
				defer wg.Done()
				lctxServices := lmlog.NewLMContextWith(log.WithFields(logrus.Fields{"res": "init-sync-services"}))
				i.initSyncNamespacedResource(lctxServices, constants.ServiceDeviceGroupName, rest.ID)
				log.Infof("Finish syncing %v", constants.ServiceDeviceGroupName)
			}()
		case constants.DeploymentDeviceGroupName:
			go func() {
				defer wg.Done()
				lctxDeployments := lmlog.NewLMContextWith(log.WithFields(logrus.Fields{"res": "init-sync-deployments"}))
				if !permission.HasDeploymentPermissions() {
					log.Warnf("Resource deployments has no permissions, ignore sync")
					return
				}
				i.initSyncNamespacedResource(lctxDeployments, constants.DeploymentDeviceGroupName, rest.ID)
				log.Infof("Finish syncing %v", constants.DeploymentDeviceGroupName)
			}()
		default:
			func() {
				defer wg.Done()
				log.Infof("Unsupported group to sync, ignore it: %v", subgroup.Name)
			}()

		}
	}
	log.Debugf("Waiting to complete sync")
	// wait the init sync processes finishing
	wg.Wait()
	log.Debugf("Completed sync")
}

func (i *InitSyncer) initSyncNodes(lctx *lmctx.LMContext, parentGroupID int32) {
	log := lmlog.Logger(lctx)
	rest, err := devicegroup.Find(parentGroupID, constants.NodeDeviceGroupName, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Warnf("Failed to get the node group")
		return
	}
	if rest.SubGroups == nil {
		return
	}

	//get node info from k8s
	nodesMap, err := node.GetNodesMap(i.DeviceManager.K8sClient)
	if err != nil || nodesMap == nil {
		log.Warnf("Failed to get the nodes from k8s, err: %v", err)
		return
	}

	for _, subGroup := range rest.SubGroups {
		// all the node device will be added to the group "ALL", so we only need to check it
		if subGroup.Name != constants.AllNodeDeviceGroupName {
			continue
		}
		i.syncDevices(lctx, constants.NodeDeviceGroupName, nodesMap, subGroup)
	}
}

func (i *InitSyncer) initSyncNamespacedResource(lctx *lmctx.LMContext, deviceType string, parentGroupID int32) {
	log := lmlog.Logger(lctx)
	rest, err := devicegroup.Find(parentGroupID, deviceType, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Warnf("Failed to get the %s group", deviceType)
		return
	}
	if rest.SubGroups == nil {
		return
	}

	// loop every namespace
	for _, subGroup := range rest.SubGroups {
		//get pod/service info from k8s
		var deviceMap map[string]string
		if deviceType == constants.PodDeviceGroupName {
			deviceMap, err = pod.GetPodsMap(i.DeviceManager.K8sClient, subGroup.Name)
		} else if deviceType == constants.ServiceDeviceGroupName {
			deviceMap, err = service.GetServicesMap(lctx, i.DeviceManager.K8sClient, subGroup.Name)
		} else if deviceType == constants.DeploymentDeviceGroupName {
			deviceMap, err = deployment.GetDeploymentsMap(lctx, i.DeviceManager.K8sClient, subGroup.Name)
		} else {
			return
		}
		if err != nil || deviceMap == nil {
			log.Warnf("Failed to get the %s from k8s, namespace: %v, err: %v", deviceType, subGroup.Name, err)
			continue
		}

		// get and check all the devices in the group
		i.syncDevices(lctx, deviceType, deviceMap, subGroup)
	}
}

func (i *InitSyncer) syncDevices(lctx *lmctx.LMContext, resourceType string, resourcesMap map[string]string, subGroup *models.DeviceGroupData) {
	log := lmlog.Logger(lctx)
	devices, err := i.DeviceManager.GetListByGroupID(lctx, strings.ToLower(resourceType), subGroup.ID)
	if err != nil {
		log.Warnf("Failed to get the devices in the group: %v", subGroup.FullPath)
		return
	}
	if len(devices) == 0 {
		log.Warnf("There is no device in the group: %v", subGroup.FullPath)
		return
	}
	for _, device := range devices {
		// the "auto.clustername" property checking is used to prevent unexpected deletion of the normal non-k8s device
		// which may be assigned to the cluster group
		autoClusterName := i.DeviceManager.GetPropertyValue(device, constants.K8sClusterNamePropertyKey)
		if autoClusterName != i.DeviceManager.Config().ClusterName {
			log.Infof("Ignore the device (%v) which does not have property %v:%v",
				*device.DisplayName, constants.K8sClusterNamePropertyKey, i.DeviceManager.Config().ClusterName)
			continue
		}
		// the displayName may be renamed, we should use the constants.K8sResourceNamePropertyKey property value
		resourceName := i.DeviceManager.GetPropertyValue(device, constants.K8sResourceNamePropertyKey)
		// for compatibility, if resourceName is empty, use display name
		if resourceName == "" {
			resourceName = *device.DisplayName
		}
		_, exist := resourcesMap[resourceName]
		if !exist {
			log.Infof("Delete the non-exist %v device: %v", resourceType, *device.DisplayName)
			err := i.DeviceManager.DeleteByID(lctx, strings.ToLower(resourceType), device.ID)
			if err != nil {
				log.Warnf("Failed to delete the device: %v", *device.DisplayName)
			}
		}
	}
}

// RunPeriodicSync runs synchronization periodically.
func (i *InitSyncer) RunPeriodicSync(syncTime int) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "periodic-sync"}))
	go func() {
		for {
			time.Sleep(time.Duration(syncTime) * time.Minute)
			i.InitSync(lctx)
		}
	}()
}
