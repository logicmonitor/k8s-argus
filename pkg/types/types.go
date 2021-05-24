package types

// go:generate mockgen -destination=../mocks/mock_types.go -package=mocks github.com/logicmonitor/k8s-argus/pkg/types LMFacade,Watcher,DeviceManager,DeviceMapper,DeviceBuilder

import (
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"k8s.io/client-go/kubernetes"
	toolscache "k8s.io/client-go/tools/cache"
)

const defaultWorkerRetryLimit = 2

// Base is a struct for embedding
type Base struct {
	LMClient  *client.LMSdkGo
	K8sClient *kubernetes.Clientset
	Config    *config.Config
}

// WConfig worker configuration
type WConfig struct {
	ID         enums.ResourceType
	Channels   map[string]chan ICommand
	RetryLimit int
}

// NewHTTPWConfig new
func NewHTTPWConfig(rt enums.ResourceType) *WConfig {
	ch := make(chan ICommand)

	return &WConfig{
		Channels: map[string]chan ICommand{
			"GET":    ch,
			"POST":   ch,
			"DELETE": ch,
			"PUT":    ch,
			"PATCH":  ch,
		},
		RetryLimit: defaultWorkerRetryLimit,
		ID:         rt,
	}
}

// GetConfig returns reference to itself. impl here to avoid duplication everywhere
func (wc *WConfig) GetConfig() *WConfig {
	return wc
}

// GetChannel Get channel for mentioned command
func (wc *WConfig) GetChannel(command ICommand) chan ICommand {
	// Convert to switch case when adding new if conditions
	if command, ok := command.(IHTTPCommand); ok {
		m := command.(IHTTPCommand).GetMethod()

		return wc.Channels[m]
	}

	return nil
}

// Watcher is the LogicMonitor Watcher interface.
type Watcher interface {
	AddFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, []DeviceOption)
	DeleteFunc() func(interface{})
	UpdateFunc() func(oldObj, newObj interface{})
	GetConfig() *WConfig
	ResourceType() enums.ResourceType
}

// ResourceWatcher is the LogicMonitor Watcher interface.
type ResourceWatcher interface {
	GetConfig() *WConfig
	ResourceType() enums.ResourceType
}

// WatcherConfigurer is the LogicMonitor Watcher interface.
type WatcherConfigurer interface {
	AddFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, DeviceBuilder) ([]DeviceOption, error)
	UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, DeviceBuilder) ([]DeviceOption, bool, error)
	DeleteFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}) []DeviceOption
}

// ResourceCache cache
type ResourceCache interface {
	Run()
	Set(cache.ResourceName, cache.ResourceMeta) bool
	Exists(*lmctx.LMContext, cache.ResourceName, string) (cache.ResourceMeta, bool)
	Get(lctx *lmctx.LMContext, name cache.ResourceName) ([]cache.ResourceMeta, bool)
	Unset(cache.ResourceName, string) bool
	Load() error
	Save() error
	List() []cache.IterItem
	UnsetLMID(enums.ResourceType, int32) bool
}

// DeviceManager is an interface that describes how resources in Kubernetes
// are mapped into LogicMonitor as devices.
type DeviceManager interface {
	DeviceMapper
	DeviceBuilder
	Actions
	GetResourceCache() ResourceCache
}

// Actions actions
type Actions interface {
	// AddFunc wrapper
	AddFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...DeviceOption)
	UpdateFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ...DeviceOption)
	DeleteFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...DeviceOption)
}

// DeviceMapper is the interface responsible for mapping a Kubernetes resource to
// a LogicMonitor device.
type DeviceMapper interface {
	// FindByDisplayName searches for a device by it's display name. It will
	// returns a device if and only if
	// one device was found, and
	// returns nil otherwise.
	FindByDisplayName(*lmctx.LMContext, enums.ResourceType, string) (*models.Device, error)
	// Add adds a device to a LogicMonitor account.
	Add(*lmctx.LMContext, enums.ResourceType, interface{}, ...DeviceOption) (*models.Device, error)
	// DeleteByID deletes a device by device ID.
	DeleteByID(*lmctx.LMContext, enums.ResourceType, int32) error
}

