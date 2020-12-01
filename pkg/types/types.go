package types

import (
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// Base is a struct for embedding
type Base struct {
	LMClient  *client.LMSdkGo
	K8sClient *kubernetes.Clientset
	Config    *config.Config
}

// WConfig worker configuration
type WConfig struct {
	ID             string
	MethodChannels map[string]chan ICommand
	RetryLimit     int
}

// GetConfig returns reference to itself. impl here to avoid duplication everywhere
func (wc *WConfig) GetConfig() *WConfig {
	return wc
}

// GetChannel Get channel for mentioned command
func (wc *WConfig) GetChannel(command ICommand) chan ICommand {
	switch command := command.(type) {
	case IHTTPCommand:
		m := command.(IHTTPCommand).GetMethod()
		return wc.MethodChannels[m]
	}
	return nil
}

// Watcher is the LogicMonitor Watcher interface.
type Watcher interface {
	APIVersion() string
	Enabled() bool
	Resource() string
	ObjType() runtime.Object
	AddFunc() func(obj interface{})
	DeleteFunc() func(obj interface{})
	UpdateFunc() func(oldObj, newObj interface{})
	GetConfig() *WConfig
}

// DeviceManager is an interface that describes how resources in Kubernetes
// are mapped into LogicMonitor as devices.
type DeviceManager interface {
	DeviceMapper
	DeviceBuilder
}

// DeviceMapper is the interface responsible for mapping a Kubernetes resource to
// a LogicMonitor device.
type DeviceMapper interface {
	// Config returns the Argus config.
	Config() *config.Config
	// FindByDisplayName searches for a device by it's display name. It will return a device if and only if
	// one device was found, and return nil otherwise.
	FindByDisplayName(*lmctx.LMContext, string, string) (*models.Device, error)
	// FindByDisplayNames searches for devices by the specified string by its display name. It will return the device list.
	FindByDisplayNames(*lmctx.LMContext, string, ...string) ([]*models.Device, error)
	// Add adds a device to a LogicMonitor account.
	Add(*lmctx.LMContext, string, map[string]string, ...DeviceOption) (*models.Device, error)
	// UpdateAndReplace updates a device using the 'replace' OpType.
	UpdateAndReplace(*lmctx.LMContext, string, *models.Device, ...DeviceOption) (*models.Device, error)
	// UpdateAndReplaceByDisplayName updates a device using the 'replace' OpType if and onlt if it does not already exist.
	UpdateAndReplaceByDisplayName(*lmctx.LMContext, string, string, string, UpdateFilter, map[string]string, ...DeviceOption) (*models.Device, error)
	// UpdateAndReplaceField updates a device using the 'replace' OpType for a
	// specific field of a device.
	UpdateAndReplaceField(*lmctx.LMContext, string, *models.Device, string, ...DeviceOption) (*models.Device, error)
	// UpdateAndReplaceFieldByDisplayName updates a device using the 'replace' OpType for a
	// specific field of a device.
	UpdateAndReplaceFieldByDisplayName(*lmctx.LMContext, string, string, string, string, ...DeviceOption) (*models.Device, error)
	// DeleteByID deletes a device by device ID.
	DeleteByID(*lmctx.LMContext, string, int32) error
	// DeleteByDisplayName deletes a device by device display name.
	DeleteByDisplayName(*lmctx.LMContext, string, string, string) error
	// GetDesiredDisplayName returns desired display name based on FullDisplayNameIncludeClusterName and FullDisplayNameIncludeNamespace properties.
	GetDesiredDisplayName(string, string, string) string
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
	// SystemCategories sets the system.categories property on the device.
	SystemCategories(string) DeviceOption
	// ResourceLabels sets custom properties for the device
	ResourceLabels(map[string]string) DeviceOption
	// Auto adds an auto property to the device.
	Auto(string, string) DeviceOption
	// System adds a system property to the device.
	System(string, string) DeviceOption
	// System adds a custom property to the device.
	Custom(string, string) DeviceOption
}

// UpdateFilter is a boolean function to run predicate and return boolean value
type UpdateFilter func() bool

// ExecRequest funnction type to point to execute fubction
type ExecRequest func() (interface{}, error)

// ParseErrResp function signature to parse error response
type ParseErrResp func(error) *models.ErrorResponse

// LMExecutor All the
type LMExecutor interface {
	AddDevice(*lm.AddDeviceParams) ExecRequest
	AddDeviceErrResp(error) *models.ErrorResponse

	UpdateDevice(*lm.UpdateDeviceParams) ExecRequest
	UpdateDeviceErrResp(error) *models.ErrorResponse

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

//Worker worker interface to provide interface method
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

// LMContext return LMContext object from command
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
	//Send(command ICommand)
	// sync
	SendReceive(*lmctx.LMContext, string, ICommand) (interface{}, error)
	RegisterWorker(string, Worker) (bool, error)
}

// RateLimitUpdateRequest struct to send new rate limits received from server to manager
type RateLimitUpdateRequest struct {
	Worker   string
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
