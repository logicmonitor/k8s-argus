package devicecache

import (
	"sync"

	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/sirupsen/logrus"
)

// Store store
type Store struct {
	InternalMap map[cache.ResourceName][]cache.ResourceMeta `json:"internal_map"`
	rwm         sync.RWMutex
}

// NewStore new
func NewStore() *Store {
	return &Store{
		InternalMap: make(map[cache.ResourceName][]cache.ResourceMeta),
		rwm:         sync.RWMutex{},
	}
}

// Set adds entry into cache map
func (store *Store) Set(name cache.ResourceName, meta cache.ResourceMeta) bool {
	logrus.Tracef("Setting cache entry %s", name)
	store.rwm.Lock()
	defer store.rwm.Unlock()
	b, done := store.setInternal(name, meta)
	if done {
		return b
	}

	return true
}

func (store *Store) setInternal(name cache.ResourceName, meta cache.ResourceMeta) (bool, bool) {
	list := store.InternalMap[name]
	for idx, m := range list {
		if m.Container == meta.Container {
			list[idx] = meta

			return true, true
		}
	}
	list = append(list, meta)
	store.InternalMap[name] = list

	return true, true
}

// AddAll adds entry into cache map
func (store *Store) AddAll(another map[cache.ResourceName][]cache.ResourceMeta) bool {
	store.rwm.Lock()
	defer store.rwm.Unlock()
	for k, v := range another {
		for _, meta := range v {
			store.setInternal(k, meta)
		}
	}

	return true
}

// Exists checks entry into cache map
func (store *Store) Exists(lctx *lmctx.LMContext, name cache.ResourceName, namespace string) (cache.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	log.Tracef("Checking cache entry %s", name)
	store.rwm.RLock()
	defer store.rwm.RUnlock()
	list, ok := store.InternalMap[name]
	if !ok {
		return cache.ResourceMeta{}, false // nolint: exhaustivestruct
	}
	for _, meta := range list {
		if meta.Container == namespace {
			return meta, true
		}
	}

	return cache.ResourceMeta{}, false // nolint: exhaustivestruct
}

// Get get list of all resource meta for mentioned resource name
func (store *Store) Get(lctx *lmctx.LMContext, name cache.ResourceName) ([]cache.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	log.Tracef("Get cache entry list %v", name)
	store.rwm.RLock()
	defer store.rwm.RUnlock()
	list, ok := store.InternalMap[name]
	log.Tracef("Cache entry list: %v: %v", list, ok)
	if !ok {
		return []cache.ResourceMeta{}, false // nolint: exhaustivestruct
	}

	return list, true
}

// Unset checks entry into cache map
func (store *Store) Unset(name cache.ResourceName, namespace string) bool {
	logrus.Tracef("Deleting cache entry %s", name)
	store.rwm.Lock()
	defer store.rwm.Unlock()
	list, ok := store.InternalMap[name]
	if !ok {
		return true
	}
	for idx, v := range list {
		if v.Container == namespace {
			list[idx] = list[len(list)-1]
			list = list[:len(list)-1]

			break
		}
	}
	store.InternalMap[name] = list

	return true
}

// Size size
func (store *Store) Size() int64 {
	size := int64(0)
	for _, v := range store.InternalMap {
		size += int64(len(v))
	}

	return size
}

// List list
func (store *Store) List() []cache.IterItem {
	store.rwm.Lock()
	defer store.rwm.Unlock()
	var list []cache.IterItem
	for k, v := range store.InternalMap {
		for _, m := range v {
			list = append(list, cache.IterItem{K: k, V: m})
		}
	}

	return list
}

// UnsetLMID unset
func (store *Store) UnsetLMID(rt enums.ResourceType, id int32) bool {
	store.rwm.RLock()
	var k cache.ResourceName
	var v []cache.ResourceMeta
	var m cache.ResourceMeta
	for k, v = range store.InternalMap {
		if k.Resource == rt {
			for _, meta := range v {
				if meta.LMID == id {
					m = meta

					break
				}
			}
		}
	}
	store.rwm.RUnlock()

	return store.Unset(k, m.Container)
}