// DeviceOption is the function definition for the functional options pattern.
type DeviceOption func(*models.Device)

// DeviceBuilder is the interface responsible for building a device struct.
type DeviceBuilder interface {
	// Name sets the device name.
	Name(string) DeviceOption
	// DisplayName sets the device name.
	DisplayName(string) DeviceOption
	// CollectorID sets the preferred collector ID for the device.
	CollectorID(int32) DeviceOption
	// SystemCategory sets the system.categories property on the device.
	SystemCategory(string, enums.BuilderAction) DeviceOption
	// ResourceLabels sets custom properties for the device
	ResourceLabels(map[string]string) DeviceOption
	// Auto adds an auto property to the device.
	Auto(string, string) DeviceOption
	// System adds a system property to the device.
	System(string, string) DeviceOption
	// Custom adds a custom property to the device.
	Custom(string, string) DeviceOption
	// DeletedOn adds kubernetes.resourceDeletedOn property to the device.
	DeletedOn(time.Time) DeviceOption
}

// UpdateFilter is a boolean function to run predicate and
// returns boolean value
type UpdateFilter func() bool

// ExecRequest function type to point to execute function
type ExecRequest func() (interface{}, error)

// ParseErrResp function signature to parse error response
type ParseErrResp func(error) *models.ErrorResponse

// LMExecutor All the
type LMExecutor interface {
	AddDevice(*lm.AddDeviceParams) ExecRequest
	AddDeviceErrResp(error) *models.ErrorResponse

	UpdateDevice(*lm.UpdateDeviceParams) ExecRequest
	UpdateDeviceErrResp(error) *models.ErrorResponse

	GetDeviceByID(params *lm.GetDeviceByIDParams) ExecRequest
	GetDeviceByIDErrResp(error) *models.ErrorResponse

	UpdateDevicePropertyByName(*lm.UpdateDevicePropertyByNameParams) ExecRequest
	UpdateDevicePropertyErrResp(error) *models.ErrorResponse

	GetDeviceList(*lm.GetDeviceListParams) ExecRequest
	GetDeviceListErrResp(error) *models.ErrorResponse

	PatchDevice(*lm.PatchDeviceParams) ExecRequest
	PatchDeviceErrResp(error) *models.ErrorResponse

	DeleteDeviceByID(*lm.DeleteDeviceByIDParams) ExecRequest
	DeleteDeviceByIDErrResp(error) *models.ErrorResponse

	GetImmediateDeviceListByDeviceGroupID(*lm.GetImmediateDeviceListByDeviceGroupIDParams) ExecRequest
	GetImmediateDeviceListByDeviceGroupIDErrResp(error) *models.ErrorResponse
}

// WorkerResponse wraps response and error
type WorkerResponse struct {
	Response interface{}
	Error    error
}

// Worker worker interface to provide interface method
type Worker interface {
	Run()
	GetConfig() *WConfig
}

// HTTPWorker specific worker to handle http requests
type HTTPWorker interface {
	Worker
	// TODO: Headers need to intercept for rate limiting the requests and for backoff
	// GetHeaders(interface{}) map[string]interface{}
}

// type GetHeaders func(response interface{}) (interface{}, error)

// ICommand based command interface
type ICommand interface {
	Execute() (interface{}, error)
	LMContext() *lmctx.LMContext
}

// Responder interface to indicate response can be sent back
type Responder interface {
	SetResponseChannel(chan *WorkerResponse)
	GetResponseChannel() chan *WorkerResponse
}

// Command base command
type Command struct {
	LMCtx       *lmctx.LMContext
	ExecFun     ExecRequest
	RespChannel chan *WorkerResponse
}

// Execute command execute
func (c *Command) Execute() (interface{}, error) {
	return c.ExecFun()
}

// LMContext returns LMContext object from command
func (c *Command) LMContext() *lmctx.LMContext {
	return c.LMCtx
}

// SetResponseChannel sets response channel  into command to send response back
func (c *Command) SetResponseChannel(rch chan *WorkerResponse) {
	c.RespChannel = rch
}

// GetResponseChannel returns response channel to send response
func (c *Command) GetResponseChannel() chan *WorkerResponse {
	return c.RespChannel
}

