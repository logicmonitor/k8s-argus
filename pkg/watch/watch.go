package watch

import "k8s.io/client-go/pkg/runtime"

// Watcher is the LogicMonitor Watcher interface.
type Watcher interface {
	Resource() string
	ObjType() runtime.Object
	AddFunc() func(obj interface{})
	DeleteFunc() func(obj interface{})
	UpdateFunc() func(oldObj, newObj interface{})
}
