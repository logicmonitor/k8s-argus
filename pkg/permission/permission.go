package permission

import (
	"sync"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

var (
	m  sync.Map
	mu sync.Mutex
)

// HasPermissions has permission
func HasPermissions(rt enums.ResourceType) bool {
	if rt == enums.ETCD || rt == enums.Unknown {
		return true
	}
	load, ok := m.Load(rt)
	if ok {
		return load.(bool)
	}
	mu.Lock()
	defer mu.Unlock()
	listWatch := cache.NewListWatchFromClient(util.GetK8sRESTClient(config.GetClientSet(), rt.K8SAPIVersion()), rt.String(), corev1.NamespaceAll, fields.Everything())
	listWatch.DisableChunking = true
	_, err := listWatch.List(constants.DefaultListOptions)
	m.Store(rt, err == nil)
	return err == nil
}
