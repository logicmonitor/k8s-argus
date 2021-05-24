package devicecache_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
)

func BenchmarkCacheSetSameObject(b *testing.B) {
	cache2 := devicecache.NewResourceCache(nil, 1*time.Second)
	for i := 0; i < b.N; i++ {
		cache2.Set(cache.ResourceName{Name: "abc", Resource: 1}, cache.ResourceMeta{ // nolint: exhaustivestruct
			Container:     "xyz",
			LMID:          0,
			DisplayName:   "abc",
			SysCategories: []string{"KubernetesPod"},
		})
	}
}

func BenchmarkCacheSetDiffObject(b *testing.B) {
	cache2 := devicecache.NewResourceCache(nil, 1*time.Second)
	for i := 0; i < b.N; i++ {
		cache2.Set(cache.ResourceName{Name: "abc", Resource: 1}, cache.ResourceMeta{ // nolint: exhaustivestruct
			Container:     fmt.Sprintf("xyz-%v", i),
			LMID:          0,
			DisplayName:   "abc",
			SysCategories: []string{"KubernetesPod"},
		})
	}
}

func BenchmarkCacheSetDiffDuplObject(b *testing.B) {
	cache2 := devicecache.NewResourceCache(nil, 1*time.Second)
	for i := 0; i < b.N; i++ {
		cache2.Set(cache.ResourceName{Name: "abc", Resource: 1}, cache.ResourceMeta{ // nolint: exhaustivestruct
			Container:     fmt.Sprintf("xyz-%v", i%5),
			LMID:          0,
			DisplayName:   "abc",
			SysCategories: []string{"KubernetesPod"},
		})
	}
}

func TestNewResourceCache(t *testing.T) {
	t.Parallel()
	m := make(map[cache.ResourceName][]cache.ResourceMeta)

	for i := 0; i < 10; i++ {
		key := cache.ResourceName{
			Name:     fmt.Sprintf("abc-%v", i),
			Resource: 1,
		}
		value := cache.ResourceMeta{ // nolint: exhaustivestruct
			Container:     "ns",
			LMID:          int32(i * 10),
			DisplayName:   "abc-ns",
			SysCategories: []string{"KubernetesPod"},
		}
		list, ok := m[key]
		if !ok {
			list = []cache.ResourceMeta{} // nolint: exhaustivestruct
		}
		list = append(list, value)
		m[key] = list
	}

	arr, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Failed to marshal map %v", err)
	}
	t.Logf("%s", arr)

	m2 := make(map[cache.ResourceName][]cache.ResourceMeta)
	if er := json.Unmarshal(arr, &m2); er != nil {
		t.Errorf("Failed to unmarshal %v", er)
	}
	t.Logf("Before %v", m)
	t.Logf("Unmarshaled map %v", m2)
}
