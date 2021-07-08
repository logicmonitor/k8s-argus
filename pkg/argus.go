package argus

import (
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/etcd"
	"github.com/logicmonitor/k8s-argus/pkg/eventprocessor"
	"github.com/logicmonitor/k8s-argus/pkg/filters"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/resource/builder"
	"github.com/logicmonitor/k8s-argus/pkg/resourcecache"
	"github.com/logicmonitor/k8s-argus/pkg/sync"
	"github.com/logicmonitor/k8s-argus/pkg/tree"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/k8s-argus/pkg/watch/namespace"
	"github.com/logicmonitor/k8s-argus/pkg/watch/resourcewatcher"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

// Argus represents the Argus cli.
type Argus struct {
	*types.LMRequester
	types.ResourceManager
	types.ResourceCache
	Watchers               []types.ResourceWatcher
	controllerStateHolders map[enums.ResourceType]*types.ControllerInitSyncStateHolder
	NSWatcher              *namespace.OldWatcher
	RunnerFacade           eventprocessor.RunnerFacade
}

func (a *Argus) Init() error {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"argus": "init"}))
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	var resourceGroupTree *types.ResourceGroupTree
	if conf.EnableNewResourceTree {
		resourceGroupTree, err = tree.GetResourceGroupTree2(lctx, a.ResourceManager, a.LMRequester)
	} else {
		resourceGroupTree, err = tree.GetResourceGroupTree(lctx, a.ResourceManager, a.LMRequester)
	}
	if err := a.CreateResourceGroupTree(lctx, resourceGroupTree, true); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	// init sync to delete the non-exist resource resources through LogicMonitor API
	initSyncer := sync.InitSyncer{
		LMRequester:     a.LMRequester,
		ResourceManager: a.ResourceManager,
	}

	lctx2 := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "init-sync"}))
	initSyncer.Sync(lctx2)

	// periodically delete the non-exist resource resources through logicmonitor API based on specified time interval.
	initSyncer.RunPeriodicSync()

	err = discoverETCDNodes(a.ResourceManager)
	if err != nil {
		return err
	}
	return nil
}

// NewArgus instantiates and returns argus.
// nolint: cyclop
func NewArgus(lmrequester *types.LMRequester, resourceManager types.ResourceManager, resourceCache *resourcecache.ResourceCache) (*Argus, error) {
	return &Argus{ // nolint: exhaustivestruct
		LMRequester:            lmrequester,
		ResourceManager:        resourceManager,
		ResourceCache:          resourceCache,
		controllerStateHolders: make(map[enums.ResourceType]*types.ControllerInitSyncStateHolder),
	}, nil
}

func discoverETCDNodes(resourceManager types.ResourceManager) error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	if conf.EtcdDiscoveryToken != "" {
		etcdController := etcd.Controller{
			ResourceManager: resourceManager,
		}
		_, err := etcdController.DiscoverByToken()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Argus) CreateWatchers(lctx *lmctx.LMContext) error {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	m := make(map[enums.ResourceType]*struct{})
	for _, d := range conf.DisableResourceMonitoring {
		m[d] = nil
	}
	for _, rt := range enums.ALLResourceTypes {
		// These need special handling
		if rt == enums.Namespaces || rt == enums.ETCD {
			continue
		}
		if _, ok := m[rt]; ok {
			log.Warnf("Resource %s is being disabled for monitoring", rt.String())
			continue
		}
		a.Watchers = append(a.Watchers, &resourcewatcher.Watcher{Resource: rt})
	}
	a.NSWatcher = namespace.NewOldWatcher(a.ResourceManager, a.ResourceCache, a.LMRequester)
	return nil
}

// Watch watches the API for events.
func (a *Argus) Watch(lctx *lmctx.LMContext) error {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	syncInterval := *conf.Intervals.PeriodicSyncInterval
	log.Debugf("Starting watchers")
	b := &builder.Builder{}

	nsRT, controller := a.RunNSWatcher(syncInterval)
	log.Debugf("Starting ns watcher of %v", nsRT.String())
	stop := make(chan struct{})
	go controller.Run(stop)

	for _, w := range a.Watchers {
		rt := w.ResourceType()
		// TODO: has permission and check for enabled flag in case if user wants to avoid all resource of specific type
		//  earlier all resources used to ignore from filter config but still it used to put pressure on k8s api-server to unnecessary polls
		if !permission.HasPermissions(rt) {
			log.Warnf("Have no permission for resource %s", rt.String())

			continue
		}
		watchlist := cache.NewListWatchFromClient(util.GetK8sRESTClient(config.GetClientSet(), rt.K8SAPIVersion()), rt.String(), corev1.NamespaceAll, fields.Everything())
		clientState, controller := a.createNewInformer(watchlist, rt, syncInterval, b)
		go watchForFilterRuleChange(rt, clientState)
		log.Debugf("Starting watcher of %s", rt)
		stop := make(chan struct{})
		stateHolder := types.NewControllerInitSyncStateHolder(controller)
		stateHolder.Run()
		a.controllerStateHolders[rt] = &stateHolder
		go controller.Run(stop)
	}
	return nil
}

