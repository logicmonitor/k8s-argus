package resourcecache

import (
	"sync"

	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
)

// Store store
type Store struct {
	InternalMap map[types.ResourceName][]types.ResourceMeta `json:"internal_map"`
	rwm         sync.RWMutex
}

// NewStore new
func NewStore() *Store {
	return &Store{
		InternalMap: make(map[types.ResourceName][]types.ResourceMeta),
		rwm:         sync.RWMutex{},
	}
}

// Set adds entry into cache map
func (store *Store) Set(name types.ResourceName, meta types.ResourceMeta) bool {
	logrus.Tracef("Setting cache entry %s", name)
	store.rwm.Lock()
	defer store.rwm.Unlock()
	b, done := store.setInternal(name, meta)
	if done {
		return b
	}

	return true
}

func (store *Store) setInternal(name types.ResourceName, meta types.ResourceMeta) (bool, bool) {
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
func (store *Store) AddAll(another map[types.ResourceName][]types.ResourceMeta) bool {
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
func (store *Store) Exists(lctx *lmctx.LMContext, name types.ResourceName, namespace string) (types.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	log.Tracef("Checking cache entry %s", name)
	store.rwm.RLock()
	defer store.rwm.RUnlock()
	list, ok := store.InternalMap[name]
	if !ok {
		return types.ResourceMeta{}, false // nolint: exhaustivestruct
	}
	for _, meta := range list {
		if meta.Container == namespace {
			return meta, true
		}
	}

	return types.ResourceMeta{}, false // nolint: exhaustivestruct
}

// Get get list of all resource meta for mentioned resource name
func (store *Store) Get(lctx *lmctx.LMContext, name types.ResourceName) ([]types.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	log.Tracef("Get cache entry list %v", name)
	store.rwm.RLock()
	defer store.rwm.RUnlock()
	list, ok := store.InternalMap[name]
	log.Tracef("Cache entry list: %v: %v", list, ok)
	if !ok || len(list) == 0 {
		return []types.ResourceMeta{}, false // nolint: exhaustivestruct
	}

	return list, true
}

// Unset checks entry into cache map
func (store *Store) Unset(name types.ResourceName, namespace string) (types.ResourceMeta, bool) {
	logrus.Tracef("Deleting cache entry %s", name)
	store.rwm.Lock()
	defer store.rwm.Unlock()
	list, ok := store.InternalMap[name]
	if !ok {
		return types.ResourceMeta{}, false
	}
	var meta types.ResourceMeta
	for idx, v := range list {
		if v.Container == namespace {
			meta = v
			list[idx] = list[len(list)-1]
			list = list[:len(list)-1]

			break
		}
	}
	store.InternalMap[name] = list

	return meta, true
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
func (store *Store) List() []types.IterItem {
	store.rwm.Lock()
	defer store.rwm.Unlock()
	var list []types.IterItem
	for k, v := range store.InternalMap {
		for _, m := range v {
			list = append(list, types.IterItem{K: k, V: m})
		}
	}

	return list
}

// ListWithFilter list
func (store *Store) ListWithFilter(f func(k types.ResourceName, v types.ResourceMeta) bool) []types.IterItem {
	store.rwm.Lock()
	defer store.rwm.Unlock()
	var list []types.IterItem
	for k, v := range store.InternalMap {
		for _, m := range v {
			if f(k, m) {
				list = append(list, types.IterItem{K: k, V: m})
			}
		}
	}

	return list
}

func (store *Store) getInternalUsingLMID(rt enums.ResourceType, id int32) (types.ResourceName, types.ResourceMeta) {
	store.rwm.RLock()
	defer store.rwm.RUnlock()

	var k types.ResourceName
	var v []types.ResourceMeta
	var m types.ResourceMeta
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

	return k, m
}

func (store *Store) getMap() map[types.ResourceName][]types.ResourceMeta {
	store.rwm.RLock()
	defer store.rwm.RUnlock()
	m := make(map[types.ResourceName][]types.ResourceMeta)
	for k, v := range store.InternalMap {
		m[k] = v
	}
	return m
}