// IHTTPCommand Http command interface
type IHTTPCommand interface {
	// GetMethod Get Http method
	GetMethod() string
	// GetCategory Get rest api category
	GetCategory() string
}

// LMHCErrParse function to parse error response
type LMHCErrParse struct {
	ParseErrResp ParseErrResp
}

// ParseErrResponse executes parse error response function
func (lhp *LMHCErrParse) ParseErrResponse(err error) *models.ErrorResponse {
	return lhp.ParseErrResp(err)
}

// HTTPCommand extended Command
type HTTPCommand struct {
	*Command
	*LMHCErrParse
	Method   string
	Category string
	IsGlobal bool
	// GetHeaderFun GetHeaders
}

// GetMethod Get Http method
func (hc *HTTPCommand) GetMethod() string {
	return hc.Method
}

// GetCategory Get rest api category
func (hc *HTTPCommand) GetCategory() string {
	return hc.Category
}

// LMHCErrParser methods specific to lm sdk
type LMHCErrParser interface {
	ParseErrResponse(err error) *models.ErrorResponse
}

// LMFacade public interface others to interact with
type LMFacade interface {
	// Async
	// Send(command ICommand)

	// SendReceive sync api call
	SendReceive(*lmctx.LMContext, enums.ResourceType, ICommand) (interface{}, error)
	// RegisterWorker registers worker to facade client to put command objects on channel
	RegisterWorker(enums.ResourceType, Worker) (bool, error)
}

// RateLimitUpdateRequest struct to send new rate limits received from server to manager
type RateLimitUpdateRequest struct {
	IsGlobal bool
	Worker   enums.ResourceType
	Category string
	Method   string
	Limit    int64
	Window   int
}

// WorkerRateLimitsUpdate struct to send new rate limits received from server to manager
type WorkerRateLimitsUpdate struct {
	Category string
	Method   string
	Limit    int64
	Window   int
}

// RateLimitManager interface for rate limit manager
type RateLimitManager interface {
	// GetUpdateRequestChannel channel to send new limits to rate limit manager
	GetUpdateRequestChannel() chan RateLimitUpdateRequest
	// GetRateLimitConfig sends config for requested resource
	GetRateLimitConfig(resource string) map[string]int
	// RegisterWorkerNotifyChannel register channel to send updates to workers
	RegisterWorkerNotifyChannel(resource string, ch chan WorkerRateLimitsUpdate) (bool, error)
}

// DeviceExists error when device is already present in LM
type DeviceExists struct{}

// Error implements error interface
func (err DeviceExists) Error() string {
	return "device already present, ignoring add event"
}

// GetCollectorIDError error when device is already present in LM
type GetCollectorIDError struct {
	Err error
}

// Error implements error interface
func (err GetCollectorIDError) Error() string {
	return fmt.Sprintf("could not get collector id: %s", err.Err)
}

// ControllerInitSyncStateHolder struct to hold state
type ControllerInitSyncStateHolder struct {
	Controller toolscache.Controller
	hasSynced  bool
	mu         sync.RWMutex
}

// NewControllerInitSyncStateHolder create
func NewControllerInitSyncStateHolder(controller toolscache.Controller) ControllerInitSyncStateHolder {
	return ControllerInitSyncStateHolder{
		Controller: controller,
		hasSynced:  false,
		mu:         sync.RWMutex{},
	}
}

// HasSynced retrieve current state
func (stateHolder *ControllerInitSyncStateHolder) HasSynced() bool {
	stateHolder.mu.RLock()
	defer stateHolder.mu.RUnlock()

	return stateHolder.hasSynced
}

func (stateHolder *ControllerInitSyncStateHolder) markSynced() {
	stateHolder.mu.Lock()
	defer stateHolder.mu.Unlock()
	stateHolder.hasSynced = true
}

// Run starts watching on controller status and stops when it completes initial sync
func (stateHolder *ControllerInitSyncStateHolder) Run() {
	go func() {
		for {
			hs := stateHolder.Controller.HasSynced()
			if hs {
				stateHolder.markSynced()

				break
			}
			// May go in back to back loop and stop Controller from running update so sleeping for a millis to give chance to other goroutines
			time.Sleep(time.Millisecond)
		}
	}()
}