func watchForFilterRuleChange(rt enums.ResourceType, clientstate cache.Store) {
	func() {
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"filter_hook": rt.String()}))
		log := lmlog.Logger(lctx)
		filters.AddFilterHook(filters.FilterHook{
			Hook: func(c *filters.Config, c2 *filters.Config) {
				err := clientstate.Resync()
				if err != nil {
					log.Errorf("Failed to run graceful resync on exclusion rules change: %s", err)
				}
				log.Infof("Resync initiatiated on exclusion rules change")
			},
			Predicate: func(c *filters.Config, c2 *filters.Config) bool {
				return true
			},
		})
	}()
}

func (a *Argus) createNewInformer(watchlist cache.ListerWatcher, rt enums.ResourceType, syncInterval time.Duration, b *builder.Builder) (cache.Store, cache.Controller) {
	return cache.NewInformer(
		watchlist,
		rt.K8SObjectType(),
		syncInterval,
		cache.FilteringResourceEventHandler{
			FilterFunc: a.genericObjectFilterFunc(),
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: resourcewatcher.AddFuncDispatcher(
					a.RunnerFacade,
					resourcewatcher.AddFuncWithExclude(
						resourcewatcher.PreprocessAddEventForOldUID(
							a.ResourceManager.GetResourceCache(),
							a.ResourceManager.DeleteFunc(),
							b,
							resourcewatcher.AddOrUpdateFunc(
								a.controllerStateHolders,
								b.AddFuncWithDefaults(
									resourcewatcher.WatcherConfigurer(rt),
									a.ResourceManager,
								),
								b.UpdateFuncWithDefaults(
									resourcewatcher.UpsertBasedOnCache(
										a.ResourceManager.GetResourceCache(),
										resourcewatcher.WatcherConfigurer(rt),
										a.ResourceManager,
										b,
									),
								),
							),
						),
						b.DeleteFuncWithDefaults(resourcewatcher.WatcherConfigurer(rt), a.ResourceManager.DeleteFunc()),
					),
				),
				UpdateFunc: resourcewatcher.UpdateFuncDispatcher(
					a.RunnerFacade,
					resourcewatcher.UpdateFuncWithExclude(
						resourcewatcher.PreprocessUpdateEventForOldUID(
							a.ResourceManager.GetResourceCache(),
							a.ResourceManager.DeleteFunc(),
							b,
							b.UpdateFuncWithDefaults(
								resourcewatcher.UpsertBasedOnCache(
									a.ResourceManager.GetResourceCache(),
									resourcewatcher.WatcherConfigurer(rt),
									a.ResourceManager,
									b,
								),
							),
						),
						b.DeleteFuncWithDefaults(resourcewatcher.WatcherConfigurer(rt), a.ResourceManager.DeleteFunc()),
					),
				),
				DeleteFunc: resourcewatcher.DeleteFuncDispatcher(
					a.RunnerFacade,
					b.DeleteFuncWithDefaults(resourcewatcher.WatcherConfigurer(rt), a.ResourceManager.DeleteFunc()),
				),
			},
		},
	)
}

func (a *Argus) genericObjectFilterFunc() func(obj interface{}) bool {
	if conf, err := config.GetConfig(); err == nil {
		if conf.RegisterGenericFilter {
			return func(obj interface{}) bool {
				if rt, ok := resourcewatcher.InferResourceType(obj); ok {
					if meta := rt.ObjectMeta(obj); meta != nil {
						val := util.EvaluateExclusion(meta.Labels)
						if !val {
							logrus.Tracef("returning exclusion for: %s-%s", meta.Name, meta.Namespace)
						}
						return val
					}
					logrus.Tracef("cannot get ObjectMeta to run exclusion filter for ResourceType %s: %v", rt.String(), obj)

				} else {
					logrus.Tracef("cannot infer object type to run exclusion filter: %v", obj)
				}
				logrus.Tracef("returning true for %v", obj)
				return true
			}
		}
	}
	return func(obj interface{}) bool {
		return true
	}
}

func (a *Argus) RunNSWatcher(syncInterval time.Duration) (enums.ResourceType, cache.Controller) {
	rt := a.NSWatcher.ResourceType()
	// start ns watcher
	listWatch := cache.NewListWatchFromClient(util.GetK8sRESTClient(config.GetClientSet(), rt.K8SAPIVersion()), rt.String(), corev1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		listWatch,
		rt.K8SObjectType(),
		syncInterval,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    a.NSWatcher.AddFunc(),
			UpdateFunc: a.NSWatcher.UpdateFunc(),
			DeleteFunc: a.NSWatcher.DeleteFunc(),
		},
	)
	return rt, controller
}

func (a *Argus) CreateParallelRunners(lctx *lmctx.LMContext) error {
	runnerFacade, err := createRunnerFacade(lctx)
	if err != nil {
		return err
	}
	a.RunnerFacade = runnerFacade
	return nil
}
