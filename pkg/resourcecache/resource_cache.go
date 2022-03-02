package resourcecache

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
)

var (
	// Version to differentiate cache flushed to configmap is with previous datastructures or not.
	// Whenever cache.ResourceName and cache.ResourceMeta struct modifies, increase the version so that it will easier to parse previously created CMs and convert it to newer data structure
	Version = "v2"

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
	flushTimeToCM time.Duration
	flushMU       sync.Mutex
	*types.LMRequester
	clusterGrpID int32
	stateLoaded  bool

	// On restart, first dump loaded from config map
	dumpID int64

	// soft refresh last time
	softRefreshLast map[string]time.Time
	softRefreshMu   sync.RWMutex

	hooks   []types.CacheHook
	hookrwm sync.RWMutex
}

// NewResourceCache create new ResourceCache object
func NewResourceCache(facadeObj *types.LMRequester) *ResourceCache {
	resourceCache := &ResourceCache{
		store:           NewStore(),
		rwm:             sync.RWMutex{},
		flushTimeToCM:   1 * time.Minute,
		flushMU:         sync.Mutex{},
		LMRequester:     facadeObj,
		clusterGrpID:    -1,
		stateLoaded:     false,
		rebuildMutex:    sync.Mutex{},
		softRefreshMu:   sync.RWMutex{},
		softRefreshLast: make(map[string]time.Time),
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
	// register after first initialization to avoid half cache to override previously stored cache
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"shutdown_hook": "resource_cache"}))
		log := lmlog.Logger(lctx)
		log.Infof("Shutdown hook registered to storing cache to configmap on exit")
		<-c
		log.Infof("Shutting down, storing cache to configmap")
		err := rc.Save(lctx)
		if err != nil {
			log.Errorf("Failed to store cache to configmap")
		}
		log.Infof("Stored cache to configmap")
		os.Exit(1)
	}()
}

func (rc *ResourceCache) AutoCacheBuilder(initialised chan<- bool) {
	if !rc.stateLoaded {
		rc.rebuildCache()
		initialised <- true
	} else {
		initialised <- true
	}
	close(initialised)
	lctx := lmlog.NewLMContextWith(nil)
	for {
		sleep := time.Hour
		conf, err := config.GetConfig(lctx)
		if err == nil {
			sleep = *conf.Intervals.CacheSyncInterval
		}
		// to keep constant interval between cache rebuild runs, as rebuild cache is heavy operation so ticker may lead back to back runs
		time.Sleep(sleep)
		rc.rebuildCache()
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
	rc.rebuildCache()
}

// Set adds entry into cache map
func (rc *ResourceCache) Set(lctx *lmctx.LMContext, name types.ResourceName, meta types.ResourceMeta) bool {
	log := lmlog.Logger(lctx)
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	if name.Resource != enums.Namespaces {
		log.Tracef("Deleting previous entry if any...")
		rc.unsetInternal(lctx, name, meta.Container)
	}
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
	if meta.IsInvalid && softRefresh {
		log.Debugf("Invalid cache entry hence refreshing its container...")
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
	log := lmlog.Logger(lctx)
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	rc.unsetInternal(lctx, name, container)
	log.Tracef("Deleting cache entry %v: %v", name, container)
	meta, ok := rc.store.Unset(lctx, name, container)
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

func (rc *ResourceCache) unsetInternal(lctx *lmctx.LMContext, name types.ResourceName, container string) bool {
	log := lmlog.Logger(lctx)
	// Special handling for resource groups, do not add another.
	// when parent resource group is deleted its all child hierarchy gets deleted from portal, hence
	if name.Resource == enums.Namespaces {
		/*if exists, ok := rc.store.Exists(lctx, name, container); ok {
			list := rc.store.ListWithFilter(func(k types.ResourceName, v types.ResourceMeta) bool {
				return v.Container == exists.Container ||
					(k.Resource.IsNamespaceScopedResource() && v.Container == name.Name)
			})
			log.Infof("Removing dependent (%d) container cache entries from cache of %s", len(list), container)
			for _, item := range list {
				rc.store.Unset(item.K, item.V.Container)
			}
		}*/

		// delete child groups cache entries if parent got deleted
		if exists, ok := rc.store.Exists(lctx, name, container); ok {
			list := rc.store.ListWithFilter(func(k types.ResourceName, v types.ResourceMeta) bool {
				return k.Resource == enums.Namespaces && v.Container == fmt.Sprintf("%v", exists.LMID)
			})
			if len(list) > 0 { // check to avoid excessive logs at info level
				log.Infof("Removing (%d) child group cache entries having container (%v) ", len(list), exists.LMID)
				for _, item := range list {
					rc.store.Unset(lctx, item.K, item.V.Container)
				}
			}
		}
	}

	return true
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
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Tracef("soft refresh: get config failed: %s", err)
		return
	}
	log.Tracef("soft refresh: starting...")
	SoftRefreshMinDuration := time.Minute
	if t, ok := rc.getSoftRefreshLastTime(container); ok && t.Before(time.Now().Add(-SoftRefreshMinDuration)) {
		log.Tracef("soft refresh: valid for soft refresh")
		list, ok := rc.Get(lctx, types.ResourceName{Name: container, Resource: enums.Namespaces})
		log.Tracef("soft refresh: container details: %v: %v", list, ok)
		if ok {
			for _, meta := range list {
				resp, err := rc.getDevices(lctx, meta.LMID)
				log.Tracef("soft refresh: resources fetched of %v: %v", meta.LMID, resp)
				if err == nil && resp != nil {
					log.Tracef("soft refresh: storing resources in cache: %v", resp)
					for _, resource := range resp {
						log.Tracef("soft refresh: storing resource in cache: %v", resource)
						rc.storeDevice(lctx, resource, conf.ClusterName, rc.store)
					}
				}
			}
		} else {
			log.Tracef("soft refresh: no cache entry found for container")
		}
	} else {
		log.Tracef("soft refresh: ignoring soft refresh less than %v", SoftRefreshMinDuration)
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

func (rc *ResourceCache) setLastUpdateTime(key types.ResourceName) {
	rc.softRefreshMu.Lock()
	defer rc.softRefreshMu.Unlock()
	t := time.Now()
	rc.softRefreshLast[key.Name] = t
}
