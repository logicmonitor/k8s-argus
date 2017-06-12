package constants

const (
	ClusterCategory = "KubernetesCluster"
	NodeCategory    = "KubernetesNode"
	ServiceCategory = "KubernetesService"
	PodCategory     = "KubernetesPod"
)

var (
	// Version is set at build time.
	Version string
)
