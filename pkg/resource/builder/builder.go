package builder

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Builder implements types.ResourceBuilder
type Builder struct{}

// Name implements types.ResourceBuilder.
func (b *Builder) Name(name string) types.ResourceOption {
	return func(resource *models.Device) {
		resource.Name = &name
	}
}

// DisplayName implements types.ResourceBuilder.
func (b *Builder) DisplayName(name string) types.ResourceOption {
	return func(resource *models.Device) {
		resource.DisplayName = &name
	}
}

// CollectorID implements types.ResourceBuilder.
func (b *Builder) CollectorID(id int32) types.ResourceOption {
	return func(resource *models.Device) {
		resource.PreferredCollectorID = &id
	}
}

// SystemCategory implements types.ResourceBuilder
// nolint: predeclared
func (b *Builder) SystemCategory(category string, action enums.BuilderAction) types.ResourceOption {
	return setProperty(constants.K8sSystemCategoriesPropertyKey, category, action)
}

// DisableResourceAlerting implements types.ResourceBuilder
func (b *Builder) DisableResourceAlerting(disable bool) types.ResourceOption {
	return func(device *models.Device) {
		device.DisableAlerting = disable
	}
}

// ResourceAnnotations implements types.ResourceBuilder
func (b *Builder) ResourceAnnotations(properties map[string]string) types.ResourceOption {
	return func(resource *models.Device) {
		if resource == nil {
			return
		}
		if resource.CustomProperties == nil {
			resource.CustomProperties = []*models.NameAndValue{}
		}
		for name, value := range properties {
			propName := constants.AnnotationCustomPropertyPrefix + name
			propValue := value
			if propValue == "" {
				propValue = constants.LabelNullPlaceholder
			}
			existed := false
			for _, prop := range resource.CustomProperties {
				if *prop.Name == propName {
					*prop.Value = propValue
					existed = true

					break
				}
			}
			if !existed {
				resource.CustomProperties = append(resource.CustomProperties, &models.NameAndValue{
					Name:  &propName,
					Value: &propValue,
				})
			}
		}
	}
}

// ResourceLabels implements types.ResourceBuilder
func (b *Builder) ResourceLabels(properties map[string]string) types.ResourceOption {
	return func(resource *models.Device) {
		if resource == nil {
			return
		}
		if resource.CustomProperties == nil {
			resource.CustomProperties = []*models.NameAndValue{}
		}
		for name, value := range properties {
			propName := constants.LabelCustomPropertyPrefix + name
			propValue := value
			if propValue == "" {
				propValue = constants.LabelNullPlaceholder
			}
			existed := false
			for _, prop := range resource.CustomProperties {
				if *prop.Name == propName {
					*prop.Value = propValue
					existed = true

					break
				}
			}
			if !existed {
				resource.CustomProperties = append(resource.CustomProperties, &models.NameAndValue{
					Name:  &propName,
					Value: &propValue,
				})
			}
		}
	}
}

// Auto implements types.ResourceBuilder
func (b *Builder) Auto(name, value string) types.ResourceOption {
	return setProperty("auto."+name, value, enums.Add)
}

// System implements types.ResourceBuilder.
func (b *Builder) System(name, value string) types.ResourceOption {
	return setProperty("system."+name, value, enums.Add)
}

// Custom implements types.ResourceBuilder.
func (b *Builder) Custom(name, value string) types.ResourceOption {
	return setProperty(name, value, enums.Add)
}

// DeletedOn implements types.ResourceBuilder
func (b *Builder) DeletedOn(value time.Time) types.ResourceOption {
	return func(resource *models.Device) {
		name := constants.K8sResourceDeletedOnPropertyKey

		if resource == nil {
			return
		}

		if resource.CustomProperties == nil {
			resource.CustomProperties = []*models.NameAndValue{}
		}

		for _, prop := range resource.CustomProperties {
			if *prop.Name == name {
				return
			}
		}

		strValue := strconv.FormatInt(value.Unix(), 10)
		resource.CustomProperties = append(resource.CustomProperties, &models.NameAndValue{
			Name:  &name,
			Value: &strValue,
		})
	}
}

