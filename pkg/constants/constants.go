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
)
