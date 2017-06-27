package constants

const (
	ClusterCategory        = "KubernetesCluster"
	NodeCategory           = "KubernetesNode"
	NodeDeletedCategory    = "KubernetesNodeDeleted"
	ServiceCategory        = "KubernetesService"
	ServiceDeletedCategory = "KubernetesServiceDeleted"
	PodCategory            = "KubernetesPod"
	PodDeletedCategory     = "KubernetesPodDeleted"
)

var (
	// Version is set at build time.
	Version string
)
