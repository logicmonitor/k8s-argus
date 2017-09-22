package builder

import (
	"github.com/logicmonitor/k8s-argus/pkg/types"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
)

// Builder implements types.DeviceBuilder
type Builder struct {
	types.DeviceBuilder
}

// Name implements types.DeviceBuilder.
func (b *Builder) Name(name string) types.DeviceOption {
	return func(device *lm.RestDevice) {
		device.Name = name
	}
}

// DisplayName implements types.DeviceBuilder.
func (b *Builder) DisplayName(name string) types.DeviceOption {
	return func(device *lm.RestDevice) {
		device.DisplayName = name
	}
}

// CollectorID implements types.DeviceBuilder.
func (b *Builder) CollectorID(id int32) types.DeviceOption {
	return func(device *lm.RestDevice) {
		device.PreferredCollectorId = id
	}
}

// SystemCategories implements types.DeviceBuilder.
func (b *Builder) SystemCategories(categories string) types.DeviceOption {
	return setProperty("system.categories", categories)
}

// Auto implements types.DeviceBuilder.
func (b *Builder) Auto(name, value string) types.DeviceOption {
	return setProperty("auto."+name, value)
}

// System implements types.DeviceBuilder.
func (b *Builder) System(name, value string) types.DeviceOption {
	return setProperty("system."+name, value)
}

func setProperty(name, value string) types.DeviceOption {
	return func(device *lm.RestDevice) {
		if value != "" {
			device.CustomProperties = append(device.CustomProperties, lm.NameAndValue{
				Name:  name,
				Value: value,
			})
		} else {
			log.Warnf("Custom property value is empty for %q, skipping", name)
		}
	}
}
