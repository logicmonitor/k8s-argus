package constants

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// Version is the Argus version and is set at build time.
	Version string

	// DefaultListOptions default list all
	DefaultListOptions = metav1.ListOptions{} // nolint: exhaustivestruct
)

const (
	// UserAgentBase is the base string for the User-Agent HTTP header.
	UserAgentBase = "LogicMonitor Argus/"
)

const (
	// RootResourceGroupID is the root ID in the resource tree.
	RootResourceGroupID = int32(1)
	// MaxResourceLength is the max length of the resource name
	MaxResourceLength = 244
)

const (
	// LabelNodeRole is the label name used to specify a node's role in the cluster
	LabelNodeRole = "node-role.kubernetes.io/"
	// LabelCustomPropertyPrefix is the prefix to use for custom properties based of labels
	LabelCustomPropertyPrefix = "kubernetes.label."
	// AnnotationCustomPropertyPrefix is the prefix to use for custom properties based of labels
	AnnotationCustomPropertyPrefix = "kubernetes.annotation."
	// SelectorCustomPropertyPrefix is the prefix to use for custom properties based of labels
	SelectorCustomPropertyPrefix = "kubernetes.selector."
	// MatchLabelsKey is prefix for MatchLabels
	MatchLabelsKey = "matchLabels"
	// NodeSelectorKey is the prefix for NodeSelector
	NodeSelectorKey = "nodeSelector"
	// PodSelectorKey is the prefix for NetworkPolicies
	PodSelectorKey = "podSelector"
	// LogicalEQUALS used for LogicalEquals
	LogicalEQUALS = " == "
	// LogicalAND used for LogicalAND
	LogicalAND = " && "
	// LabelNullPlaceholder is the string used to represent null values in custom properties
	LabelNullPlaceholder = "null"
	// LabelFargateProfile is the label name used for fargate profile to distinguish between fargate & other Pods
	LabelFargateProfile = "eks.amazonaws.com/fargate-profile"
)

const (
	// AllNodeResourceGroupName is the service resource group name in the cluster resource group.
	AllNodeResourceGroupName = "All"
	// EtcdResourceGroupName is the etcd resource group name in the cluster resource group.
	EtcdResourceGroupName = "Etcd"
	// NodeResourceGroupName is the top-level resource group name in the cluster resource group.
	NodeResourceGroupName = "Nodes"
	// NamespacesGroupName is the namespaces resource group name in the cluster resource group.
	NamespacesGroupName = "Namespaces"
	// ClusterScopedGroupName is the resource group name for Cluster Scoped Resources in the cluster resource group.
	ClusterScopedGroupName = "ClusterScoped"
)

const (
	// ClusterCategory is the system.category used to identity the Kubernetes cluster in LogicMonitor.
	ClusterCategory = "KubernetesCluster"
	// DeletedResourceGroup is the name of the resource group where deleted resources are optionally moved to.
	DeletedResourceGroup = "_deleted"
	// ClusterResourceGroupPrefix is the prefix for the top level cluster resource group
	ClusterResourceGroupPrefix = "Kubernetes Cluster: "
)

const (
	// K8sClusterNamePropertyKey is the key of the unique auto property kubernetes cluster name
	K8sClusterNamePropertyKey = "auto.clustername"
	// K8sResourceCreatedOnPropertyKey is the key of the custom property used to record resource create timestamp
	K8sResourceCreatedOnPropertyKey = "kubernetes.resourceCreatedOn"
	// K8sResourceDeletedOnPropertyKey is the key of the custom property used to record resource deleted timestamp
	K8sResourceDeletedOnPropertyKey = "kubernetes.resourceDeletedOn"
	// K8sResourceType is the type value of the k8s resource
	K8sResourceType = 8
	// K8sSystemCategoriesPropertyKey is the key of the unique custom property kubernetes system categories
	K8sSystemCategoriesPropertyKey = "system.categories"
	// K8sSystemIPsPropertyKey is the key of the system ips property
	K8sSystemIPsPropertyKey = "system.ips"

	// K8sResourceNamePropertyKey is the key of the unique auto property kubernetes resource name.
	K8sResourceNamePropertyKey = "auto.name"
	// K8sResourceNamespacePropertyKey is the key of the unique auto property kubernetes resource namespace.
	K8sResourceNamespacePropertyKey = "auto.namespace"
	// K8sResourceUIDPropertyKey is the key of the unique auto property kubernetes resource uid.
	K8sResourceUIDPropertyKey = "auto.uid"
)

