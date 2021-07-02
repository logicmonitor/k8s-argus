package tree

import (
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup/dgbuilder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
)

// GetResourceGroupTree creates the ResourceGroup tree that will represent the cluster in LogicMonitor.
func GetResourceGroupTree(lctx *lmctx.LMContext, dgBuilder types.ResourceManager, requester *types.LMRequester) (*types.ResourceGroupTree, error) {
	conf, err2 := getConf(lctx, requester)
	if err2 != nil {
		return nil, err2
	}
	nodes := enums.Nodes
	pods := enums.Pods
	etcd := enums.ETCD
	deployments := enums.Deployments
	services := enums.Services
	hpas := enums.Hpas
	doNotCreateDeletedGroup := conf.DeleteResources
	return &types.ResourceGroupTree{
		Options: []types.ResourceGroupOption{
			dgBuilder.GroupName(util.ClusterGroupName(conf.ClusterName)),
			dgBuilder.ParentID(conf.ClusterGroupID),
			dgBuilder.DisableAlerting(conf.DisableAlerting),
			dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(constants.ClusterCategory).And().Auto("clustername").Equals(conf.ClusterName)),
			dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.Cluster)),
		},
		ChildGroups: []*types.ResourceGroupTree{
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(deployments.TitlePlural()),
					dgBuilder.DisableAlerting(true),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.Deployments)),
				},
				ChildGroups: []*types.ResourceGroupTree{
					{
						DontCreate: doNotCreateDeletedGroup,
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.DeletedResourceGroup),
							dgBuilder.DisableAlerting(true),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(deployments.GetDeletedCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
				},
			},
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(hpas.TitlePlural()),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.HPA)),
				},
				ChildGroups: []*types.ResourceGroupTree{
					{
						DontCreate: doNotCreateDeletedGroup,
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.DeletedResourceGroup),
							dgBuilder.DisableAlerting(true),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(hpas.GetDeletedCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
				},
			},
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(nodes.TitlePlural()),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.Nodes)),
				},
				ChildGroups: []*types.ResourceGroupTree{
					{
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.AllNodeResourceGroupName),
							dgBuilder.DisableAlerting(conf.DisableAlerting),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(nodes.GetCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
					{
						DontCreate: doNotCreateDeletedGroup,
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.DeletedResourceGroup),
							dgBuilder.DisableAlerting(true),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(nodes.GetDeletedCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
				},
			},
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(pods.TitlePlural()),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.Pods)),
				},
				ChildGroups: []*types.ResourceGroupTree{
					{
						DontCreate: doNotCreateDeletedGroup,
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.DeletedResourceGroup),
							dgBuilder.DisableAlerting(true),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(pods.GetDeletedCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
				},
			},
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(services.TitlePlural()),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.Services)),
				},
				ChildGroups: []*types.ResourceGroupTree{
					{
						DontCreate: doNotCreateDeletedGroup,
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.DeletedResourceGroup),
							dgBuilder.DisableAlerting(true),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(services.GetDeletedCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
				},
			},
			{
				ChildGroups: nil,
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(constants.EtcdResourceGroupName),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
					dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(etcd.GetCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.ETCD)),
				},
			},
		},
	}, nil
}

// GetResourceGroupTree2 creates the Resource tree that will represent the cluster in LogicMonitor.
// nolint: dupl
func GetResourceGroupTree2(lctx *lmctx.LMContext, dgBuilder types.ResourceManager, requester *types.LMRequester) (*types.ResourceGroupTree, error) {
	conf, err2 := getConf(lctx, requester)
	if err2 != nil {
		return nil, err2
	}
	etcd := enums.ETCD
	nodes := enums.Nodes
	doNotCreateDeletedGroup := conf.DeleteResources
	deletedBuilder := dgbuilder.NewAppliesToBuilder().
		Auto("clustername").Equals(conf.ClusterName).And().
		OpenBracket()
	for _, e := range enums.ALLResourceTypes {
		if e == enums.Namespaces || !e.IsNamespaceScopedResource() {
			continue
		}
		deletedBuilder = deletedBuilder.HasCategory(e.GetDeletedCategory()).Or()
	}
	deletedBuilder.TrimOrCloseBracket()
	return &types.ResourceGroupTree{
		Options: []types.ResourceGroupOption{
			dgBuilder.GroupName(util.ClusterGroupName(conf.ClusterName)),
			dgBuilder.ParentID(conf.ClusterGroupID),
			dgBuilder.DisableAlerting(conf.DisableAlerting),
			dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(constants.ClusterCategory).And().Auto("clustername").Equals(conf.ClusterName)),
			dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.Cluster)),
		},
		ChildGroups: []*types.ResourceGroupTree{
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(constants.EtcdResourceGroupName),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
					dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(etcd.GetCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.ETCD)),
				},
			},
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(nodes.TitlePlural()),
					dgBuilder.CustomProperties(dgbuilder.NewPropertyBuilder().AddProperties(conf.ResourceGroupProperties.Nodes)),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
				},
				ChildGroups: []*types.ResourceGroupTree{
					{
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.AllNodeResourceGroupName),
							dgBuilder.DisableAlerting(conf.DisableAlerting),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(nodes.GetCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
					{
						DontCreate: doNotCreateDeletedGroup,
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.DeletedResourceGroup),
							dgBuilder.DisableAlerting(true),
							dgBuilder.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(nodes.GetDeletedCategory()).And().Auto("clustername").Equals(conf.ClusterName)),
						},
					},
				},
			},
			{
				Options: []types.ResourceGroupOption{
					dgBuilder.GroupName(constants.NamespacesGroupName),
					dgBuilder.DisableAlerting(conf.DisableAlerting),
				},
				ChildGroups: []*types.ResourceGroupTree{
					{
						DontCreate: doNotCreateDeletedGroup,
						Options: []types.ResourceGroupOption{
							dgBuilder.GroupName(constants.DeletedResourceGroup),
							dgBuilder.DisableAlerting(true),
							dgBuilder.AppliesTo(deletedBuilder),
						},
					},
				},
			},
		},
	}, nil
}

func getConf(lctx *lmctx.LMContext, requester *types.LMRequester) (*config.Config, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	// check and update the params
	if err := checkAndUpdateClusterGroup(lctx, conf, requester); err != nil {
		return nil, err
	}
	return conf, nil
}

// check the cluster group ID, if the group does not exist, just use the root group
func checkAndUpdateClusterGroup(lctx *lmctx.LMContext, config *config.Config, lmClient *types.LMRequester) error {
	log := lmlog.Logger(lctx)
	// do not need to check the root group
	if config.ClusterGroupID == constants.RootResourceGroupID {
		return nil
	}

	rg, err := resourcegroup.GetByID(lctx, config.ClusterGroupID, lmClient)
	if err != nil {
		log.Errorf("Failed to search cluster resource group [%d]: %s", config.ClusterGroupID, err)
		return err
	}
	// if the group does not exist anymore, we will add the cluster to the root group
	if rg == nil {
		log.Warnf("The resource group (id=%v) does not exist, the cluster will be added to the root group", config.ClusterGroupID)
		config.ClusterGroupID = constants.RootResourceGroupID
	}
	return nil
}
