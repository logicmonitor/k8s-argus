package types

// go:generate mockgen -destination=../mocks/mock_types.go -package=mocks github.com/logicmonitor/k8s-argus/pkg/types LMFacade,Watcher,ResourceManager,ResourceMapper,ResourceBuilder

import (
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	ID       int
	inCh     chan *WorkerCommand
	MaxRetry int
}

// NewWConfig new
func NewWConfig(id int) *WConfig {
	ch := make(chan *WorkerCommand)

	return &WConfig{
		inCh:     ch,
		MaxRetry: defaultWorkerRetryLimit,
		ID:       id,
	}
}

// GetConfig returns reference to itself. impl here to avoid duplication everywhere
func (wc *WConfig) GetConfig() *WConfig {
	return wc
}

// GetChannel Get channel for mentioned command
func (wc *WConfig) GetChannel() chan *WorkerCommand {
	return wc.inCh
}

// Watcher is the LogicMonitor Watcher interface.
type Watcher interface {
	AddFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, []ResourceOption)
	DeleteFunc() func(interface{})
	UpdateFunc() func(oldObj, newObj interface{})
	GetConfig() *WConfig
	ResourceType() enums.ResourceType
}

// ResourceWatcher is the LogicMonitor Watcher interface.
type ResourceWatcher interface {
	ResourceType() enums.ResourceType
}

// WatcherConfigurer is the LogicMonitor Watcher interface.
type WatcherConfigurer interface {
	AddFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, ResourceBuilder) ([]ResourceOption, error)
	UpdateFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ResourceBuilder) ([]ResourceOption, bool, error)
	DeleteFuncOptions() func(*lmctx.LMContext, enums.ResourceType, interface{}) []ResourceOption
}

// ResourceCache cache
type ResourceCache interface {
	Run()
	Set(*lmctx.LMContext, ResourceName, ResourceMeta) bool
	Exists(*lmctx.LMContext, ResourceName, string, bool) (ResourceMeta, bool)
	Get(lctx *lmctx.LMContext, name ResourceName) ([]ResourceMeta, bool)
	Unset(*lmctx.LMContext, ResourceName, string) bool
	Load(*lmctx.LMContext) error
	Save(*lmctx.LMContext) error
	List() []IterItem
	UnsetLMID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) bool
	SoftRefresh(*lmctx.LMContext, string)
	AddCacheHook(hook CacheHook)
	ListWithFilter(f func(k ResourceName, v ResourceMeta) bool) []IterItem
}

// ResourceGroupManager interface for resource group operations
type ResourceGroupManager interface {
	ResourceGroupBuilder
	CreateResourceGroupTree(lctx *lmctx.LMContext, tree *ResourceGroupTree, update bool) error
	GetResourceGroupByID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) (*models.DeviceGroup, error)
	DeleteResourceGroup(lctx *lmctx.LMContext, rt enums.ResourceType, id int32, deleteIfEmpty bool) error
}

// ResourceManager is an interface that describes how resources in Kubernetes
// are mapped into LogicMonitor as resources.
type ResourceManager interface {
	ResourceMapper
	ResourceBuilder
	Actions
	ResourceGroupManager
	GetResourceCache() ResourceCache
}

// Actions actions
type Actions interface {
	// AddFunc wrapper
	AddFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...ResourceOption) (*models.Device, error)
	UpdateFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ...ResourceOption) (*models.Device, error)
	DeleteFunc() func(*lmctx.LMContext, enums.ResourceType, interface{}, ...ResourceOption) error
}

// SyncUpdater methods for syncer only
type SyncUpdater interface {
	UpdateResourceByID(*lmctx.LMContext, enums.ResourceType, int32, ...ResourceOption) (*models.Device, error)
	DeleteResourceByID(*lmctx.LMContext, enums.ResourceType, int32) error
}

// ResourceMapper is the interface responsible for mapping a Kubernetes resource to
// a LogicMonitor resource.
type ResourceMapper interface {
	SyncUpdater
	// FindByDisplayName searches for a resource by it's display name. It will
	// returns a resource if and only if
	// one resource was found, and
	// returns nil otherwise.
	FindByDisplayName(*lmctx.LMContext, enums.ResourceType, string) (*models.Device, error)
	// DeleteByID deletes a resource by resource ID.
	DeleteByID(*lmctx.LMContext, enums.ResourceType, int32) error
}

// ResourceOption is the function definition for the functional options pattern.
type ResourceOption func(*models.Device)