// nolint: predeclared,gocognit,cyclop
func setProperty(name, value string, action enums.BuilderAction) types.ResourceOption {
	return func(resource *models.Device) {
		if resource == nil {
			return
		}
		if resource.CustomProperties == nil {
			resource.CustomProperties = []*models.NameAndValue{}
		}
		if action.Is(enums.Delete) {
			props := make([]*models.NameAndValue, 0)
			if strings.HasPrefix(name, "system.") {
				props = resource.SystemProperties
			} else if strings.HasPrefix(name, "auto.") {
				props = resource.AutoProperties
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
				resource.SystemProperties = props
			} else if strings.HasPrefix(name, "auto.") {
				resource.AutoProperties = props
			}
		}
		for idx, prop := range resource.CustomProperties {
			if *prop.Name == name && (value != "" || action.Is(enums.Delete)) {
				if *prop.Name == constants.K8sSystemCategoriesPropertyKey {
					value = getUpdatedSystemCategories(*prop.Value, value, action)
					*prop.Value = value
					return
				}
				if action.Is(enums.Delete) {
					resource.CustomProperties[idx] = resource.CustomProperties[len(resource.CustomProperties)-1]
					resource.CustomProperties = resource.CustomProperties[:len(resource.CustomProperties)-1]
					return
				}
				*prop.Value = value
				return
			}
		}
		if value != "" {
			resource.CustomProperties = append(resource.CustomProperties, &models.NameAndValue{
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
func (b *Builder) AddFuncWithDefaults(
	configurer types.WatcherConfigurer,
	actions types.Actions,
) types.AddPreprocessFunc {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig(lctx)
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		accessor, err := meta.Accessor(obj)
		if err != nil {
			log.Errorf("failed to get objectmeta: %s", err)
			return
		}
		options := b.GetDefaultsResourceOptions(rt, meta.AsPartialObjectMetadata(accessor), conf)
		additionalOptions, err := configurer.AddFuncOptions()(lctx, rt, obj, b)
		if err != nil {
			if errors.Is(err, aerrors.ErrPodSucceeded) {
				log.Warnf("pod having succeeded status will not be considered for monitoring: %s", err)
			} else {
				log.Errorf("failed to get resource additional options: %s", err)
			}
			return
		}

		options = append(options, additionalOptions...)
		actions.AddFunc()(lctx, rt, obj, options...) // nolint: errcheck
	}
}

// UpdateFuncWithDefaults update
func (b *Builder) UpdateFuncWithDefaults(
	target types.UpdateProcessFunc,
) types.UpdatePreprocessFunc {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj, newObj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig(lctx)
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		objectMeta, _ := rt.ObjectMeta(newObj)
		oldObjectMeta, _ := rt.ObjectMeta(oldObj)

		options := b.GetDefaultsResourceOptions(rt, objectMeta, conf)
		oldObjOptions := b.GetDefaultsResourceOptions(rt, oldObjectMeta, conf)

		target(lctx, rt, oldObj, newObj, oldObjOptions, options)
	}
}

// DeleteFuncWithDefaults delete
func (b *Builder) DeleteFuncWithDefaults(
	configurer types.WatcherConfigurer,
	deleteFun types.ExecDeleteFunc,
) types.DeletePreprocessFunc {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig(lctx)
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		objectMeta, _ := rt.ObjectMeta(obj)
		options := b.GetDefaultsResourceOptions(rt, objectMeta, conf)
		additionalOptions := configurer.DeleteFuncOptions()(lctx, rt, obj)
		options = append(options, additionalOptions...)
		deleteFun(lctx, rt, obj, options...) // nolint: errcheck
	}
}

// MarkDeleteFunc mark
func (b *Builder) MarkDeleteFunc(
	configurer types.WatcherConfigurer,
	updateFun func(*lmctx.LMContext, enums.ResourceType, interface{}, interface{}, ...types.ResourceOption),
) func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	return func(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
		log := lmlog.Logger(lctx)
		conf, err := config.GetConfig(lctx)
		if err != nil {
			log.Errorf("Failed to get config: %s", err)

			return
		}
		objectMeta, _ := rt.ObjectMeta(obj)
		options := b.GetDefaultsResourceOptions(rt, objectMeta, conf)
		additionalOptions := configurer.DeleteFuncOptions()(lctx, rt, obj)
		options = append(options, additionalOptions...)
		options = append(options, b.GetMarkDeleteOptions(lctx, rt, objectMeta)...)
		updateFun(lctx, rt, obj, obj, options...)
	}
}

// GetDefaultsResourceOptions returns default options for resource
func (b *Builder) GetDefaultsResourceOptions(rt enums.ResourceType, objectMeta *metav1.PartialObjectMetadata, conf *config.Config) []types.ResourceOption {
	if objectMeta == nil {
		return []types.ResourceOption{}
	}
	options := []types.ResourceOption{
		b.Name(rt.LMName(objectMeta)),
		b.ResourceLabels(objectMeta.Labels),
		b.ResourceAnnotations(objectMeta.Annotations),
		b.DisplayName(util.GetDisplayName(rt, objectMeta, conf)),
		b.SystemCategory(rt.GetCategory(), enums.Add),
		b.Auto("name", objectMeta.Name),
		b.Auto("selflink", util.SelfLink(rt.IsNamespaceScopedResource(), rt.K8SAPIVersion(), rt.String(), objectMeta)),
		b.Auto("uid", string(objectMeta.UID)),
		b.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(objectMeta.CreationTimestamp.Unix(), 10)),
	}
	if rt.IsNamespaceScopedResource() {
		options = append(options, b.Auto("namespace", objectMeta.Namespace))
	}

	return options
}

