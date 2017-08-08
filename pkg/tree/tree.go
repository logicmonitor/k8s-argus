package tree

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/tree/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	log "github.com/sirupsen/logrus"
)

// DeviceTree manages the device tree representation of a Kubernetes cluster in LogicMonitor.
type DeviceTree struct {
	*types.Base
}

// CreateDeviceTree creates the Device tree that will represent the cluster in the LogicMonitor portal.
func (d *DeviceTree) CreateDeviceTree() (map[string]int32, error) {
	deviceGroups := make(map[string]int32)
	clusterDeviceGroupID, err := d.createClusterDeviceGroup()
	if err != nil {
		return nil, err
	}
	log.Infof("Using cluster device group with id %d", clusterDeviceGroupID)

	serviceDeviceGroupID, err := d.createServiceDeviceGroup(clusterDeviceGroupID)
	if err != nil {
		return nil, err
	}
	deviceGroups["services"] = serviceDeviceGroupID
	_, err = d.createServiceDeletedDeviceGroup(serviceDeviceGroupID)
	if err != nil {
		return nil, err
	}
	log.Infof("Using service device group with id %d", serviceDeviceGroupID)

	etcdDeviceGroupID, err := d.createEtcdDeviceGroup(clusterDeviceGroupID)
	if err != nil {
		return nil, err
	}
	_, err = d.createEtcdDeletedDeviceGroup(etcdDeviceGroupID)
	if err != nil {
		return nil, err
	}
	log.Infof("Using etcd device group with id %d", etcdDeviceGroupID)

	nodeDeviceGroupID, err := d.createNodeDeviceGroup(clusterDeviceGroupID)
	if err != nil {
		return nil, err
	}
	_, err = d.createNodeDeletedDeviceGroup(nodeDeviceGroupID)
	if err != nil {
		return nil, err
	}
	log.Infof("Using node device group with id %d", nodeDeviceGroupID)

	podDeviceGroupID, err := d.createPodDeviceGroup(clusterDeviceGroupID)
	if err != nil {
		return nil, err
	}
	deviceGroups["pods"] = podDeviceGroupID
	_, err = d.createPodDeletedDeviceGroup(podDeviceGroupID)
	if err != nil {
		return nil, err
	}
	log.Infof("Using pod device group with id %d", podDeviceGroupID)

	return deviceGroups, nil
}

func (d *DeviceTree) createClusterDeviceGroup() (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(constants.ClusterCategory).And().Auto("clustername").Equals(d.Config.ClusterName)
	opts := &devicegroup.Options{
		Name:            "Kubernetes Cluster: " + d.Config.ClusterName,
		ParentID:        constants.RootDeviceGroupID,
		DisableAlerting: d.Config.DisableAlerting,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createServiceDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder()
	opts := &devicegroup.Options{
		Name:            constants.ServiceDeviceGroupName,
		ParentID:        parentID,
		DisableAlerting: d.Config.DisableAlerting,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createServiceDeletedDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName)
	opts := &devicegroup.Options{
		Name:            constants.DeletedDeviceGroup,
		ParentID:        parentID,
		DisableAlerting: true,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createEtcdDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(constants.EtcdCategory).And().Auto("clustername").Equals(d.Config.ClusterName)
	opts := &devicegroup.Options{
		Name:            constants.EtcdDeviceGroupName,
		ParentID:        parentID,
		DisableAlerting: d.Config.DisableAlerting,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createEtcdDeletedDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(constants.EtcdDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName)
	opts := &devicegroup.Options{
		Name:            constants.DeletedDeviceGroup,
		ParentID:        parentID,
		DisableAlerting: true,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createNodeDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeCategory).And().Auto("clustername").Equals(d.Config.ClusterName)
	opts := &devicegroup.Options{
		Name:            constants.NodeDeviceGroupName,
		ParentID:        parentID,
		DisableAlerting: d.Config.DisableAlerting,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createNodeDeletedDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName)
	opts := &devicegroup.Options{
		Name:            constants.DeletedDeviceGroup,
		ParentID:        parentID,
		DisableAlerting: true,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createPodDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder()
	opts := &devicegroup.Options{
		Name:            constants.PodDeviceGroupName,
		ParentID:        parentID,
		DisableAlerting: d.Config.DisableAlerting,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DeviceTree) createPodDeletedDeviceGroup(parentID int32) (int32, error) {
	appliesTo := devicegroup.NewAppliesToBuilder().HasCategory(constants.PodDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName)
	opts := &devicegroup.Options{
		Name:            constants.DeletedDeviceGroup,
		ParentID:        parentID,
		DisableAlerting: true,
		AppliesTo:       appliesTo,
		Client:          d.LMClient,
	}
	id, err := devicegroup.Create(opts)
	if err != nil {
		return 0, err
	}

	return id, nil
}
