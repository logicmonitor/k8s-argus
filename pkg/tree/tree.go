package tree

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/types"
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
			Name:            "Kubernetes Cluster: " + d.Config.ClusterName,
			ParentID:        constants.RootDeviceGroupID,
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
			Name:                  constants.NodeDeviceGroupName,
			DisableAlerting:       d.Config.DisableAlerting,
			AppliesTo:             devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:                d.LMClient,
			DeleteDevices:         d.Config.DeleteDevices,
			AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
		},
		{
			Name:                  constants.ServiceDeviceGroupName,
			DisableAlerting:       d.Config.DisableAlerting,
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
	deviceGroups := make(map[string]int32)
	for _, opts := range d.buildOptsSlice() {
		if opts.Name != "Kubernetes Cluster: "+d.Config.ClusterName {
			opts.ParentID = deviceGroups["Kubernetes Cluster: "+d.Config.ClusterName]
		}
		id, err := devicegroup.Create(opts)
		if err != nil {
			return nil, err
		}
		deviceGroups[opts.Name] = id
	}

	return deviceGroups, nil
}
