package resourcewatcher

import (
	"github.com/logicmonitor/k8s-argus/pkg/enums"
)

// Watcher represents a watcher type that watches pods.
type Watcher struct {
	Resource enums.ResourceType
}

// ResourceType resource
func (w *Watcher) ResourceType() enums.ResourceType {
	return w.Resource
}
