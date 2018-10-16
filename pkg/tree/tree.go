package tree

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	log "github.com/sirupsen/logrus"
)

// DeviceTree manages the device tree representation of a Kubernetes cluster in LogicMonitor.
type DeviceTree struct {
	*types.Base
}

// nolint: dupl
func (d *DeviceTree) buildOptsSlice() []*devicegroup.Options {
	// The device group at index 0 will be the root device group for all subsequent device groups.
	return []*devicegroup.Options{
		{
			Name:            constants.ClusterDeviceGroupPrefix + d.Config.ClusterName,
			ParentID:        d.Config.ClusterGroupID,
			DisableAlerting: d.Config.DisableAlerting,
			AppliesTo:       devicegroup.NewAppliesToBuilder().HasCategory(constants.ClusterCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:          d.LMClient,
			DeleteDevices:   d.Config.DeleteDevices,
		},
		{
			Name:                  constants.EtcdDeviceGroupName,
			DisableAlerting:       d.Config.DisableAlerting,
			AppliesTo:             devicegroup.NewAppliesToBuilder().HasCategory(constants.EtcdCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:                d.LMClient,
			DeleteDevices:         d.Config.DeleteDevices,
			AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(constants.EtcdDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
		},

		{
			Name:            constants.NodeDeviceGroupName,
			DisableAlerting: d.Config.DisableAlerting,
			AppliesTo:       devicegroup.NewAppliesToBuilder(),
			Client:          d.LMClient,
			DeleteDevices:   d.Config.DeleteDevices,
		},
		{
			Name:                  constants.AllNodeDeviceGroupName,
			DisableAlerting:       d.Config.DisableAlerting,
			AppliesTo:             devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:                d.LMClient,
			DeleteDevices:         d.Config.DeleteDevices,
			AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
		},

		{
			Name: constants.ServiceDeviceGroupName,
			// Services are a WIP in the product, disable alerting for now,
			DisableAlerting:       true,
			AppliesTo:             devicegroup.NewAppliesToBuilder(),
			Client:                d.LMClient,
			DeleteDevices:         d.Config.DeleteDevices,
			AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(constants.ServiceDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
		},
		{
			Name:                  constants.PodDeviceGroupName,
			DisableAlerting:       d.Config.DisableAlerting,
			AppliesTo:             devicegroup.NewAppliesToBuilder(),
			Client:                d.LMClient,
			DeleteDevices:         d.Config.DeleteDevices,
			AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(constants.PodDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
		},
	}
}

// CreateDeviceTree creates the Device tree that will represent the cluster in LogicMonitor.
func (d *DeviceTree) CreateDeviceTree() (map[string]int32, error) {
	// create the parent cluster group first
	d.checkAndUpdateClusterGroup()

	deviceGroups := make(map[string]int32)
	for _, opts := range d.buildOptsSlice() {
		switch opts.Name {
		case constants.AllNodeDeviceGroupName:
			// the all nodes group should be nested in 'Nodes'
			opts.ParentID = deviceGroups[constants.NodeDeviceGroupName]
		case constants.ClusterDeviceGroupPrefix + d.Config.ClusterName:
			// don't do anything for the root cluster group
		default:
			opts.ParentID = deviceGroups[constants.ClusterDeviceGroupPrefix+d.Config.ClusterName]
		}

		id, err := devicegroup.Create(opts)
		if err != nil {
			return nil, err
		}
		deviceGroups[opts.Name] = id
	}

	return deviceGroups, nil
}

func (d *DeviceTree) checkAndUpdateClusterGroup() {
	// do not need to check the root group
	if d.Config.ClusterGroupID == constants.RootDeviceGroupID {
		return
	}

	// if the group does not exist anymore, we will add the cluster to the root group
	if !devicegroup.ExistsByID(d.Config.ClusterGroupID, d.LMClient) {
		log.Warnf("The device group (id=%v) does not exist, the cluster will be added to the root group", d.Config.ClusterGroupID)
		d.Config.ClusterGroupID = constants.RootDeviceGroupID
	}
}
