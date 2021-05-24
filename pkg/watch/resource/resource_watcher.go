package resource

import (
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/types"
)

// Watcher represents a watcher type that watches pods.
type Watcher struct {
	Resource enums.ResourceType
	*types.WConfig
}

// ResourceType resource
func (w *Watcher) ResourceType() enums.ResourceType {
	return w.Resource
}
