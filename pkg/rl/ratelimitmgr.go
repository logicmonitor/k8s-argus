package ratelimiter

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	deployments = "deployments"
	pods        = "pods"
	services    = "services"
	nodes       = "nodes"
)

var (
	rlm   = NewManager()
	mutex sync.Mutex
)

func readConfig() *Policies {
	configBytes, err := ioutil.ReadFile("/etc/argus/rl-policy.yaml")
	if err != nil {
		log.Fatalf("Failed to read rl policy config file: /etc/argus/rl-policy.yaml")
	}
	log.Debugf("rl policy raw: %s", configBytes)
	m := &Policies{}
	err = yaml.Unmarshal(configBytes, m)
	if err != nil {
		log.Fatalf("Couldn't parse rl-policy.yaml file")
	}
	log.Infof("Policies read: %v", m)
	return m
}

func init() {
	rlm.SetPolicies(readConfig())
	rlm.Run()
}

// VerbLimits limits of http verb
type VerbLimits struct {
	POD  int64 `yaml:"pod"`
	DEP  int64 `yaml:"dep"`
	SVC  int64 `yaml:"svc"`
	NODE int64 `yaml:"node"`
}

// Get returns values from verblimit for mentioned resource
func (vl *VerbLimits) Get(resource string) int64 {
	switch resource {
	case pods:
		return vl.POD
	case deployments:
		return vl.DEP
	case services:
		return vl.SVC
	case nodes:
		return vl.NODE
	}
	return -1
}

// Set sets value to mentioned resource with the new value
func (vl *VerbLimits) Set(resource string, limit int64) {
	switch resource {
	case pods:
		vl.POD = limit
	case deployments:
		vl.DEP = limit
	case services:
		vl.SVC = limit
	case nodes:
		vl.NODE = limit
	}
}

// VerbPolicy limits of http verb
type VerbPolicy struct {
	POD  int64 `yaml:"pod"`
	DEP  int64 `yaml:"dep"`
	SVC  int64 `yaml:"svc"`
	NODE int64 `yaml:"node"`
}

// Get return policy for mentioned resource
func (vl VerbPolicy) Get(resource string) int64 {
	switch resource {
	case pods:
		return vl.POD
	case deployments:
		return vl.DEP
	case services:
		return vl.SVC
	case nodes:
		return vl.NODE
	}
	return -1
}

// APIResource resource category like device, devicegroup, etc.
type APIResource struct {
	GET    *VerbLimits `yaml:"get"`
	POST   *VerbLimits `yaml:"post"`
	PUT    *VerbLimits `yaml:"put"`
	PATCH  *VerbLimits `yaml:"patch"`
	DELETE *VerbLimits `yaml:"delete"`
}

// Get returns verblimits for mentioned http verb
func (apires *APIResource) Get(method string) *VerbLimits {
	switch method {
	case http.MethodGet:
		return apires.GET
	case http.MethodPost:
		return apires.POST
	case http.MethodPut:
		return apires.PUT
	case http.MethodPatch:
		return apires.PATCH
	case http.MethodDelete:
		return apires.DELETE
	}
	return nil
}

// APIPolicy resource category like device, devicegroup, etc.
type APIPolicy struct {
	GET    VerbPolicy `yaml:"get"`
	POST   VerbPolicy `yaml:"post"`
	PUT    VerbPolicy `yaml:"put"`
	PATCH  VerbPolicy `yaml:"patch"`
	DELETE VerbPolicy `yaml:"delete"`
}

// Get returns verbpolicy for mentioned http verb
func (apires APIPolicy) Get(method string) VerbPolicy {
	switch method {
	case http.MethodGet:
		return apires.GET
	case http.MethodPost:
		return apires.POST
	case http.MethodPut:
		return apires.PUT
	case http.MethodPatch:
		return apires.PATCH
	case http.MethodDelete:
		return apires.DELETE
	}
	return VerbPolicy{}
}

//Policies rate limit policies
type Policies struct {
	Device APIPolicy `yaml:"device"`
}

// Get rturns policy for mentioned api resource category
func (p Policies) Get(resource string) APIPolicy {
	switch resource {
	case "device":
		return p.Device
	}
	return APIPolicy{}
}

// Manager struct to hold manager entities
type Manager struct {
	Policies                Policies
	WorkerBroadcastChannels map[string]chan types.WorkerRateLimitsUpdate
	CurrentLimits           map[string]*APIResource
	UpReqChan               chan types.RateLimitUpdateRequest
}

