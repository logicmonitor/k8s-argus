package storage

import (
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/policy"
)

// CollectorSetName is a type used for keys in a PolicyMap.
type CollectorSetName = string

// PolicyStore is an interface for handling the policy cache.
type PolicyStore interface {
	GetPolicy(CollectorSetName) (*policy.Policy, bool)
	SetPolicy(CollectorSetName, *policy.Policy) error
	DeletePolicy(CollectorSetName) error
	IterPolicies() <-chan *policy.Policy
}

// Storage is an interface used for storing cached state.
type Storage interface {
	PolicyStore
}
