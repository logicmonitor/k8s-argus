package dgbuilder

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

const (
	hasCategoryOpen = "hasCategory("
	existsOpen      = "exists("
	closingBracket  = ")"
)

// Builder implements types.ResourceBuilder
type Builder struct{}

// GroupName implements types.ResourceBuilder.
func (b *Builder) GroupName(name string) types.ResourceGroupOption {
	return func(resourceGroup *models.DeviceGroup) {
		resourceGroup.Name = &name
	}
}

// Description implements types.ResourceBuilder.
func (b *Builder) Description(description string) types.ResourceGroupOption {
	return func(resourceGroup *models.DeviceGroup) {
		resourceGroup.Description = description
	}
}

// ParentID implements types.ResourceBuilder.
func (b *Builder) ParentID(parentID int32) types.ResourceGroupOption {
	return func(resourceGroup *models.DeviceGroup) {
		resourceGroup.ParentID = parentID
	}
}

// AppliesTo implements types.ResourceBuilder.
func (b *Builder) AppliesTo(appliesTo types.AppliesToBuilder) types.ResourceGroupOption {
	return func(resourceGroup *models.DeviceGroup) {
		resourceGroup.AppliesTo = appliesTo.Build()
	}
}

// CustomProperties implements types.ResourceBuilder.
func (b *Builder) CustomProperties(appliesTo types.PropertyBuilder) types.ResourceGroupOption {
	return func(resourceGroup *models.DeviceGroup) {
		resourceGroup.CustomProperties = appliesTo.Build(resourceGroup.CustomProperties)
	}
}

// DisableAlerting implements types.ResourceBuilder.
func (b *Builder) DisableAlerting(disableAlerting bool) types.ResourceGroupOption {
	return func(resourceGroup *models.DeviceGroup) {
		resourceGroup.DisableAlerting = disableAlerting
	}
}

// Custom implements types.ResourceBuilder.
func (b *Builder) Custom(name, value string) types.ResourceGroupOption {
	return setProperty(name, value, enums.Add)
}

// nolint: predeclared,gocognit,cyclop
func setProperty(name, value string, action enums.BuilderAction) types.ResourceGroupOption {
	return func(resource *models.DeviceGroup) {
		if resource == nil {
			return
		}
		if resource.CustomProperties == nil {
			resource.CustomProperties = []*models.NameAndValue{}
		}
		for idx, prop := range resource.CustomProperties {
			if *prop.Name == name && (value != "" || action.Is(enums.Delete)) {
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

// AppliesToBuilderImpl impl
type AppliesToBuilderImpl struct {
	value string
}

// NewAppliesToBuilder is the builder for appliesTo.
func NewAppliesToBuilder() types.AppliesToBuilder {
	return &AppliesToBuilderImpl{value: ""} // nolint: exhaustivestruct
}

// And and
func (a *AppliesToBuilderImpl) And() types.AppliesToBuilder {
	a.value += " && "

	return a
}

// OpenBracket open (
func (a *AppliesToBuilderImpl) OpenBracket() types.AppliesToBuilder {
	a.value += " ( "

	return a
}

// TrimOrCloseBracket removes last or and close )
func (a *AppliesToBuilderImpl) TrimOrCloseBracket() types.AppliesToBuilder {
	a.value = strings.TrimSuffix(a.value, " || ")
	a.value += " ) "

	return a
}

// Or or
func (a *AppliesToBuilderImpl) Or() types.AppliesToBuilder {
	a.value += " || "

	return a
}

// Equals equals
func (a *AppliesToBuilderImpl) Equals(val string) types.AppliesToBuilder {
	a.value += " == " + fmt.Sprintf(`"%s"`, val)

	return a
}

// HasCategory has
func (a *AppliesToBuilderImpl) HasCategory(category string) types.AppliesToBuilder {
	a.value += hasCategoryOpen + fmt.Sprintf(`"%s"`, category) + closingBracket

	return a
}

// Exists exists
func (a *AppliesToBuilderImpl) Exists(property string) types.AppliesToBuilder {
	a.value += existsOpen + fmt.Sprintf(`"%s"`, property) + closingBracket

	return a
}

// Auto auto
func (a *AppliesToBuilderImpl) Auto(property string) types.AppliesToBuilder {
	a.value += "auto." + property

	return a
}

// Custom custom
func (a *AppliesToBuilderImpl) Custom(property string) types.AppliesToBuilder {
	a.value += property

	return a
}

// Build string
func (a *AppliesToBuilderImpl) Build() string {
	return a.value
}

type propertyBuilder struct {
	properties []config.PropOpts
}

// NewPropertyBuilder is the builder for properties
func NewPropertyBuilder() types.PropertyBuilder {
	return &propertyBuilder{} // nolint: exhaustivestruct
}

func (p *propertyBuilder) Add(key string, value string, override bool) types.PropertyBuilder {
	opts := config.PropOpts{
		Name:     key,
		Value:    value,
		Override: &override,
	}
	p.properties = append(p.properties, opts)

	return p
}

func (p *propertyBuilder) AddProperties(properties []config.PropOpts) types.PropertyBuilder {
	p.properties = append(p.properties, properties...)

	return p
}

func (p *propertyBuilder) Build(existingProps []*models.NameAndValue) []*models.NameAndValue {
	exProps := make(map[string]int)

	for idx, elm := range existingProps {
		exProps[*elm.Name] = idx
	}

	for _, prop := range p.properties {
		key := prop.Name
		value := prop.Value
		override := true
		if prop.Override != nil {
			override = *prop.Override
		}
		val, ok := exProps[key]
		if ok && override {
			existingProps[val].Value = &value
		} else if !ok {
			existingProps = append(existingProps, &models.NameAndValue{Name: &key, Value: &value})
		}
	}

	return existingProps
}
