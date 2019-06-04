package sync

import (
	"sync"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/watch/deployment"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/logicmonitor/lm-sdk-go/models"
	log "github.com/sirupsen/logrus"
)

// InitSyncer implements the initial sync through logicmonitor API
type InitSyncer struct {
	DeviceManager *device.Manager
}

// InitSync implements the initial sync through logicmonitor API
func (i *InitSyncer) InitSync() {
	log.Infof("Start to sync the resource devices")
	clusterName := i.DeviceManager.Base.Config.ClusterName
	// get the cluster info
	parentGroupID := i.DeviceManager.Config().ClusterGroupID
	groupName := constants.ClusterDeviceGroupPrefix + clusterName
	rest, err := devicegroup.Find(parentGroupID, groupName, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Infof("Failed to get the cluster group: %v, parentID: %v", groupName, parentGroupID)
		return
	}

	// get the node, pod, service, deployment info
	if rest.SubGroups != nil && len(rest.SubGroups) != 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(rest.SubGroups))
		for _, subgroup := range rest.SubGroups {
			switch subgroup.Name {
			case constants.NodeDeviceGroupName:
				go func() {
					defer wg.Done()
					i.initSyncNodes(rest.ID)
					log.Infof("Finish syncing %v", constants.NodeDeviceGroupName)
				}()
			case constants.PodDeviceGroupName:
				go func() {
					defer wg.Done()
					i.initSyncPodsOrServicesOrDeploys(constants.PodDeviceGroupName, rest.ID)
					log.Infof("Finish syncing %v", constants.PodDeviceGroupName)
				}()
			case constants.ServiceDeviceGroupName:
				go func() {
					defer wg.Done()
					i.initSyncPodsOrServicesOrDeploys(constants.ServiceDeviceGroupName, rest.ID)
					log.Infof("Finish syncing %v", constants.ServiceDeviceGroupName)
				}()
			case constants.DeploymentDeviceGroupName:
				go func() {
					defer wg.Done()
					i.initSyncPodsOrServicesOrDeploys(constants.DeploymentDeviceGroupName, rest.ID)
					log.Infof("Finish syncing %v", constants.DeploymentDeviceGroupName)
				}()
			default:
				func() {
					defer wg.Done()
					log.Infof("Unsupported group to sync, ignore it: %v", subgroup.Name)
				}()

			}
		}

		// wait the init sync processes finishing
		wg.Wait()
	}
	log.Infof("Finished syncing the resource devices")
}

func (i *InitSyncer) initSyncNodes(parentGroupID int32) {
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
		i.syncDevices(constants.NodeDeviceGroupName, nodesMap, subGroup)
	}
}

func (i *InitSyncer) initSyncPodsOrServicesOrDeploys(deviceType string, parentGroupID int32) {
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
			deviceMap, err = service.GetServicesMap(i.DeviceManager.K8sClient, subGroup.Name)
		} else if deviceType == constants.DeploymentDeviceGroupName {
			deviceMap, err = deployment.GetDeploymentsMap(i.DeviceManager.K8sClient, subGroup.Name)
		} else {
			return
		}
		if err != nil || deviceMap == nil {
			log.Warnf("Failed to get the %s from k8s, namespace: %v, err: %v", deviceType, subGroup.Name, err)
			continue
		}

		// get and check all the devices in the group
		i.syncDevices(deviceType, deviceMap, subGroup)
	}
}

func (i *InitSyncer) syncDevices(resourceType string, resourcesMap map[string]string, subGroup *models.DeviceGroupData) {
	devices, err := i.DeviceManager.GetListByGroupID(subGroup.ID)
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
		cps := device.CustomProperties
		autoClusterName := ""
		for _, cp := range cps {
			if *cp.Name == constants.K8sClusterNamePropertyKey {
				autoClusterName = *cp.Value
				break
			}
		}
		if autoClusterName != i.DeviceManager.Config().ClusterName {
			log.Infof("Ignore the device (%v) which does not have property %v:%v",
				*device.DisplayName, constants.K8sClusterNamePropertyKey, i.DeviceManager.Config().ClusterName)
			continue
		}

		_, exist := resourcesMap[*device.DisplayName]
		if !exist {
			log.Infof("Delete the non-exist %v device: %v", resourceType, *device.DisplayName)
			err := i.DeviceManager.DeleteByID(device.ID)
			if err != nil {
				log.Warnf("Failed to delete the device: %v", *device.DisplayName)
			}
		}
	}
}
