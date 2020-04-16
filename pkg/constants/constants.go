package constants

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
)

const (
	// LabelNodeRole is the label name used to specify a node's role in the cluster
	LabelNodeRole = "node-role.kubernetes.io/"
	// LabelCustomPropertyPrefix is the prefix to use for custom properties based of labels
	LabelCustomPropertyPrefix = "kubernetes.label."
	// LabelNullPlaceholder is the string used to represent null values in custom properties
	LabelNullPlaceholder = "null"
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
	// ServiceCategory is the system.category used to identity a Kubernetes Service resource type in LogicMonitor.
	ServiceCategory = "KubernetesService"
	// ServiceDeletedCategory is the system.category used to identity a deleted Kubernetes Service resource type in LogicMonitor.
	ServiceDeletedCategory = "KubernetesServiceDeleted"
	// DeploymentCategory is the system.category used to identity a Kubernetes Service resource type in LogicMonitor.
	DeploymentCategory = "KubernetesDeployment"
	// DeploymentDeletedCategory is the system.category used to identity a deleted Kubernetes Service resource type in LogicMonitor.
	DeploymentDeletedCategory = "KubernetesDeploymentDeleted"
	// PodCategory is the system.category used to identity the Kubernetes Pod resource type in LogicMonitor.
	PodCategory = "KubernetesPod"
	// PodDeletedCategory is the system.category used to identity a deleted Kubernetes Pod resource type in LogicMonitor.
	PodDeletedCategory = "KubernetesPodDeleted"
	// DeletedDeviceGroup is the name of the device group where deleted devices are optionally moved to.
	DeletedDeviceGroup = "_deleted"
	// ClusterDeviceGroupPrefix is the prefix for the top level cluster device group
	ClusterDeviceGroupPrefix = "Kubernetes Cluster: "
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
	// K8sDeviceType is the type value of the k8s device
	K8sDeviceType = 8
	// K8sSystemCategoriesPropertyKey is the key of the unique custom property kubernetes system categories
	K8sSystemCategoriesPropertyKey = "system.categories"
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
)
