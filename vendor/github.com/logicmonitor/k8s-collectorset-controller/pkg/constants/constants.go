package constants

var (
	// Version is the Argus version and is set at build time.
	Version string
)

const (
	// UserAgentBase is the base string for the User-Agent HTTP header.
	UserAgentBase = "LogicMonitor CollectorSet/"
)

const (
	// AccessID is the environment variable name to lookup for the LogicMonitor access ID.
	AccessID = "ARGUS_ACCESS_ID"
	// AccessKey is the environment variable name to lookup for the LogicMonitor access key.
	AccessKey = "ARGUS_ACCESS_KEY"
	// Account is the environment variable name to lookup for the LogicMonitor account.
	Account = "ARGUS_ACCOUNT"
)

const (
	// ArgusSecretName is the service account name with the proper RBAC policies to allow a collector to monitor the cluster.
	ArgusSecretName = "argus"
	// CollectorServiceAccountName is the service account name with the proper RBAC policies to allow a collector to monitor the cluster.
	CollectorServiceAccountName = "collector"
	// HealthServerServiceName is the gRPC service name for the health checks.
	HealthServerServiceName = "grpc.health.v1.Health"
)