// NewManager creates a object
func NewManager() *Manager {
	return &Manager{
		WorkerBroadcastChannels: make(map[string]chan types.WorkerRateLimitsUpdate),
		CurrentLimits:           make(map[string]*APIResource),
		UpReqChan:               make(chan types.RateLimitUpdateRequest),
	}
}

// GetUpdateRequestChannel channel to send new limits to rate limit manager
func (m *Manager) GetUpdateRequestChannel() chan types.RateLimitUpdateRequest {
	return m.UpReqChan
}

// GetUpdateRequestChannel channel to send new limits to rate limit manager
func GetUpdateRequestChannel() chan types.RateLimitUpdateRequest {
	return rlm.UpReqChan
}

// RegisterWorkerNotifyChannel registers worker to receive updates
func (m *Manager) RegisterWorkerNotifyChannel(resource string, ch chan types.WorkerRateLimitsUpdate) {
	m.WorkerBroadcastChannels[resource] = ch
}

// RegisterWorkerNotifyChannel registers worker to receive updates
func RegisterWorkerNotifyChannel(resource string, ch chan types.WorkerRateLimitsUpdate) {
	rlm.RegisterWorkerNotifyChannel(resource, ch)
}

// SetPolicies sets policies
func (m *Manager) SetPolicies(policies *Policies) {
	m.Policies = *policies
}

func (m *Manager) saveLimits(request types.RateLimitUpdateRequest) {
	apires, ok := m.CurrentLimits[request.Category]
	if !ok {
		apires = m.initNewCurrentLimit()
		m.CurrentLimits[request.Category] = apires
	}
	log.Infof("%v", apires)
	httpVerbLimits := apires.Get(request.Method)
	httpVerbLimits.Set(request.Worker, request.Limit)
}

func (m *Manager) hasDelta(request types.RateLimitUpdateRequest) bool {
	apires, ok := m.CurrentLimits[request.Category]
	if !ok {
		return true
	}
	if apires.Get(request.Method).Get(request.Worker) != request.Limit {
		return true
	}
	return false
}
func (m *Manager) initNewCurrentLimit() *APIResource {
	return &APIResource{
		GET: &VerbLimits{
			POD:  0,
			DEP:  0,
			SVC:  0,
			NODE: 0,
		},
		PUT: &VerbLimits{
			POD:  0,
			DEP:  0,
			SVC:  0,
			NODE: 0,
		},
		POST: &VerbLimits{
			POD:  0,
			DEP:  0,
			SVC:  0,
			NODE: 0,
		},
		PATCH: &VerbLimits{
			POD:  0,
			DEP:  0,
			SVC:  0,
			NODE: 0,
		},
		DELETE: &VerbLimits{
			POD:  0,
			DEP:  0,
			SVC:  0,
			NODE: 0,
		},
	}
}

func (m *Manager) distributeWorkerLimits(request types.RateLimitUpdateRequest) {
	vl := m.Policies.Get(request.Category).Get(request.Method)
	log.Debugf("reform policies : %v", m.Policies)
	log.Debugf("reform limits for : %v", vl)
	v := reflect.ValueOf(vl)
	typeOfS := v.Type()

	/*
		Distribute limits among workers according to configured percentage limits for those workers
		for ex:
		POD: 70%
		NODES: 10%
		Say, rate limit threshold values is 500, then 350 limit will be given to POD worker and 50 will be given to NODES.
	*/
	for i := 0; i < v.NumField(); i++ {
		log.Infof("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		rlch := types.WorkerRateLimitsUpdate{
			Category: request.Category,
			Method:   request.Method,
			// divide by 100 since configured value is in percentage
			Limit:  request.Limit * v.Field(i).Interface().(int64) / 100,
			Window: request.Window,
		}
		switch typeOfS.Field(i).Name {
		case "POD":
			m.WorkerBroadcastChannels["pods"] <- rlch
		case "DEP":
			m.WorkerBroadcastChannels["deployments"] <- rlch
		case "SVC":
			m.WorkerBroadcastChannels["services"] <- rlch
		case "NODE":
			m.WorkerBroadcastChannels["nodes"] <- rlch
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

// Run starts ratelimit manager
func (m *Manager) Run() {
	log.Debugf("policies initialised: %v", m.Policies)
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case req := <-m.UpReqChan:
				log.Infof("Update Limit Request received to ratelimiter %v", req)
				m.handleRequest(req)
			case <-ticker.C:
				log.Debugf("Rate limit manager is on standby mode %+v", m.CurrentLimits)
			}
		}
	}()
}
