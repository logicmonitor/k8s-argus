package constants

import "time"

var (
	// Version is the Argus version and is set at build time.
	Version string
)

const (
	// UserAgentBase is the base string for the User-Agent HTTP header.
	UserAgentBase = "LogicMonitor Argus/"
)

const (
	// RootDeviceGroupID is the root ID in the device tree.
	RootDeviceGroupID = 1
	// CustomPropertiesFieldName is the field name for a device's custom properties.
	CustomPropertiesFieldName = "customProperties"
	// NameFieldName is the field name for a device's name.
	NameFieldName = "name"
	// DisplayNameFieldName is the field name for a device's display name.
	DisplayNameFieldName = "displayName"
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
	// ServiceDeletedCategory is the system.category used to identity a deleted Kubernetes Service resource type in LogicMonitor.
	ServiceDeletedCategory = "KubernetesServiceDeleted"
	// ServiceConflictCategory is the system.category used to identity a conflicting Kubernetes Service resource type in LogicMonitor.
	ServiceConflictCategory = "KubernetesServiceConflict"
	// DeploymentCategory is the system.category used to identity a Kubernetes Service resource type in LogicMonitor.
	DeploymentCategory = "KubernetesDeployment"
	// DeploymentDeletedCategory is the system.category used to identity a deleted Kubernetes Service resource type in LogicMonitor.
	DeploymentDeletedCategory = "KubernetesDeploymentDeleted"
	// DeploymentConflictCategory is the system.category used to identity a conflicting Kubernetes Deployment resource type in LogicMonitor.
	DeploymentConflictCategory = "KubernetesDeploymentConflict"
	// PodCategory is the system.category used to identity the Kubernetes Pod resource type in LogicMonitor.
	PodCategory = "KubernetesPod"
	// PodDeletedCategory is the system.category used to identity a deleted Kubernetes Pod resource type in LogicMonitor.
	PodDeletedCategory = "KubernetesPodDeleted"
	// PodConflictCategory is the system.category used to identity a conflicting Kubernetes Pod resource type in LogicMonitor.
	PodConflictCategory = "KubernetesPodConflict"
	// HorizontalPodAutoscalerCategory is the system.category used to identity the Kubernetes HorizontalPodAutoscaler resource type in LogicMonitor.
	HorizontalPodAutoscalerCategory = "KubernetesHorizontalPodAutoscaler"
	// HorizontalPodAutoscalerDeletedCategory is the system.category used to identity a deleted Kubernetes HorizontalPodAutoscaler resource type in LogicMonitor.
	HorizontalPodAutoscalerDeletedCategory = "KubernetesHorizontalPodAutoscalerDeleted"
	// HorizontalPodAutoscalerConflictCategory is the system.category used to identity a conflicting Kubernetes HorizontalPodAutoscaler resource type in LogicMonitor.
	HorizontalPodAutoscalerConflictCategory = "KubernetesHorizontalPodAutoscalerConflict"
	// DeletedDeviceGroup is the name of the device group where deleted devices are optionally moved to.
	DeletedDeviceGroup = "_deleted"
	// ClusterDeviceGroupPrefix is the prefix for the top level cluster device group
	ClusterDeviceGroupPrefix = "Kubernetes Cluster: "
	//ConflictDeviceGroup is the name of the device group where conflicting devices are optionally moved to.
	ConflictDeviceGroup = "_conflict"
)

const (
	// ConfigPath is the path used to read the config.yaml file from.
	ConfigPath = "/etc/argus/config.yaml"
	// AccessID is the environment variable name to lookup for the LogicMonitor access ID.
	AccessID = "ARGUS_ACCESS_ID"
	// AccessKey is the environment variable name to lookup for the LogicMonitor access key.
	AccessKey = "ARGUS_ACCESS_KEY"
	// Account is the environment variable name to lookup for the LogicMonitor account.
	Account = "ARGUS_ACCOUNT"
)

const (
	// K8sClusterNamePropertyKey is the key of the unique auto property kubernetes cluster name
	K8sClusterNamePropertyKey = "auto.clustername"
	// K8sResourceNamePropertyKey is the key of the custom property used to record resource name
	K8sResourceNamePropertyKey = "auto.resourcename"
	// K8sResourceCreatedOnPropertyKey is the key of the custom property used to record resource create timestamp
	K8sResourceCreatedOnPropertyKey = "kubernetes.resourceCreatedOn"
	// K8sResourceDeletedOnPropertyKey is the key of the custom property used to record resource deleted timestamp
	K8sResourceDeletedOnPropertyKey = "kubernetes.resourceDeletedOn"
	// K8sResourceDeleteAfterDurationPropertyKey is the key of the custom property used to delete resources from the portal after specified time
	K8sResourceDeleteAfterDurationPropertyKey = "kubernetes.resourcedeleteafterduration"
	// K8sResourceDeleteAfterDurationPropertyValue is the default value of the custom property used to delete resources from the portal after specified time
	K8sResourceDeleteAfterDurationPropertyValue = "P1DT0H0M0S"
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
	// K8sAPIVersionAppsV1beta1 is the version 'apps/v1beta1' of k8s api
	K8sAPIVersionAppsV1beta1 = "apps/v1beta1"
	// K8sAPIVersionAppsV1beta2 is the version 'apps/v1beta2' of k8s api
	K8sAPIVersionAppsV1beta2 = "apps/v1beta2"
	// K8sAPIVersionAppsV1 is the version 'apps/v1' of k8s api
	K8sAPIVersionAppsV1 = "apps/v1"
	// K8sAutoscalingV1 is the version 'autoscaling/v1' of k8s api
	K8sAutoscalingV1 = "autoscaling/v1"
)

const (
	// Deployments deployments generic
	Deployments = "deployments"
	// Pods pods generic
	Pods = "pods"
	// Services Services generic
	Services = "services"
	// Nodes Nodes generic
	Nodes = "nodes"
	// HorizontalPodAutoScalers hpa generic
	HorizontalPodAutoScalers = "horizontalpodautoscalers"
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
	// IsPingDevice is the key used in watcher context to pass metadata
	IsPingDevice = "ispingdevice"

	// ResyncPodsClusterProperty is a cluster property for graceful pod updates
	ResyncPodsClusterProperty = "resync.pods"
)

const (
	// DefaultPeriodicSyncInterval Default interval for Periodic Discovery.
	DefaultPeriodicSyncInterval = time.Minute * 30

	// DefaultPeriodicDeleteInterval Default interval for Periodic delete.
	DefaultPeriodicDeleteInterval = time.Minute * 30

	// DefaultCacheResyncInterval Default interval for cache resync.
	DefaultCacheResyncInterval = time.Minute * 5
)
