package sync

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InitSyncer implements the initial sync with Santaba
type InitSyncer struct {
	DeviceManager *device.Manager
}

// InitSync implements the initial sync with Santaba
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

	// get the node, pod, service info
	if rest.SubGroups != nil {
		c := make(chan string, 3)
		syncNum := 0
		for _, subgroup := range rest.SubGroups {
			switch subgroup.Name {
			case constants.NodeDeviceGroupName:
				go i.intSyncNodes(rest.Id, c)
				syncNum++
			case constants.PodDeviceGroupName:
				go i.initSyncPods(rest.Id, c)
				syncNum++
			case constants.ServiceDeviceGroupName:
				go i.initSyncServices(rest.Id, c)
				syncNum++
			default:
				log.Infof("Unsupported group to sync, ignore it: %v", subgroup.Name)
			}
		}

		for i := 0; i < syncNum; i++ {
			log.Infof("Finish syncing %v", <-c)
		}
	}
}

func (i *InitSyncer) intSyncNodes(parentGroupID int32, c chan string) {
	defer i.sendInfoToChan(constants.NodeDeviceGroupName, c)

	rest, err := devicegroup.Find(parentGroupID, constants.NodeDeviceGroupName, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Warnf("Failed to get the node group")
		return
	}
	if rest.SubGroups == nil {
		return
	}

	//get node info from k8s
	nodesMap := make(map[string]string)
	nodeList, err := i.DeviceManager.K8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil || nodeList == nil {
		log.Warnf("Failed to get the nodes from k8s")
		return
	}
	for _, nodeInfo := range nodeList.Items {
		nodesMap[nodeInfo.Name] = node.GetInternalAddress(nodeInfo.Status.Addresses).Address
	}

	for _, subGroup := range rest.SubGroups {
		// all the node device will be added to the group "ALL", so we only need to check it
		if subGroup.Name != constants.AllNodeDeviceGroupName {
			continue
		}
		i.syncDevices(constants.NodeDeviceGroupName, nodesMap, subGroup)
	}

}

func (i *InitSyncer) initSyncPods(parentGroupID int32, c chan string) {
	defer i.sendInfoToChan(constants.PodDeviceGroupName, c)

	rest, err := devicegroup.Find(parentGroupID, constants.PodDeviceGroupName, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Warnf("Failed to get the pod group")
		return
	}
	if rest.SubGroups == nil {
		return
	}

	// loop every namesplace
	for _, subGroup := range rest.SubGroups {
		//get pod info from k8s
		podsMap := make(map[string]string)
		podList, err := i.DeviceManager.K8sClient.CoreV1().Pods(subGroup.Name).List(metav1.ListOptions{})
		if err != nil || podList == nil {
			log.Warnf("Failed to get the pods from k8s")
			return
		}
		for _, podInfo := range podList.Items {
			// TODO: we should improve the value of the map to the ip of the pod when changing the name of the device to the ip
			podsMap[podInfo.Name] = podInfo.Name
		}

		// get and check all the devices in the group
		i.syncDevices(constants.PodDeviceGroupName, podsMap, subGroup)
	}
}

func (i *InitSyncer) initSyncServices(parentGroupID int32, c chan string) {
	defer i.sendInfoToChan(constants.ServiceDeviceGroupName, c)

	rest, err := devicegroup.Find(parentGroupID, constants.ServiceDeviceGroupName, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Warnf("Failed to get the pod group")
		return
	}
	if rest.SubGroups == nil {
		return
	}

	// loop every namesplace
	for _, subGroup := range rest.SubGroups {
		//get service info from k8s
		servicesMap := make(map[string]string)
		serviceList, err := i.DeviceManager.K8sClient.CoreV1().Services(subGroup.Name).List(metav1.ListOptions{})
		if err != nil || serviceList == nil {
			log.Warnf("Failed to get the services from k8s")
			return
		}
		for _, serviceInfo := range serviceList.Items {
			servicesMap[service.FmtServiceDisplayName(&serviceInfo)] = service.FmtServiceName(&serviceInfo)
		}

		// get and check all the devices in the group
		i.syncDevices(constants.ServiceDeviceGroupName, servicesMap, subGroup)
	}
}

func (i *InitSyncer) sendInfoToChan(info string, c chan string) {
	c <- info
}

func (i *InitSyncer) syncDevices(resourceType string, resourcesMap map[string]string, subGroup logicmonitor.GroupData) {
	devices, err := i.DeviceManager.GetListByGroupID(subGroup.Id)
	if err != nil || devices == nil {
		log.Warnf("Failed to get the devices in the group: %v", subGroup.FullPath)
		return
	}
	for _, device := range devices {
		name, exist := resourcesMap[device.DisplayName]
		if !exist || name != device.Name {
			log.Infof("Delete the non-exist %v device: %v", resourceType, device.DisplayName)
			err := i.DeviceManager.DeleteByID(device.Id)
			if err != nil {
				log.Warnf("Failed to delete the device: %v", device.DisplayName)
			}
		}
	}
}