// GetMarkDeleteOptions mark delete
func (b *Builder) GetMarkDeleteOptions(lctx *lmctx.LMContext, rt enums.ResourceType, meta *metav1.PartialObjectMetadata) []types.ResourceOption {
	if meta.DeletionTimestamp == nil {
		t := metav1.Now()
		meta.DeletionTimestamp = &t
	}

	options := []types.ResourceOption{
		b.SystemCategory(rt.GetDeletedCategory(), enums.Add),
		b.DisableResourceAlerting(true),
		b.DeletedOn(meta.DeletionTimestamp.Time),
		b.ChangePrimaryKeysToMarkDelete(),
	}
	if val, ok := meta.Labels["logicmonitor/deleteafterduration"]; ok {
		options = append(options, b.Custom("kubernetes.resourcedeleteafterduration", val))
	}
	// We are not deleting argus pod, as we need argus pod logs for troubleshooting
	if util.IsArgusPodObject(lctx, rt, meta) {
		// defaults to 10 days
		scheduledDeleteTime := "P10DT0H0M0S"
		conf, err := config.GetConfig(lctx)
		if err == nil {
			scheduledDeleteTime = *conf.DeleteInfraPodsAfter
		}
		options = append(options, b.Custom("kubernetes.resourcedeleteafterduration", scheduledDeleteTime))
	}
	return options
}

// ChangePrimaryKeysToMarkDelete change
func (b *Builder) ChangePrimaryKeysToMarkDelete() types.ResourceOption {
	return func(resource *models.Device) {
		if resource == nil {
			return
		}
		shortUUID := strconv.FormatUint(uint64(util.GetShortUUID()), 10)
		resourceName := util.TrimName(*resource.Name)
		b.Name(resourceName + "-" + shortUUID)(resource)
		b.DisplayName(*resource.DisplayName + "-" + shortUUID)(resource)
	}
}
