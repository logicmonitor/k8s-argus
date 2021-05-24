package tree

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
)

// DeviceTree2 manages the device tree representation of a Kubernetes cluster in LogicMonitor.
type DeviceTree2 struct {
	*types.Base
	ResourceCache *devicecache.ResourceCache
}

// nolint: dupl, unused
func (d *DeviceTree2) buildOptsSlice() []*devicegroup.Options {
	// The device group at index 0 will be the root device group for all subsequent device groups.

	return []*devicegroup.Options{
		{
			Name:             util.ClusterGroupName(d.Config.ClusterName),
			ParentID:         d.Config.ClusterGroupID,
			DisableAlerting:  d.Config.DisableAlerting,
			AppliesTo:        devicegroup.NewAppliesToBuilder().HasCategory(constants.ClusterCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:           d.LMClient,
			DeleteDevices:    d.Config.DeleteDevices,
			CustomProperties: devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.Cluster),
		},
		{
			Name:                  constants.EtcdDeviceGroupName,
			DisableAlerting:       d.Config.DisableAlerting,
			AppliesTo:             devicegroup.NewAppliesToBuilder().HasCategory(constants.EtcdCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:                d.LMClient,
			DeleteDevices:         d.Config.DeleteDevices,
			AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(constants.EtcdDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			CustomProperties:      devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.ETCD),
		},

		{
			Name:                              constants.NodeDeviceGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder(),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToConflict:                 devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeConflictCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.Nodes),
		},
		{
			Name:                              constants.AllNodeDeviceGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToDeletedGroup:             devicegroup.NewAppliesToBuilder().HasCategory(constants.NodeDeletedCategory).And().Auto("clustername").Equals(d.Config.ClusterName),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder(),
		},
		{
			Name:                              constants.NamespacesGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder(),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToDeletedGroup:             devicegroup.NewAppliesToBuilder(),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder(),
		},
	}
}

// CreateDeviceTree creates the Device tree that will represent the cluster in LogicMonitor.
// nolint: dupl
func (d *DeviceTree2) CreateDeviceTree(lctx *lmctx.LMContext) (map[string]int32, error) {
	deviceGroups := make(map[string]int32)
	/*for _, opts := range d.buildOptsSlice() {
			switch opts.Name {
			case constants.AllNodeDeviceGroupName:
				// the all nodes group should be nested in 'Nodes'
				opts.ParentID = deviceGroups[constants.NodeDeviceGroupName]
			case util.ClusterGroupName(d.Config.ClusterName):
				// don't do anything for the root cluster group
			default:
				opts.ParentID = deviceGroups[util.ClusterGroupName(d.Config.ClusterName)]
			}

			id, err := devicegroup.Create(lctx, opts, d.ResourceCache)
			if err != nil {

	return nil, err
			}
			deviceGroups[opts.Name] = id
		}*/

	return deviceGroups, nil
}
