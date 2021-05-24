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
func (dc *Store) Set(name cache.ResourceName, meta cache.ResourceMeta) bool {
	logrus.Tracef("Setting cache entry %s", name)
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	b, done := dc.setInternal(name, meta)
	if done {
		return b
	}

	return true
}

func (dc *Store) setInternal(name cache.ResourceName, meta cache.ResourceMeta) (bool, bool) {
	list := dc.InternalMap[name]
	for idx, m := range list {
		if m.Container == meta.Container {
			list[idx] = meta

			return true, true
		}
	}
	list = append(list, meta)
	dc.InternalMap[name] = list

	return true, true
}

// AddAll adds entry into cache map
func (dc *Store) AddAll(another map[cache.ResourceName][]cache.ResourceMeta) bool {
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	for k, v := range another {
		for _, meta := range v {
			dc.setInternal(k, meta)
		}
	}

	return true
}

// Exists checks entry into cache map
func (dc *Store) Exists(lctx *lmctx.LMContext, name cache.ResourceName, namespace string) (cache.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	log.Tracef("Checking cache entry %s", name)
	dc.rwm.RLock()
	defer dc.rwm.RUnlock()
	list, ok := dc.InternalMap[name]
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
func (dc *Store) Get(lctx *lmctx.LMContext, name cache.ResourceName) ([]cache.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	log.Tracef("Get cache entry list %v", name)
	dc.rwm.RLock()
	defer dc.rwm.RUnlock()
	list, ok := dc.InternalMap[name]
	log.Tracef("Cache entry list: %v: %v", list, ok)
	if !ok {
		return []cache.ResourceMeta{}, false // nolint: exhaustivestruct
	}

	return list, true
}

// Unset checks entry into cache map
func (dc *Store) Unset(name cache.ResourceName, namespace string) bool {
	logrus.Tracef("Deleting cache entry %s", name)
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	list, ok := dc.InternalMap[name]
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
	dc.InternalMap[name] = list

	return true
}

// Size size
func (dc *Store) Size() int64 {
	size := int64(0)
	for _, v := range dc.InternalMap {
		size += int64(len(v))
	}

	return size
}

// List list
func (dc *Store) List() []cache.IterItem {
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	var list []cache.IterItem
	for k, v := range dc.InternalMap {
		for _, m := range v {
			list = append(list, cache.IterItem{K: k, V: m})
		}
	}

	return list
}

// UnsetLMID unset
func (dc *Store) UnsetLMID(rt enums.ResourceType, id int32) bool {
	dc.rwm.RLock()
	var k cache.ResourceName
	var v []cache.ResourceMeta
	var m cache.ResourceMeta
	for k, v = range dc.InternalMap {
		if k.Resource == rt {
			for _, meta := range v {
				if meta.LMID == id {
					m = meta

					break
				}
			}
		}
	}
	dc.rwm.RUnlock()

	return dc.Unset(k, m.Container)
}
