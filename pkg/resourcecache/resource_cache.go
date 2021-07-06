package resourcecache

import (
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var (
	// Version to differentiate cache flushed to configmap is with previous datastructures or not.
	// Whenever cache.ResourceName and cache.ResourceMeta struct modifies, increase the version so that it will easier to parse previously created CMs and convert it to newer data structure
	Version = "v1"

	// CMMaxBytes constant to set max threshold for byte size of configmap
	// 1MiB is the capacity of configmap data, but we are storing 80% of the data in a chunk
	CMMaxBytes = 1048576 * 80 / 100
)

const rateLimitBackoffTime = 10 * time.Second

// ResourceCache to maintain a resource cache to calculate delta between resource presence on server and on cluster
type ResourceCache struct {
	store         *Store
	rwm           sync.RWMutex
	rebuildMutex  sync.Mutex
	resyncPeriod  time.Duration
	flushTimeToCM time.Duration
	flushMU       sync.Mutex
	*types.LMRequester
	clusterGrpID int32
	stateLoaded  bool

	// On restart, first dump loaded from config map
	dumpID int64

	// soft refresh last time
	softRefreshLast map[string]time.Time
	softRefreshMu   sync.Mutex

	hooks   []types.CacheHook
	hookrwm sync.RWMutex
}

// NewResourceCache create new ResourceCache object
func NewResourceCache(facadeObj *types.LMRequester, rp time.Duration) *ResourceCache {
	resourceCache := &ResourceCache{
		store:         NewStore(),
		rwm:           sync.RWMutex{},
		resyncPeriod:  rp,
		flushTimeToCM: 1 * time.Minute,
		flushMU:       sync.Mutex{},
		LMRequester:   facadeObj,
		clusterGrpID:  -1,
		stateLoaded:   false,
		rebuildMutex:  sync.Mutex{},
	}

	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"resource_cache": "init"}))
	log := lmlog.Logger(lctx)
	if err := resourceCache.Load(lctx); err != nil {
		log.Warnf("ResourceCache is not loaded from its last state stored in configmaps, it may take more time to rebuild cache on first run. Error: %s", err)
	} else {
		resourceCache.stateLoaded = true
		log.Info("ResourceCache loaded from configmaps")
	}

	return resourceCache
}

func (rc *ResourceCache) AddCacheHook(hook types.CacheHook) {
	rc.hookrwm.Lock()
	defer rc.hookrwm.Unlock()
	rc.hooks = append(rc.hooks, hook)
	// run hook on existing items
	for _, item := range rc.List() {
		if hook.Predicate(types.CacheSet, item.K, item.V) {
			hook.Hook(item.K, item.V)
		}
	}
}

// Run start a goroutine to resync cache periodically with server
func (rc *ResourceCache) Run() {
	initialised := make(chan bool)
	go rc.AutoCacheBuilder(initialised)

	// Periodically flush cache to configmap for recovery after application restart or pod restart
	go rc.CacheToConfigMapDumper()

	// Wait for first initialization
	<-initialised
}

func (rc *ResourceCache) AutoCacheBuilder(initialised chan<- bool) {
	gauge := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "resource",
			Subsystem:   "cache",
			Name:        "rebuild_time",
			Help:        "Time taken to rebuild cache",
			ConstLabels: prometheus.Labels{"run": "auto"},
		})

	if !rc.stateLoaded {
		rc.rebuildCache(gauge)
		initialised <- true
	} else {
		initialised <- true
	}
	close(initialised)

	for {
		// to keep constant interval between cache rebuild runs, as rebuild cache is heavy operation so ticker may lead back to back runs
		time.Sleep(rc.resyncPeriod)
		rc.rebuildCache(gauge)
	}
}

func (rc *ResourceCache) CacheToConfigMapDumper() {
	func(resourceCache *ResourceCache) {
		for range time.Tick(resourceCache.flushTimeToCM) {
			debugID := util.GetShortUUID()
			lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"debug_id": debugID}))
			log := lmlog.Logger(lctx)
			start := time.Now()
			err := resourceCache.Save(lctx)
			if err != nil {
				log.Errorf("Flush cache to cm failed: %s", err)
			} else {
				log.Infof("Cache flushed to configmaps in time %v", time.Since(start))
			}
		}
	}(rc)
}

// Rebuild graceful rebuild cache
func (rc *ResourceCache) Rebuild(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	log.Infof("Gracefully building cache")
	gauge := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "cache",
			Subsystem:   "resource",
			Name:        "rebuild_time",
			Help:        "Time taken to rebuild cache",
			ConstLabels: prometheus.Labels{"run": "graceful"},
		})
	rc.rebuildCache(gauge)
}

// Set adds entry into cache map
func (rc *ResourceCache) Set(lctx *lmctx.LMContext, name types.ResourceName, meta types.ResourceMeta) bool {
	log := lmlog.Logger(lctx)
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	log.Tracef("Deleting previous entry if any...")
	rc.unsetInternal(lctx, name, meta.Container)
	log.Tracef("Setting cache entry %v: %v", name, meta)
	ok := rc.store.Set(lctx, name, meta)
	if ok {
		go func() {
			rc.hookrwm.RLock()
			defer rc.hookrwm.RUnlock()
			for _, hook := range rc.hooks {
				if hook.Predicate(types.CacheSet, name, meta) {
					hook.Hook(name, meta)
				}
			}
		}()
	}
	return ok
}

