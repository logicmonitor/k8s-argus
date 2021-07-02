package resourcecache_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resourcecache"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
)

var lctx = lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"resource_cache": "test"}))

func BenchmarkCacheSetSameObject(b *testing.B) {
	cache2 := resourcecache.NewResourceCache(nil, 1*time.Second)
	for i := 0; i < b.N; i++ {
		cache2.Set(lctx, types.ResourceName{Name: "abc", Resource: 1}, types.ResourceMeta{ // nolint: exhaustivestruct
			Container:     "xyz",
			LMID:          0,
			DisplayName:   "abc",
			SysCategories: []string{"KubernetesPod"},
		})
	}
}

func BenchmarkCacheSetDiffObject(b *testing.B) {
	cache2 := resourcecache.NewResourceCache(nil, 1*time.Second)
	for i := 0; i < b.N; i++ {
		cache2.Set(lctx, types.ResourceName{Name: "abc", Resource: 1}, types.ResourceMeta{ // nolint: exhaustivestruct
			Container:     fmt.Sprintf("xyz-%v", i),
			LMID:          0,
			DisplayName:   "abc",
			SysCategories: []string{"KubernetesPod"},
		})
	}
}

func BenchmarkCacheSetDiffDuplObject(b *testing.B) {
	cache2 := resourcecache.NewResourceCache(nil, 1*time.Second)
	for i := 0; i < b.N; i++ {
		cache2.Set(lctx, types.ResourceName{Name: "abc", Resource: 1}, types.ResourceMeta{ // nolint: exhaustivestruct
			Container:     fmt.Sprintf("xyz-%v", i%5),
			LMID:          0,
			DisplayName:   "abc",
			SysCategories: []string{"KubernetesPod"},
		})
	}
}

func TestNewResourceCache(t *testing.T) {
	t.Parallel()
	m := make(map[types.ResourceName][]types.ResourceMeta)

	for i := 0; i < 10; i++ {
		key := types.ResourceName{
			Name:     fmt.Sprintf("abc-%v", i),
			Resource: 1,
		}
		value := types.ResourceMeta{ // nolint: exhaustivestruct
			Container:     "ns",
			LMID:          int32(i * 10),
			DisplayName:   "abc-ns",
			SysCategories: []string{"KubernetesPod"},
		}
		list, ok := m[key]
		if !ok {
			list = []types.ResourceMeta{} // nolint: exhaustivestruct
		}
		list = append(list, value)
		m[key] = list
	}

	arr, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Failed to marshal map %v", err)
	}
	t.Logf("%s", arr)

	m2 := make(map[types.ResourceName][]types.ResourceMeta)
	if er := json.Unmarshal(arr, &m2); er != nil {
		t.Errorf("Failed to unmarshal %v", er)
	}
	t.Logf("Before %v", m)
	t.Logf("Unmarshaled map %v", m2)
}