// ResourceBuilder is the interface responsible for building a resource struct.
type ResourceBuilder interface {
	// Name sets the resource name.
	Name(string) ResourceOption
	// DisplayName sets the resource name.
	DisplayName(string) ResourceOption
	// CollectorID sets the preferred collector ID for the resource.
	CollectorID(int32) ResourceOption
	// SystemCategory sets the system.categories property on the resource.
	SystemCategory(string, enums.BuilderAction) ResourceOption
	// ResourceLabels sets custom properties for the resource
	ResourceLabels(map[string]string) ResourceOption
	// Auto adds an auto property to the resource.
	Auto(string, string) ResourceOption
	// System adds a system property to the resource.
	System(string, string) ResourceOption
	// Custom adds a custom property to the resource.
	Custom(string, string) ResourceOption
	// DeletedOn adds kubernetes.resourceDeletedOn property to the resource.
	DeletedOn(time.Time) ResourceOption

	GetMarkDeleteOptions(*lmctx.LMContext, enums.ResourceType, *metav1.PartialObjectMetadata) []ResourceOption
}

// AppliesToBuilder is an interface for building an appliesTo string.
type AppliesToBuilder interface {
	HasCategory(string) AppliesToBuilder
	Auto(string) AppliesToBuilder
	And() AppliesToBuilder
	OpenBracket() AppliesToBuilder
	TrimOrCloseBracket() AppliesToBuilder
	Custom(string) AppliesToBuilder
	Or() AppliesToBuilder
	Equals(string) AppliesToBuilder
	Exists(string) AppliesToBuilder
	Build() string
}

// PropertyBuilder is an interface for building properties
type PropertyBuilder interface {
	Add(string, string, bool) PropertyBuilder
	AddProperties([]config.PropOpts) PropertyBuilder
	Build([]*models.NameAndValue) []*models.NameAndValue
}

// ResourceGroupOption is the function definition for the functional options pattern.
type ResourceGroupOption func(group *models.DeviceGroup)

// ResourceGroupBuilder is the interface responsible for building a resource struct.
type ResourceGroupBuilder interface {
	// GroupName sets the resource name.
	GroupName(string) ResourceGroupOption
	// ParentID parent group id
	ParentID(int32) ResourceGroupOption
	DisableAlerting(bool) ResourceGroupOption
	AppliesTo(AppliesToBuilder) ResourceGroupOption
	CustomProperties(builder PropertyBuilder) ResourceGroupOption
}

// UpdateFilter is a boolean function to run predicate and
// returns boolean value
type UpdateFilter func() bool

// ExecRequest function type to point to execute function
type ExecRequest func() (interface{}, error)

// ParseErrResp function signature to parse error response
type ParseErrResp func(error) *models.ErrorResponse

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
	SendReceive(*lmctx.LMContext, *WorkerCommand) (interface{}, error)
	// RegisterWorker registers worker to facade client to put command objects on channel
	RegisterWorker(Worker) (bool, error)
	// UnregisterWorker registers worker to facade client to put command objects on channel
	UnregisterWorker(Worker) (bool, error)
	// Count registers worker to facade client to put command objects on channel
	Count() int
}

// RateLimits struct to send new rate limits received from server to manager
type RateLimits struct {
	Limit  int64
	Window int
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
	GetUpdateRequestChannel() chan RateLimits
	// GetRateLimitConfig sends config for requested resource
	GetRateLimitConfig(resource string) map[string]int
	// RegisterWorkerNotifyChannel register channel to send updates to workers
	RegisterWorkerNotifyChannel(resource string, ch chan WorkerRateLimitsUpdate) (bool, error)
}

// GetCollectorIDError error when resource is already present in LM
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

type (
	AddPreprocessFunc func(*lmctx.LMContext, enums.ResourceType, interface{})
	AddProcessFunc    func(*lmctx.LMContext, enums.ResourceType, interface{}, ...ResourceOption)
)

type (
	UpdatePreprocessFunc func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{})
	UpdateProcessFunc    func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, []ResourceOption, []ResourceOption)
)

type (
	DeletePreprocessFunc func(*lmctx.LMContext, enums.ResourceType, interface{})
	DeleteProcessFunc    func(*lmctx.LMContext, enums.ResourceType, interface{}, ...ResourceOption)
)

type (
	WatcherAddFunc    func(interface{})
	WatcherUpdateFunc func(interface{}, interface{})
	WatcherDeleteFunc func(interface{})
)

type (
	ExecAddFunc    func(*lmctx.LMContext, enums.ResourceType, interface{}, ...ResourceOption) (*models.Device, error)
	ExecUpdateFunc func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ...ResourceOption) (*models.Device, error)
	ExecDeleteFunc func(*lmctx.LMContext, enums.ResourceType, interface{}, ...ResourceOption) error
)

// LMRequester this is just to tiw facade and executor together, never mix or attach executor with facade
type LMRequester struct {
	LMFacade
	LMExecutor
}