// Exists checks entry into cache map
func (rc *ResourceCache) Exists(lctx *lmctx.LMContext, name types.ResourceName, container string, softRefresh bool) (types.ResourceMeta, bool) {
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	log := lmlog.Logger(lctx)
	log.Tracef("Checking cache entry %v: %v", name, container)
	meta, ok := rc.store.Exists(lctx, name, container)
	if !ok && softRefresh {
		rc.SoftRefresh(lctx, container)
	}
	if softRefresh {
		return rc.store.Exists(lctx, name, container)
	}
	return meta, ok
}

// Get checks entry into cache map
func (rc *ResourceCache) Get(lctx *lmctx.LMContext, name types.ResourceName) ([]types.ResourceMeta, bool) {
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	log := lmlog.Logger(lctx)
	log.Tracef("Get cache entry list %v", name)

	return rc.store.Get(lctx, name)
}

// Unset checks entry into cache map
func (rc *ResourceCache) Unset(lctx *lmctx.LMContext, name types.ResourceName, container string) bool {
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	meta, ok := rc.unsetInternal(lctx, name, container)
	if ok {
		go func() {
			rc.hookrwm.RLock()
			defer rc.hookrwm.RUnlock()
			for _, hook := range rc.hooks {
				if hook.Predicate(types.CacheUnset, name, meta) {
					hook.Hook(name, meta)
				}
			}
		}()
	}
	return ok
}

func (rc *ResourceCache) unsetInternal(lctx *lmctx.LMContext, name types.ResourceName, container string) (types.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	/*// Special handling for resource groups, do not add another.
	// when parent resource group is deleted its all child hierarchy gets deleted from portal, hence
	if name.Resource == enums.Namespaces {
		if exists, ok := rc.store.Exists(lctx, name, container); ok {
			list := rc.store.ListWithFilter(func(k cache.ResourceName, v cache.ResourceMeta) bool {
				return v.Container == exists.Container ||
					(k.Resource.IsNamespaceScopedResource() && v.Container == name.GroupName)
			})
			log.Infof("Removing dependent (%d) container cache entries from cache of %s", len(list), container)
			for _, item := range list {
				rc.store.Unset(item.K, item.V.Container)
			}
		}
	}*/
	log.Tracef("Deleting cache entry %v: %v", name, container)

	return rc.store.Unset(lctx, name, container)
}

// Load loads cache from configmaps
func (rc *ResourceCache) Load(lctx *lmctx.LMContext) error {
	cmList, err := rc.listAllCacheConfigMaps()
	if err != nil {
		return err
	}
	if cmList.Size() == 0 {
		return fmt.Errorf("no config maps found")
	}
	selectedDumpID := rc.selectDumpID(lctx, cmList)
	tmpCache := NewStore()
	err2 := rc.populateCacheStore(lctx, cmList, selectedDumpID, tmpCache)
	if err2 != nil {
		return err2
	}
	if tmpCache.Size() == 0 {
		return fmt.Errorf("no cache data found in present configmaps")
	}
	rc.resetCacheStore(tmpCache)
	rc.dumpID = selectedDumpID
	rc.updateIncrementalCache(lctx)

	return nil
}

// List returns slices of all data present in cache - For ex: cache for periodic dangling resource etc uses this
func (rc *ResourceCache) List() []types.IterItem {
	return rc.store.List()
}

// ListWithFilter returns slices of all data present in cache - For ex: cache for periodic dangling resource etc uses this
func (rc *ResourceCache) ListWithFilter(f func(k types.ResourceName, v types.ResourceMeta) bool) []types.IterItem {
	return rc.store.ListWithFilter(f)
}

// UnsetLMID unsets value using id only
// WARN: This is not an o(1) map operation hence do not use until necessary, this is heavy operation to find cache entry with id
func (rc *ResourceCache) UnsetLMID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) bool {
	k, m := rc.store.getInternalUsingLMID(rt, id)
	return rc.Unset(lctx, k, m.Container)
}

func (rc *ResourceCache) SoftRefresh(lctx *lmctx.LMContext, container string) {
	conf, err := config.GetConfig()
	if err != nil {
		return
	}
	if t, ok := rc.getSoftRefreshLastTime(container); ok && t.Before(time.Now().Add(-time.Minute)) {
		list, ok := rc.Get(lctx, types.ResourceName{Name: container, Resource: enums.Namespaces})
		if ok {
			for _, meta := range list {
				resp, err := rc.getDevices(lctx, meta.LMID)
				if err != nil && resp != nil {
					for _, resource := range resp.Payload.Items {
						rc.storeDevice(lctx, resource, conf.ClusterName, rc.store)
					}
				}
			}
		}
	}
}

func (rc *ResourceCache) getSoftRefreshLastTime(key string) (time.Time, bool) {
	rc.softRefreshMu.Lock()
	defer rc.softRefreshMu.Unlock()
	if val, ok := rc.softRefreshLast[key]; ok {
		return val, ok
	}

	return time.Time{}, false
}
