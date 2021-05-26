package tree

import (
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
)

// DeviceTree manages the device tree representation of a Kubernetes cluster in LogicMonitor.
type DeviceTree struct {
	*types.Base
	ResourceCache *devicecache.ResourceCache
}

// nolint: dupl
func (d *DeviceTree) buildOptsSlice() []*devicegroup.Options {
	// The device group at index 0 will be the root device group for all subsequent device groups.

	nodes := enums.Nodes
	pods := enums.Pods
	etcd := enums.ETCD
	deployments := enums.Deployments
	services := enums.Services
	hpas := enums.Hpas
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
			AppliesTo:             devicegroup.NewAppliesToBuilder().HasCategory(etcd.GetCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:                d.LMClient,
			DeleteDevices:         d.Config.DeleteDevices,
			AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(etcd.GetDeletedCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			CustomProperties:      devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.ETCD),
		},

		{
			Name:                              constants.NodeDeviceGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder(),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.Nodes),
		},
		{
			Name:                              constants.AllNodeDeviceGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder().HasCategory(nodes.GetCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToDeletedGroup:             devicegroup.NewAppliesToBuilder().HasCategory(nodes.GetDeletedCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder(),
		},

		{
			Name:                              constants.ServiceDeviceGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder(),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToDeletedGroup:             devicegroup.NewAppliesToBuilder().HasCategory(services.GetDeletedCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.Services),
		},
		{
			Name:                              constants.PodDeviceGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder(),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToDeletedGroup:             devicegroup.NewAppliesToBuilder().HasCategory(pods.GetDeletedCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.Pods),
		},
		{
			Name:                              constants.DeploymentDeviceGroupName,
			DisableAlerting:                   true,
			AppliesTo:                         devicegroup.NewAppliesToBuilder(),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToDeletedGroup:             devicegroup.NewAppliesToBuilder().HasCategory(deployments.GetDeletedCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.Deployments),
		},
		{
			Name:                              constants.HorizontalPodAutoscalerDeviceGroupName,
			DisableAlerting:                   d.Config.DisableAlerting,
			AppliesTo:                         devicegroup.NewAppliesToBuilder(),
			Client:                            d.LMClient,
			DeleteDevices:                     d.Config.DeleteDevices,
			AppliesToDeletedGroup:             devicegroup.NewAppliesToBuilder().HasCategory(hpas.GetDeletedCategory()).And().Auto("clustername").Equals(d.Config.ClusterName),
			FullDisplayNameIncludeClusterName: d.Config.FullDisplayNameIncludeClusterName,
			CustomProperties:                  devicegroup.NewPropertyBuilder().AddProperties(d.Config.DeviceGroupProperties.HPA),
		},
	}
}

// CreateDeviceTree creates the Device tree that will represent the cluster in LogicMonitor.
func (d *DeviceTree) CreateDeviceTree(lctx *lmctx.LMContext) (map[string]int32, error) {
	deviceGroups := make(map[string]int32)
	for _, opts := range d.buildOptsSlice() {
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
	}

	return deviceGroups, nil
}
