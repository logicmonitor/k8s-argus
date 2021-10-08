package types

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/enums"
	k8stypes "k8s.io/apimachinery/pkg/types"
)

var separator = "##"

// ResourceName is key used in store map, so json specification needs to have string keys only, so custom Marshal & Unmarshal needs to be implemented to convert struct to string
type ResourceName struct {
	Name     string
	Resource enums.ResourceType
}

// MarshalText implements marshaler interface
func (i ResourceName) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%s%s%s", i.Name, separator, i.Resource.String())), nil
}

// UnmarshalText implements unmarshaler
func (i *ResourceName) UnmarshalText(text []byte) error {
	in := fmt.Sprintf("%s", text)
	arr := strings.Split(in, separator)
	i.Name = arr[0]
	rt, err := enums.ParseResourceType(strings.ToLower(arr[1]))
	if err != nil {
		return err
	}
	i.Resource = rt

	return nil
}

// ResourceMeta meta
type ResourceMeta struct {
	Container     string            `json:"container"`
	LMID          int32             `json:"lmid"`
	DisplayName   string            `json:"display_name"`
	Name          string            `json:"name"`
	Labels        map[string]string `json:"labels"`
	SysCategories []string          `json:"sys_categories"`
	UID           k8stypes.UID      `json:"uid"`
	CreatedOn     int64             `json:"created_on"`
}

func (resourceMeta ResourceMeta) HasSysCategory(category string) bool {
	for _, cat := range resourceMeta.SysCategories {
		if cat == category {
			return true
		}
	}
	return false
}

// IterItem intermediate structure to hold map entry
type IterItem struct {
	K ResourceName
	V ResourceMeta
}

type cacheHook func(rn ResourceName, meta ResourceMeta)

type cacheHookPredicate func(action CacheAction, rn ResourceName, meta ResourceMeta) bool

type CacheHook struct {
	Hook      cacheHook
	Predicate cacheHookPredicate
}

type CacheAction uint

const (
	// CacheSet new item
	CacheSet = iota

	// CacheUnset delete item
	CacheUnset
)

func (action CacheAction) String() string {
	switch action {
	case CacheSet:
		return "Set"
	case CacheUnset:
		return "Unset"
	}
	return "Unknown"
}
