package types

import (
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/models"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// Base is a struct for embedding
type Base struct {
	LMClient  *client.LMSdkGo
	K8sClient *kubernetes.Clientset
	Config    *config.Config
}

// Watcher is the LogicMonitor Watcher interface.
type Watcher interface {
	APIVersion() string
	CheckRBAC() bool
	Resource() string
	ObjType() runtime.Object
	AddFunc() func(obj interface{})
	DeleteFunc() func(obj interface{})
	UpdateFunc() func(oldObj, newObj interface{})
}

// DeviceManager is an interface that describes how resources in Kubernetes
// are mapped into LogicMonitor as devices.
type DeviceManager interface {
	DeviceMapper
	DeviceBuilder
}

// DeviceMapper is the interface responsible for mapping a Kubernetes resource to
// a LogicMonitor device.
type DeviceMapper interface {
	// Config returns the Argus config.
	Config() *config.Config
	// FindByDisplayName searches for a device by it's display name. It will return a device if and only if
	// one device was found, and return nil otherwise.
	FindByDisplayName(string) (*models.Device, error)
	// Add adds a device to a LogicMonitor account.
	Add(...DeviceOption) (*models.Device, error)
	// UpdateAndReplaceByID updates a device using the 'replace' OpType.
	UpdateAndReplaceByID(int32, ...DeviceOption) (*models.Device, error)
	// UpdateAndReplaceByDisplayName updates a device using the 'replace' OpType if and onlt if it does not already exist.
	UpdateAndReplaceByDisplayName(string, ...DeviceOption) (*models.Device, error)
	// UpdateAndReplaceFieldByID updates a device using the 'replace' OpType for a
	// specific field of a device.
	UpdateAndReplaceFieldByID(int32, string, ...DeviceOption) (*models.Device, error)
	// UpdateAndReplaceFieldByDisplayName updates a device using the 'replace' OpType for a
	// specific field of a device.
	UpdateAndReplaceFieldByDisplayName(string, string, ...DeviceOption) (*models.Device, error)
	// DeleteByID deletes a device by device ID.
	DeleteByID(int32) error
	// DeleteByDisplayName deletes a device by device display name.
	DeleteByDisplayName(string) error
}

// DeviceOption is the function definition for the functional options pattern.
type DeviceOption func(*models.Device)

// DeviceBuilder is the interface responsible for building a device struct.
type DeviceBuilder interface {
	// Name sets the device name.
	Name(string) DeviceOption
	// DisplayName sets the device name.
	DisplayName(string) DeviceOption
	// CollectorID sets the preferred collector ID for the device.
	CollectorID(int32) DeviceOption
	// SystemCategories sets the system.categories property on the device.
	SystemCategories(string) DeviceOption
	// ResourceLabels sets custom properties for the device
	ResourceLabels(map[string]string) DeviceOption
	// Auto adds an auto property to the device.
	Auto(string, string) DeviceOption
	// System adds a system property to the device.
	System(string, string) DeviceOption
	// System adds a custom property to the device.
	Custom(string, string) DeviceOption
}
