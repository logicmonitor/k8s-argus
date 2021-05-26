package builder

import (
	"strconv"
	"strings"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Builder implements types.DeviceBuilder
type Builder struct{}

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

// SystemCategory implements types.DeviceBuilder
// nolint: predeclared
func (b *Builder) SystemCategory(category string, action enums.BuilderAction) types.DeviceOption {
	return setProperty(constants.K8sSystemCategoriesPropertyKey, category, action)
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
	return setProperty("auto."+name, value, enums.Add)
}

// System implements types.DeviceBuilder.
func (b *Builder) System(name, value string) types.DeviceOption {
	return setProperty("system."+name, value, enums.Add)
}

// Custom implements types.DeviceBuilder.
func (b *Builder) Custom(name, value string) types.DeviceOption {
	return setProperty(name, value, enums.Add)
}

// DeletedOn implements types.DeviceBuilder
func (b *Builder) DeletedOn(value time.Time) types.DeviceOption {
	return func(device *models.Device) {
		name := constants.K8sResourceDeletedOnPropertyKey

		if device == nil {
			return
		}

		if device.CustomProperties == nil {
			device.CustomProperties = []*models.NameAndValue{}
		}

		for _, prop := range device.CustomProperties {
			if *prop.Name == name {
				return
			}
		}

		strValue := strconv.FormatInt(value.Unix(), 10)
		device.CustomProperties = append(device.CustomProperties, &models.NameAndValue{
			Name:  &name,
			Value: &strValue,
		})
	}
}

// nolint: predeclared,gocognit,cyclop
func setProperty(name, value string, action enums.BuilderAction) types.DeviceOption {
	return func(device *models.Device) {
		if device == nil {
			return
		}
		if device.CustomProperties == nil {
			device.CustomProperties = []*models.NameAndValue{}
		}
		if action.Is(enums.Delete) {
			props := make([]*models.NameAndValue, 0)
			if strings.HasPrefix(name, "system.") {
				props = device.SystemProperties
			} else if strings.HasPrefix(name, "auto.") {
				props = device.AutoProperties
			}
			for idx, prop := range props {
				if *prop.Name == name && (value != "" || action.Is(enums.Delete)) {
					if *prop.Name == constants.K8sSystemCategoriesPropertyKey {
						value2 := getUpdatedSystemCategories(*prop.Value, value, action)
						*prop.Value = value2
					} else {
						if action.Is(enums.Delete) {
							props[idx] = props[len(props)-1]
							props = props[:len(props)-1]
						} else {
							*prop.Value = value
						}
					}
				}
			}
			if strings.HasPrefix(name, "system.") {
				device.SystemProperties = props
			} else if strings.HasPrefix(name, "auto.") {
				device.AutoProperties = props
			}
		}
		for idx, prop := range device.CustomProperties {
			if *prop.Name == name && (value != "" || action.Is(enums.Delete)) {
				if *prop.Name == constants.K8sSystemCategoriesPropertyKey {
					value = getUpdatedSystemCategories(*prop.Value, value, action)
					*prop.Value = value
					return
				}
				if action.Is(enums.Delete) {
					device.CustomProperties[idx] = device.CustomProperties[len(device.CustomProperties)-1]
					device.CustomProperties = device.CustomProperties[:len(device.CustomProperties)-1]
					return
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
			logrus.Warnf("Custom property value is empty for %q, skipping", name)
		}
	}
}

// nolint: predeclared
func getUpdatedSystemCategories(oldValue, newValue string, action enums.BuilderAction) string {
	// we do not use strings.contain, because it may be matched as substring of some prop
	oldValues := strings.Split(strings.TrimSpace(oldValue), ",")
	for idx, ov := range oldValues {
		if ov == newValue {
			if action.Is(enums.Delete) {
				oldValues[idx] = oldValues[len(oldValues)-1]
				oldValues = oldValues[:len(oldValues)-1]
				oldValue = strings.Join(oldValues, ",")
			}
			return oldValue
		}
	}
	oldValue = oldValue + "," + newValue

	return oldValue
}

// AddFuncWithDefaults add
func (b *Builder) AddFuncWithDefaults(cache types.ResourceCache, configurer types.WatcherConfigurer, actions types.Actions) func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		objectMeta := rt.ObjectMeta(obj)
		options := b.getDefaultsDeviceOptions(rt, objectMeta, conf)
		additionalOptions, err := configurer.AddFuncOptions()(lctx, rt, obj, b)
		if err != nil {
			log.Errorf("failed to get device additional options: %s", err)

			return
		}

		options = append(options, additionalOptions...)
		actions.AddFunc()(lctx, rt, obj, options...)
	}
}

// UpdateFuncWithDefaults update
func (b *Builder) UpdateFuncWithDefaults(
	target func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, []types.DeviceOption, []types.DeviceOption),
) func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		objectMeta := rt.ObjectMeta(newObj)
		oldObjectMeta := rt.ObjectMeta(oldObj)

		options := b.getDefaultsDeviceOptions(rt, objectMeta, conf)
		oldObjOptions := b.getDefaultsDeviceOptions(rt, oldObjectMeta, conf)

		target(lctx, rt, oldObj, newObj, oldObjOptions, options)
	}
}

