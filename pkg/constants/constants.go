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
	// RootDeviceGroupID is the root ID in the device tree.
	RootDeviceGroupID = int32(1)
	// MaxResourceLength is the max length of the resource name
	MaxResourceLength = 244
)

const (
	// LabelNodeRole is the label name used to specify a node's role in the cluster
	LabelNodeRole = "node-role.kubernetes.io/"
	// LabelCustomPropertyPrefix is the prefix to use for custom properties based of labels
	LabelCustomPropertyPrefix = "kubernetes.label."
	// LabelNullPlaceholder is the string used to represent null values in custom properties
	LabelNullPlaceholder = "null"
	// LabelFargateProfile is the label name used for fargate profile to distinguish between fargate & other Pods
	LabelFargateProfile = "eks.amazonaws.com/fargate-profile"
)

const (
	// AllNodeDeviceGroupName is the service device group name in the cluster device group.
	AllNodeDeviceGroupName = "All"
	// EtcdDeviceGroupName is the etcd device group name in the cluster device group.
	EtcdDeviceGroupName = "Etcd"
	// NodeDeviceGroupName is the top-level device group name in the cluster device group.
	NodeDeviceGroupName = "Nodes"
	// NamespacesGroupName is the namespaces device group name in the cluster device group.
	NamespacesGroupName = "Namespaces"
	// PodDeviceGroupName is the pod device group name in the cluster device group.
	PodDeviceGroupName = "Pods"
	// ServiceDeviceGroupName is the service device group name in the cluster device group.
	ServiceDeviceGroupName = "Services"
	// DeploymentDeviceGroupName is the deployment device group name in the cluster device group.
	DeploymentDeviceGroupName = "Deployments"
	// HorizontalPodAutoscalerDeviceGroupName is the deployment device group name in the cluster device group.
	HorizontalPodAutoscalerDeviceGroupName = "HorizontalPodAutoscalers"
)

const (
	// ClusterCategory is the system.category used to identity the Kubernetes cluster in LogicMonitor.
	ClusterCategory = "KubernetesCluster"
	// NamespaceCategory is the system.category used to identity the Kubernetes Namespace resource type in LogicMonitor.
	NamespaceCategory = "KubernetesNamespace"
	// EtcdCategory is the system.category used to identity the Kubernetes Pod resource type in LogicMonitor.
	EtcdCategory = "KubernetesEtcd"
	// EtcdDeletedCategory is the system.category used to identity a deleted Kubernetes Etcd node in LogicMonitor.
	EtcdDeletedCategory = "KubernetesEtcdDeleted"
	// NodeCategory is the system.category used to identity the Kubernetes Node resource type in LogicMonitor.
	NodeCategory = "KubernetesNode"
	// NodeDeletedCategory is the system.category used to identity a deleted Kubernetes Node resource type in LogicMonitor.
	NodeDeletedCategory = "KubernetesNodeDeleted"
	// NodeConflictCategory is the system.category used to identity a conflicting Kubernetes Node resource type in LogicMonitor.
	NodeConflictCategory = "KubernetesNodeConflict"
	// ServiceCategory is the system.category used to identity a Kubernetes Service resource type in LogicMonitor.
	ServiceCategory = "KubernetesService"
	// DeploymentCategory is the system.category used to identity a Kubernetes Service resource type in LogicMonitor.
	DeploymentCategory = "KubernetesDeployment"
	// PodCategory is the system.category used to identity the Kubernetes Pod resource type in LogicMonitor.
	PodCategory = "KubernetesPod"
	// HorizontalPodAutoscalerCategory is the system.category used to identity the Kubernetes HorizontalPodAutoscaler resource type in LogicMonitor.
	HorizontalPodAutoscalerCategory = "KubernetesHorizontalPodAutoscaler"
	// DeletedDeviceGroup is the name of the device group where deleted devices are optionally moved to.
	DeletedDeviceGroup = "_deleted"
	// ClusterDeviceGroupPrefix is the prefix for the top level cluster device group
	ClusterDeviceGroupPrefix = "Kubernetes Cluster: "
	// ConflictDeviceGroup is the name of the device group where conflicting devices are optionally moved to.
	ConflictDeviceGroup = "_conflict"
)

const (
	// K8sClusterNamePropertyKey is the key of the unique auto property kubernetes cluster name
	K8sClusterNamePropertyKey = "auto.clustername"
	// K8sResourceCreatedOnPropertyKey is the key of the custom property used to record resource create timestamp
	K8sResourceCreatedOnPropertyKey = "kubernetes.resourceCreatedOn"
	// K8sResourceDeletedOnPropertyKey is the key of the custom property used to record resource deleted timestamp
	K8sResourceDeletedOnPropertyKey = "kubernetes.resourceDeletedOn"
	// K8sDeviceType is the type value of the k8s device
	K8sDeviceType = 8
	// K8sSystemCategoriesPropertyKey is the key of the unique custom property kubernetes system categories
	K8sSystemCategoriesPropertyKey = "system.categories"
	// K8sSystemIPsPropertyKey is the key of the system ips property
	K8sSystemIPsPropertyKey = "system.ips"

	// K8sDeviceNamePropertyKey is the key of the unique auto property kubernetes device name.
	K8sDeviceNamePropertyKey = "auto.name"
	// K8sDeviceNamespacePropertyKey is the key of the unique auto property kubernetes device namespace.
	K8sDeviceNamespacePropertyKey = "auto.namespace"
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
	// DeviceGroupCustomType is the device group of custom type
	DeviceGroupCustomType = "custom"
	// HistorySuffix is the key suffix used for maintaining history
	HistorySuffix = ".history"
	// PropertySeparator is the property separator
	PropertySeparator = ", "
)

const (
	// DefaultPeriodicSyncInterval Default interval for Periodic Discovery.
	DefaultPeriodicSyncInterval = time.Minute * 1

	// DefaultPeriodicDeleteInterval Default interval for Periodic delete.
	DefaultPeriodicDeleteInterval = time.Minute * 30

	// DefaultCacheResyncInterval Default interval for cache resync.
	DefaultCacheResyncInterval = time.Minute * 5
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
	// EnvVarArgusConfigPrefix prefix to parse environment variables into config struct
	EnvVarArgusConfigPrefix = "argus"
)
