package builder

import (
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/models"
	log "github.com/sirupsen/logrus"
)

// Builder implements types.DeviceBuilder
type Builder struct {
	types.DeviceBuilder
}

// Name implements types.DeviceBuilder.
func (b *Builder) Name(name string) types.DeviceOption {
	return func(device *models.Device) {
		device.Name = &name
	}
}

// DisplayName implements types.DeviceBuilder.
func (b *Builder) DisplayName(name string) types.DeviceOption {
	return func(device *models.Device) {
		device.DisplayName = &name
	}
}

// CollectorID implements types.DeviceBuilder.
func (b *Builder) CollectorID(id int32) types.DeviceOption {
	return func(device *models.Device) {
		device.PreferredCollectorID = &id
	}
}

// SystemCategories implements types.DeviceBuilder
func (b *Builder) SystemCategories(categories string) types.DeviceOption {
	return setProperty(constants.K8sSystemCategoriesPropertyKey, categories)
}

// ResourceLabels implements types.DeviceBuilder
func (b *Builder) ResourceLabels(properties map[string]string) types.DeviceOption {
	return func(device *models.Device) {
		if device == nil {
			return
		}
		if device.CustomProperties == nil {
			device.CustomProperties = []*models.NameAndValue{}
		}
		for name, value := range properties {
			propName := constants.LabelCustomPropertyPrefix + name
			propValue := value
			if propValue == "" {
				propValue = constants.LabelNullPlaceholder
			}
			existed := false
			for _, prop := range device.CustomProperties {
				if *prop.Name == propName {
					*prop.Value = propValue
					existed = true
					break
				}
			}
			if !existed {
				device.CustomProperties = append(device.CustomProperties, &models.NameAndValue{
					Name:  &propName,
					Value: &propValue,
				})
			}
		}
	}
}

// Auto implements types.DeviceBuilder
func (b *Builder) Auto(name, value string) types.DeviceOption {
	return setProperty("auto."+name, value)
}

// System implements types.DeviceBuilder.
func (b *Builder) System(name, value string) types.DeviceOption {
	return setProperty("system."+name, value)
}

// Custom implements types.DeviceBuilder.
func (b *Builder) Custom(name, value string) types.DeviceOption {
	return setProperty(name, value)
}

func setProperty(name, value string) types.DeviceOption {
	return func(device *models.Device) {
		if device == nil {
			return
		}
		if device.CustomProperties == nil {
			device.CustomProperties = []*models.NameAndValue{}
		}
		for _, prop := range device.CustomProperties {
			if *prop.Name == name && value != "" {
				if *prop.Name == constants.K8sSystemCategoriesPropertyKey {
					value = getUpdatedSystemCategories(*prop.Value, value)
				}
				*prop.Value = value
				return
			}
		}
		if value != "" {
			device.CustomProperties = append(device.CustomProperties, &models.NameAndValue{
				Name:  &name,
				Value: &value,
			})
		} else {
			log.Warnf("Custom property value is empty for %q, skipping", name)
		}
	}
}

func getUpdatedSystemCategories(oldValue, newValue string) string {
	// we do not use strings.contain, because it may be matched as substring of some prop
	oldValues := strings.Split(strings.TrimSpace(oldValue), ",")
	for _, ov := range oldValues {
		if ov == newValue {
			return oldValue
		}
	}
	oldValue = oldValue + "," + newValue
	return oldValue
}