const (
	// K8sAPIVersionV1 is the version 'v1' of k8s api
	K8sAPIVersionV1 = "v1"
	// K8sAPIVersionAppsV1beta2 is the version 'apps/v1beta2' of k8s api
	K8sAPIVersionAppsV1beta2 = "apps/v1beta2"
	// K8sAPIVersionAppsV1 is the version 'apps/v1' of k8s api
	K8sAPIVersionAppsV1 = "apps/v1"
	// K8sAutoscalingV1 is the version 'autoscaling/v1' of k8s api
	K8sAutoscalingV1 = "autoscaling/v1"
	// K8sAPIVersionBatchV1 is the version 'batch/v1' of k8s api
	K8sAPIVersionBatchV1 = "batch/v1"
	// K8sAPIVersionBatchV1Beta1 is the version 'batch/v1beta1' of k8s api
	K8sAPIVersionBatchV1Beta1 = "batch/v1beta1"
	// K8sAPIVersionExtensionsV1Beta1 is the version 'extensions/v1beta1' of k8s api
	K8sAPIVersionExtensionsV1Beta1 = "extensions/v1beta1"
	// K8sAPIVersionNetworkingV1 is the version 'networking/v1' of k8s api
	K8sAPIVersionNetworkingV1 = "networking/v1"
	// K8sAPIVersionNetworkingV1Beta1 is the version 'networking/v1beta1' of k8s api
	K8sAPIVersionNetworkingV1Beta1 = "networking/v1beta1"
)

const (
	// ArgusAppVersion is the key for Argus app version
	ArgusAppVersion = "argus.app-version"
	// HelmChart is the key for Argus & Collectoeset-controller label
	HelmChart = "helm-chart"
	// HelmRevision is the key for Argus & Collectoeset-controller label
	HelmRevision = "helm-revision"
	// Chart is the label key in Argus & Collectoeset-controller Deployment
	Chart = "chart"
	// Argus is the Argus Deployment label
	Argus = "argus"
	// CollectorsetController is the Collectorset-controller Deployment label
	CollectorsetController = "collectorset-controller"
	// KubernetesVersionKey is the key for customProperties
	KubernetesVersionKey = "kubernetes.version"
	// ResourceGroupCustomType is the resource group of custom type
	ResourceGroupCustomType = "custom"
	// HistorySuffix is the key suffix used for maintaining history
	HistorySuffix = ".history"
	// PropertySeparator is the property separator
	PropertySeparator = ", "
	// ArgusHelmChartAuditKey audit entry key
	ArgusHelmChartAuditKey = Argus + "." + HelmChart
	// CSCHelmChartAuditKey audit entry key
	CSCHelmChartAuditKey = CollectorsetController + "." + HelmChart
	// ArgusHelmRevisionAuditKey audit entry key
	ArgusHelmRevisionAuditKey = Argus + "." + HelmRevision
	// CSCHelmRevisionAuditKey audit entry key
	CSCHelmRevisionAuditKey = CollectorsetController + "." + HelmRevision
)

const (
	// DefaultPeriodicSyncInterval Default interval for Periodic Discovery.
	DefaultPeriodicSyncInterval = time.Minute * 30

	// DefaultPeriodicDeleteInterval Default interval for Periodic delete.
	DefaultPeriodicDeleteInterval = time.Minute * 5

	// DefaultCacheResyncInterval Default interval for cache resync.
	DefaultCacheResyncInterval = time.Minute * 10
)

const (
	// IsLocal flag used to set when application is loaded using external kubeconfig file to indicate it is running outside cluster
	IsLocal = "IS_LOCAL"
)

const (
	// ConfigInitK8sClientExitCode exit 1
	ConfigInitK8sClientExitCode = 1
	// ConfigInitExitCode exit 2
	ConfigInitExitCode = 2
	// GetConfigExitCode exit 3
	GetConfigExitCode = 3
)

const (
	// ResyncConflictingResourcesProp graceful resync flag on cluster group
	ResyncConflictingResourcesProp = "resync.conflicting_resources"
	// ResyncCacheProp graceful resync flag on cluster group
	ResyncCacheProp = "resync.cache"
)

const (
	// ConfigFileName config file name to read from configmap
	ConfigFileName = "config.yaml"
	// FiltersConfigFileName config file name to read from configmap
	FiltersConfigFileName = "filters-config.yaml"
	// EnvVarArgusConfigPrefix prefix to parse environment variables into config struct
	EnvVarArgusConfigPrefix = "argus"
)

const (
	// PartitionKey partition key used to send lm requests to a single worker all the time.
	PartitionKey = "partition_key"
)

const (
	// AutoPropCreatedBy auto property to identify who created property
	AutoPropCreatedBy = "auto.createdBy"

	// DGCustomPropCreatedBy resource group does not have auto prop, so adding it as custom prop
	DGCustomPropCreatedBy = "createdBy"

	// CreatedByPrefix created by prefix
	CreatedByPrefix = "LogicMonitor/Argus: "
)

const (
	// NameFieldName is the field name for a device's name.
	NameFieldName = "name"
)

const HeadlessServiceIPNone = "None"
