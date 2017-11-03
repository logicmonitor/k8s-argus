package inmem

import (
	"sync"

	"github.com/logicmonitor/k8s-collectorset-controller/pkg/policy"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/storage"
)

// PolicyMap is an in-memory data store that is used to cache CollectorSet
// policies.
type PolicyMap map[storage.CollectorSetName]*policy.Policy

// InMem represents the in-memory storage by implementing the storage.Storage
// interface.
type InMem struct {
	sync.RWMutex
	policies  PolicyMap
	countChan chan<- int
}

// New instantiates and returns an InMem storage.
func New(count chan<- int) *InMem {
	return &InMem{
		policies:  PolicyMap{},
		countChan: count,
	}
}

// GetPolicy provides thread-safe reads from a PolicyMap.
func (m *InMem) GetPolicy(name storage.CollectorSetName) (*policy.Policy, bool) {
	m.Lock()
	defer m.Unlock()

	p, ok := m.policies[name]

	return p, ok
}

// SetPolicy provides thread-safe writes to a PolicyMap.
func (m *InMem) SetPolicy(name storage.CollectorSetName, policy *policy.Policy) error {
	m.Lock()
	defer m.Unlock()

	m.policies[name] = policy
	m.countChan <- len(m.policies)

	return nil
}

// DeletePolicy provides thread-safe writes to a PolicyMap.
func (m *InMem) DeletePolicy(name storage.CollectorSetName) error {
	m.Lock()
	defer m.Unlock()

	delete(m.policies, name)
	m.countChan <- len(m.policies)

	return nil
}

// IterPolicies provides thread-safe iteration of a PolicyMap.
func (m *InMem) IterPolicies() <-chan *policy.Policy {
	c := make(chan *policy.Policy)

	f := func() {
		m.Lock()
		defer m.Unlock()

		for _, v := range m.policies {
			c <- v
		}
		close(c)
	}
	go f()

	return c
}
