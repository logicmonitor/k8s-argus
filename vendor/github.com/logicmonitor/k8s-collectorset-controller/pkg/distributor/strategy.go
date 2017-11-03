package distributor

import "github.com/logicmonitor/k8s-collectorset-controller/api"

// DistributionStrategy represents an algorithm used to decide which collector
// in a set should be used to monitor a given device.
type DistributionStrategy interface {
	ID(*api.CollectorIDRequest) (*api.CollectorIDReply, error)
	SetIDs([]int32) error
}
