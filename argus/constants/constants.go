package constants

const (
	// ClusterCategory is the system.category used to identity the Kubernetes cluster in LogicMonitor.
	ClusterCategory = "KubernetesCluster"
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
)
