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
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/k8s-argus/pkg/watch/deployment"
	"github.com/logicmonitor/k8s-argus/pkg/watch/hpa"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

// InitSyncer implements the initial sync through logicmonitor API
type InitSyncer struct {
	DeviceManager *device.Manager
}

// InitSync implements the initial sync through logicmonitor API
func (i *InitSyncer) InitSync(lctx *lmctx.LMContext, isRestart bool) {
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
		i.runSync(lctx, rest, isRestart)
	}
	log.Infof("Finished syncing the resource devices")
}

func (i *InitSyncer) runSync(lctx *lmctx.LMContext, rest *models.DeviceGroup, isRestart bool) {
	log := lmlog.Logger(lctx)
	wg := sync.WaitGroup{}
	wg.Add(len(rest.SubGroups))
	for _, subgroup := range rest.SubGroups {
		switch subgroup.Name {
		case constants.NodeDeviceGroupName:
			go func() {
				defer wg.Done()
				lctxNodes := lmlog.NewLMContextWith(log.WithFields(logrus.Fields{"res": "init-sync-nodes"}))
				i.initSyncNodes(lctxNodes, rest.ID, isRestart)
				log.Infof("Finish syncing %v", constants.NodeDeviceGroupName)
			}()
		// nolint: dupl
		case constants.PodDeviceGroupName:
			go func() {
				defer wg.Done()
				lctxPods := lmlog.NewLMContextWith(log.WithFields(logrus.Fields{"res": "init-sync-pods"}))
				i.initSyncNamespacedResource(lctxPods, constants.PodDeviceGroupName, rest.ID, isRestart)
				log.Infof("Finish syncing %v", constants.PodDeviceGroupName)
			}()
		case constants.ServiceDeviceGroupName:
			go func() {
				defer wg.Done()
				lctxServices := lmlog.NewLMContextWith(log.WithFields(logrus.Fields{"res": "init-sync-services"}))
				i.initSyncNamespacedResource(lctxServices, constants.ServiceDeviceGroupName, rest.ID, isRestart)
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
				i.initSyncNamespacedResource(lctxDeployments, constants.DeploymentDeviceGroupName, rest.ID, isRestart)
				log.Infof("Finish syncing %v", constants.DeploymentDeviceGroupName)
			}()
		case constants.HorizontalPodAutoscalerDeviceGroupName:
			go func() {
				defer wg.Done()
				if !permission.HasHorizontalPodAutoscalerPermissions() {
					log.Warnf("Resource HorizontalPodAutoscaler has no permissions, ignore sync")
					return
				}
				i.initSyncHPA(rest.ID, isRestart)
				log.Infof("Finish syncing %v", constants.HorizontalPodAutoscalerDeviceGroupName)
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

func (i *InitSyncer) initSyncNodes(lctx *lmctx.LMContext, parentGroupID int32, isRestart bool) {
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
	nodesMap, err := node.GetNodesMap(i.DeviceManager.K8sClient, i.DeviceManager.Config().ClusterName)
	if err != nil || nodesMap == nil {
		log.Warnf("Failed to get the nodes from k8s, err: %v", err)
		return
	}

	for _, subGroup := range rest.SubGroups {
		// all the node device will be added to the group "ALL", so we only need to check it
		if subGroup.Name != constants.AllNodeDeviceGroupName {
			continue
		}
		i.syncDevices(lctx, constants.NodeDeviceGroupName, nodesMap, subGroup, isRestart)
	}
}

func (i *InitSyncer) initSyncNamespacedResource(lctx *lmctx.LMContext, deviceType string, parentGroupID int32, isRestart bool) {
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
		//get pod/service/deployment info from k8s
		var deviceMap map[string]string
		clusterName := i.DeviceManager.Config().ClusterName

		if deviceType == constants.PodDeviceGroupName {
			deviceMap, err = pod.GetPodsMap(i.DeviceManager.K8sClient, subGroup.Name, clusterName)
		} else if deviceType == constants.ServiceDeviceGroupName {
			deviceMap, err = service.GetServicesMap(lctx, i.DeviceManager.K8sClient, subGroup.Name, clusterName)
		} else if deviceType == constants.DeploymentDeviceGroupName {
			deviceMap, err = deployment.GetDeploymentsMap(lctx, i.DeviceManager.K8sClient, subGroup.Name, clusterName)
		} else {
			return
		}
		if err != nil || deviceMap == nil {
			log.Warnf("Failed to get the %s from k8s, namespace: %v, err: %v", deviceType, subGroup.Name, err)
			continue
		}

		// get and check all the devices in the group
		i.syncDevices(lctx, deviceType, deviceMap, subGroup, isRestart)
	}
}

func (i *InitSyncer) syncDevices(lctx *lmctx.LMContext, resourceType string, resourcesMap map[string]string, subGroup *models.DeviceGroupData, isRestart bool) {
	log := lmlog.Logger(lctx)
	if len(resourcesMap) == 0 {
		log.Debugf("Ignoring sub group %v for synchronization", subGroup.FullPath)
		return
	}

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
		autoClusterName := util.GetPropertyValue(device, constants.K8sClusterNamePropertyKey)
		if autoClusterName != i.DeviceManager.Config().ClusterName {
			log.Infof("Ignore the device (%v) which does not have property %v:%v",
				*device.DisplayName, constants.K8sClusterNamePropertyKey, i.DeviceManager.Config().ClusterName)
			continue
		}

		// ignore the device if it is moved in _deleted group
		if util.GetPropertyValue(device, constants.K8sResourceDeletedOnPropertyKey) != "" {
			log.Debugf("Ignore the device (%v) for synchronization as it is moved in _deleted group", *device.DisplayName)
			continue
		}

		// the displayName may be renamed, we should use the complete displayName for comparison.
		fullDisplayName := util.GetFullDisplayName(device, resourceType, autoClusterName)
		_, exist := resourcesMap[fullDisplayName]
		if !exist {
			log.Infof("Delete the non-exist %v device: %v", resourceType, *device.DisplayName)
			err := i.DeviceManager.DeleteByID(lctx, strings.ToLower(resourceType), device.ID)
			if err != nil {
				log.Warnf("Failed to delete the device: %v", *device.DisplayName)
			}
			continue
		}
		// Rename devices as per config parameters only on Argus restart.
		if isRestart {
			i.renameDeviceToDesiredName(lctx, device, resourceType)
		}
	}
}

func (i *InitSyncer) renameDeviceToDesiredName(lctx *lmctx.LMContext, device *models.Device, resourceType string) {
	log := lmlog.Logger(lctx)
	// get name and namespace prop values from devices
	autoName := util.GetPropertyValue(device, constants.K8sDeviceNamePropertyKey)
	namespace := util.GetPropertyValue(device, constants.K8sDeviceNamespacePropertyKey)
	desiredDisplayName := i.DeviceManager.GetDesiredDisplayName(autoName, namespace, resourceType)

	if i.DeviceManager.Config().FullDisplayNameIncludeClusterName || *device.DisplayName != desiredDisplayName {
		log.Infof("Renaming existing %v device: %v to new name %s", resourceType, *device.DisplayName, desiredDisplayName)
		err := i.DeviceManager.RenameAndUpdateDevice(lctx, strings.ToLower(resourceType), device, desiredDisplayName)
		if err != nil {
			log.Errorf("Failed to rename the existing device %s", *device.DisplayName)
			return
		}
	}
}

// RunPeriodicSync runs synchronization periodically.
func (i *InitSyncer) RunPeriodicSync(syncTime time.Duration) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "periodic-sync"}))
	go func() {
		for {
			time.Sleep(syncTime)
			i.InitSync(lctx, false)
		}
	}()
}

func (i *InitSyncer) initSyncHPA(parentGroupID int32, isRestart bool) {

	deviceType := "HorizontalPodAutoscalers"
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "init-sync-hpa"}))
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
		//get hpa info from k8s
		var deviceMap map[string]string

		deviceMap, err = hpa.GetHorizontalPodAutoscalersMap(lctx, i.DeviceManager.K8sClient, subGroup.Name, i.DeviceManager.Config().ClusterName)

		if err != nil || deviceMap == nil {
			log.Warnf("Failed to get the %s from k8s, namespace: %v, err: %v", deviceType, subGroup.Name, err)
			continue
		}

		// get and check all the devices in the group
		i.syncDevices(lctx, deviceType, deviceMap, subGroup, isRestart)
	}
}
