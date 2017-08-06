package tree

import (
	"fmt"
	"net/url"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
)

// DeviceTree manages the device tree representation of a Kubernetes cluster in LogicMonitor.
type DeviceTree struct {
	*types.Base
}

// CreateDeviceTree creates the Device tree that will represent the cluster in the LogicMonitor portal.
func (d *DeviceTree) CreateDeviceTree() (map[string]int32, error) {
	deviceGroups := make(map[string]int32)
	clusterDeviceGroup, err := d.createClusterDeviceGroup()
	if err != nil {
		return nil, err
	}
	log.Infof("Using cluster device group with id %d", clusterDeviceGroup.Id)

	serviceDeviceGroup, err := d.createServiceDeviceGroup(clusterDeviceGroup)
	if err != nil {
		return nil, err
	}
	deviceGroups["services"] = serviceDeviceGroup.Id
	_, err = d.createServiceDeletedDeviceGroup(serviceDeviceGroup)
	if err != nil {
		return nil, err
	}
	log.Infof("Using service device group with id %d", serviceDeviceGroup.Id)

	etcdDeviceGroup, err := d.createEtcdDeviceGroup(clusterDeviceGroup)
	if err != nil {
		return nil, err
	}
	_, err = d.createEtcdDeletedDeviceGroup(etcdDeviceGroup)
	if err != nil {
		return nil, err
	}
	log.Infof("Using etcd device group with id %d", etcdDeviceGroup.Id)

	nodeDeviceGroup, err := d.createNodeDeviceGroup(clusterDeviceGroup)
	if err != nil {
		return nil, err
	}
	_, err = d.createNodeDeletedDeviceGroup(nodeDeviceGroup)
	if err != nil {
		return nil, err
	}
	log.Infof("Using node device group with id %d", nodeDeviceGroup.Id)

	podDeviceGroup, err := d.createPodDeviceGroup(clusterDeviceGroup)
	if err != nil {
		return nil, err
	}
	deviceGroups["pods"] = podDeviceGroup.Id
	_, err = d.createPodDeletedDeviceGroup(podDeviceGroup)
	if err != nil {
		return nil, err
	}
	log.Infof("Using pod device group with id %d", podDeviceGroup.Id)

	return deviceGroups, nil
}

func (d *DeviceTree) findDeviceGroup(parentID int32, name string) (deviceGroup *lm.RestDeviceGroup, err error) {
	filter := fmt.Sprintf("name:%s", url.QueryEscape(name))
	restResponse, apiResponse, err := d.LMClient.GetDeviceGroupList("name,id,parentId", -1, 0, filter)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to find device group %q: %v", name, _err)
	}

	log.Debugf("%#v", restResponse)

	for _, d := range restResponse.Data.Items {
		if d.ParentId == parentID {
			log.Infof("Found device group %q with id %d", name, parentID)
			deviceGroup = &d

			return
		}
	}

	return
}

func (d *DeviceTree) createDeviceGroup(name, appliesTo string, disableAlerting bool, parentID int32) (*lm.RestDeviceGroup, error) {
	restResponse, apiResponse, err := d.LMClient.AddDeviceGroup(lm.RestDeviceGroup{
		Name:            name,
		Description:     "A dynamic device group for Kubernetes.",
		ParentId:        parentID,
		AppliesTo:       appliesTo,
		DisableAlerting: disableAlerting,
	})
	if e := utilities.CheckAllErrors(restResponse, apiResponse, err); e != nil {
		return nil, fmt.Errorf("Failed to add device group: %v", e)
	}

	deviceGroup := &restResponse.Data
	log.Infof("Created device group with id %d", deviceGroup.Id)

	return deviceGroup, nil
}

func (d *DeviceTree) createClusterDeviceGroup() (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := "Kubernetes Cluster: " + d.Config.ClusterName
	appliesTo := "hasCategory(\"" + constants.ClusterCategory + "\") && auto.clustername ==\"" + d.Config.ClusterName + "\""

	clusterDeviceGroup, err = d.findDeviceGroup(1, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, d.Config.DisableAlerting, 1)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createServiceDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := "Services"
	appliesTo := ""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, d.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createServiceDeletedDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := constants.DeletedDeviceGroup
	appliesTo := "hasCategory(\"" + constants.NodeDeletedCategory + "\") && auto.clustername ==\"" + d.Config.ClusterName + "\""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, d.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createEtcdDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := "Etcd"
	appliesTo := "hasCategory(\"" + constants.EtcdCategory + "\") && auto.clustername ==\"" + d.Config.ClusterName + "\""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, d.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createEtcdDeletedDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := constants.DeletedDeviceGroup
	appliesTo := "hasCategory(\"" + constants.EtcdDeletedCategory + "\") && auto.clustername ==\"" + d.Config.ClusterName + "\""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, true, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createNodeDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := "Nodes"
	appliesTo := "hasCategory(\"" + constants.NodeCategory + "\") && auto.clustername ==\"" + d.Config.ClusterName + "\""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, d.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createNodeDeletedDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := constants.DeletedDeviceGroup
	appliesTo := "hasCategory(\"" + constants.NodeDeletedCategory + "\") && auto.clustername ==\"" + d.Config.ClusterName + "\""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, true, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createPodDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := "Pods"
	appliesTo := ""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, d.Config.DisableAlerting, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}

func (d *DeviceTree) createPodDeletedDeviceGroup(parentDeviceGroup *lm.RestDeviceGroup) (clusterDeviceGroup *lm.RestDeviceGroup, err error) {
	name := constants.DeletedDeviceGroup
	appliesTo := "hasCategory(\"" + constants.PodDeletedCategory + "\") && auto.clustername ==\"" + d.Config.ClusterName + "\""

	clusterDeviceGroup, err = d.findDeviceGroup(parentDeviceGroup.Id, name)
	if err != nil {
		return
	}

	if clusterDeviceGroup == nil {
		clusterDeviceGroup, err = d.createDeviceGroup(name, appliesTo, true, parentDeviceGroup.Id)
		if err != nil {
			return
		}
	}

	return
}
