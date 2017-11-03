package policy

import (
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/distributor"
)

// Policy represents an algorithm used to decide which collectorset
// should be used.
type Policy struct {
	DistributionStrategy distributor.DistributionStrategy
}

// Validated validates a policy.
// TODO: Validated should rank a policy based on the request on a scale that
// can be used to compare requests that match multiple polices.
func (p *Policy) Validated(req *api.CollectorIDRequest) bool {
	return true
}
