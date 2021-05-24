package ratelimiter

import (
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	ratelimiter = NewManager()
	mutex       sync.Mutex
)

const rateLimitIdleNotifierTime = 1 * time.Minute

// Init package init block so that policies will be loaded on application start
func Init() {
	ratelimiter.SetPolicies(readConfig())
	ratelimiter.Run()
}

// RegisterWorkerNotifyChannel registers worker to receive updates
func RegisterWorkerNotifyChannel(resource enums.ResourceType, ch chan types.WorkerRateLimitsUpdate) {
	ratelimiter.RegisterWorkerNotifyChannel(resource, ch)
}

// GetUpdateRequestChannel channel to send new limits to rate limit manager
func GetUpdateRequestChannel() chan types.RateLimitUpdateRequest {
	return ratelimiter.UpReqChan
}

// Package internal structures and methods

// Policies rate limit policies
type Policies struct {
	Resource map[string]map[string]map[enums.ResourceType]int64 `yaml:"resource"`
	Global   map[string]map[string]int64                        `yaml:"global,omitempty"`
}

// RateLimits struct
type RateLimits struct {
	limits map[string]map[string]int64
	mu     sync.Mutex
}

// Store store
func (l *RateLimits) Store(category string, method string, limit int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	m, ok := l.limits[category]
	if !ok {
		m = make(map[string]int64)
	}
	m[method] = limit
	l.limits[category] = m
}

// Load load
func (l *RateLimits) Load(category string, method string) (int64, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if m, ok := l.limits[category]; ok {
		if val, ok := m[method]; ok {
			return val, true
		}
	}

	return -1, false
}

// Manager struct to hold manager entities
type Manager struct {
	Policies                Policies
	CurrentLimits           *RateLimits
	WorkerBroadcastChannels map[enums.ResourceType]chan types.WorkerRateLimitsUpdate
	UpReqChan               chan types.RateLimitUpdateRequest
}

// NewManager creates a object
func NewManager() *Manager {
	return &Manager{
		WorkerBroadcastChannels: make(map[enums.ResourceType]chan types.WorkerRateLimitsUpdate),
		CurrentLimits:           &RateLimits{limits: make(map[string]map[string]int64), mu: sync.Mutex{}},
		UpReqChan:               make(chan types.RateLimitUpdateRequest),
		Policies:                Policies{}, // nolint: exhaustivestruct
	}
}

// GetUpdateRequestChannel channel to send new limits to rate limit manager
func (m *Manager) GetUpdateRequestChannel() chan types.RateLimitUpdateRequest {
	return m.UpReqChan
}

// RegisterWorkerNotifyChannel registers worker to receive updates
func (m *Manager) RegisterWorkerNotifyChannel(resource enums.ResourceType, ch chan types.WorkerRateLimitsUpdate) {
	m.WorkerBroadcastChannels[resource] = ch
}

// SetPolicies sets policies
func (m *Manager) SetPolicies(policies *Policies) {
	m.Policies = *policies
}

func (m *Manager) saveLimits(request types.RateLimitUpdateRequest) {
	m.CurrentLimits.Store(request.Category, request.Method, request.Limit)
}

func (m *Manager) hasDelta(request types.RateLimitUpdateRequest) bool {
	if oldVal, ok := m.CurrentLimits.Load(request.Category, request.Method); ok {
		if oldVal == request.Limit {
			return false
		}
	}

	return true
}

func (m *Manager) distributeWorkerLimits(request types.RateLimitUpdateRequest) {
	if !request.IsGlobal {
		policies := m.Policies.Resource[request.Category][request.Method]
		for k, v := range policies {
			logrus.Infof("K: %v V: %v", k, v)
			workerRateLimitsUpdate := types.WorkerRateLimitsUpdate{
				Category: request.Category,
				Method:   request.Method,
				// divide by 100 since configured value is in percentage
				Limit:  request.Limit * v / 100, // nolint: gomnd
				Window: request.Window,
			}
			m.WorkerBroadcastChannels[k] <- workerRateLimitsUpdate
		}
	}
}

func (m *Manager) handleRequest(request types.RateLimitUpdateRequest) {
	mutex.Lock()
	defer mutex.Unlock()
	if !m.hasDelta(request) {
		return
	}
	m.saveLimits(request)
	m.distributeWorkerLimits(request)
}

// Run starts RateLimitManager
func (m *Manager) Run() {
	logrus.Debugf("policies initialised: %v", m.Policies)

	go func() {
		ticker := time.NewTicker(rateLimitIdleNotifierTime)

		for {
			select {
			case req := <-m.UpReqChan:
				logrus.Infof("Update Limit Request received to ratelimiter %v", req)
				m.handleRequest(req)
			case <-ticker.C:
				logrus.Debugf("Rate limit manager is on standby mode %+v", m.CurrentLimits)
			}
		}
	}()
}

func readConfig() *Policies {
	// configBytes, err := ioutil.ReadFile("/etc/argus/rl-policy.yaml")
	configBytes, err := config.GetWatchConfig("rl-policy.yaml")
	if err != nil {
		logrus.Fatalf("Failed to read rl policy config file: /etc/argus/rl-policy.yaml: %s", err)
	}
	logrus.Debugf("rl policy raw: %s", configBytes)
	m := &Policies{}
	err = yaml.Unmarshal([]byte(configBytes), m)
	if err != nil {
		logrus.Fatalf("Couldn't parse rl-policy.yaml file, %s", err)
	}
	logrus.Infof("Policies read: %v", m)

	return m
}