// DeleteFuncWithDefaults delete
func (b *Builder) DeleteFuncWithDefaults(
	configurer types.WatcherConfigurer,
	deleteFun func(*lmctx.LMContext, enums.ResourceType, interface{}, ...types.DeviceOption),
) func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		objectMeta := rt.ObjectMeta(obj)
		options := b.getDefaultsDeviceOptions(rt, objectMeta, conf)
		additionalOptions := configurer.DeleteFuncOptions()(lctx, rt, obj)
		options = append(options, additionalOptions...)
		deleteFun(lctx, rt, obj, options...)
	}
}

// MarkDeleteFunc mark
func (b *Builder) MarkDeleteFunc(
	configurer types.WatcherConfigurer,
	updateFun func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ...types.DeviceOption),
) func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		objectMeta := rt.ObjectMeta(obj)
		options := b.getDefaultsDeviceOptions(rt, objectMeta, conf)
		additionalOptions := configurer.DeleteFuncOptions()(lctx, rt, obj)
		options = append(options, additionalOptions...)
		options = append(options, b.GetMarkDeleteOptions(lctx, rt, objectMeta)...)
		updateFun(lctx, rt, obj, obj, options...)
	}
}

func (b *Builder) getDefaultsDeviceOptions(rt enums.ResourceType, objectMeta *metav1.ObjectMeta, conf *config.Config) []types.DeviceOption {
	options := []types.DeviceOption{
		b.Name(rt.LMName(objectMeta)),
		b.ResourceLabels(objectMeta.Labels),
		b.DisplayName(util.GetDisplayNameNew(rt, objectMeta, conf)),
		b.SystemCategory(rt.GetCategory(), enums.Add),
		b.Auto("name", objectMeta.Name),
		b.Auto("selflink", util.SelfLink(rt.IsNamespaceScopedResource(), rt.K8SAPIVersion(), rt.String(), *objectMeta)),
		b.Auto("uid", string(objectMeta.UID)),
		b.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(objectMeta.CreationTimestamp.Unix(), 10)),
	}
	if rt.IsNamespaceScopedResource() {
		options = append(options, b.Auto("namespace", objectMeta.Namespace))
	}

	return options
}

// GetMarkDeleteOptions mark delete
func (b *Builder) GetMarkDeleteOptions(lctx *lmctx.LMContext, rt enums.ResourceType, meta *metav1.ObjectMeta) []types.DeviceOption {
	if meta.DeletionTimestamp == nil {
		t := metav1.Now()
		meta.DeletionTimestamp = &t
	}

	return []types.DeviceOption{
		b.SystemCategory(rt.GetDeletedCategory(), enums.Add),
		b.DeletedOn(meta.DeletionTimestamp.Time),
		b.ChangePrimaryKeysToMarkDelete(),
	}
}

// ChangePrimaryKeysToMarkDelete change
func (b *Builder) ChangePrimaryKeysToMarkDelete() types.DeviceOption {
	return func(device *models.Device) {
		if device == nil {
			return
		}
		shortUUID := strconv.FormatUint(uint64(util.GetShortUUID()), 10)
		deviceName := util.TrimName(*device.Name)
		b.Name(deviceName + "-" + shortUUID)(device)
		b.DisplayName(*device.DisplayName + "-" + shortUUID)(device)
	}
}
