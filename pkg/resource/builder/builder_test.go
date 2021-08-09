package builder_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/resource/builder"
	"github.com/logicmonitor/lm-sdk-go/models"
)

var (
	propName   = "name"
	prop2Name  = "name2"
	propValue1 = "value1"
	propValue2 = "value2"
)

func TestResourceLabels_NilDevice(t *testing.T) {
	t.Parallel()
	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	b := builder.Builder{} // nolint: exhaustivestruct

	resourceLabel := b.ResourceLabels(properties)
	resourceLabel(nil)
	if resourceLabel == nil {
		t.Errorf("TestResourceLabels_NilDevice - invalid inputs")
	}
}

func TestResourceLabels_NilCustomProperties(t *testing.T) {
	t.Parallel()
	resource := &models.Device{
		CustomProperties: []*models.NameAndValue{},
	}

	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	b := builder.Builder{} // nolint: exhaustivestruct

	resourceLabel := b.ResourceLabels(properties)
	resourceLabel(resource)

	kubernetesPropName := constants.LabelCustomPropertyPrefix + propName
	if propValue1 != getResourcePropValueByName(resource, kubernetesPropName) {
		t.Errorf("TestResourceLabels_NilCustomProperties - failed to set resource property %s", kubernetesPropName)
	}

	kubernetesPropName2 := constants.LabelCustomPropertyPrefix + prop2Name
	if propValue2 != getResourcePropValueByName(resource, kubernetesPropName2) {
		t.Errorf("TestResourceLabels_NilCustomProperties - failed to set resource property %s", kubernetesPropName2)
	}
}

func TestResourceLabels_ExistingCustomProperties(t *testing.T) {
	t.Parallel()
	resource := &models.Device{
		CustomProperties: []*models.NameAndValue{
			{
				Name:  &propName,
				Value: &propValue1,
			}, {
				Name:  &prop2Name,
				Value: &propValue2,
			},
		},
	}

	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	b := builder.Builder{} // nolint: exhaustivestruct
	resourceLabel := b.ResourceLabels(properties)
	resourceLabel(resource)
	if propValue1 != getResourcePropValueByName(resource, propName) {
		t.Errorf("TestResourceLabels_ExistingCustomProperties - failed to set resource property %s", propName)
	}
	if propValue2 != getResourcePropValueByName(resource, prop2Name) {
		t.Errorf("TestResourceLabels_ExistingCustomProperties- failed to set resource property %s", prop2Name)
	}
}

func TestResourceAnnotations_NilDevice(t *testing.T) {
	t.Parallel()
	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	b := builder.Builder{} // nolint: exhaustivestruct

	resourceAnnotation := b.ResourceAnnotations(properties)
	resourceAnnotation(nil)
	if resourceAnnotation == nil {
		t.Errorf("TestResourceAnnotations_NilDevice - invalid inputs")
	}
}

func TestResourceAnnotations_NilCustomProperties(t *testing.T) {
	t.Parallel()
	resource := &models.Device{
		CustomProperties: []*models.NameAndValue{},
	}

	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	b := builder.Builder{} // nolint: exhaustivestruct

	resourceAnnotation := b.ResourceAnnotations(properties)
	resourceAnnotation(resource)

	kubernetesPropName := constants.AnnotationCustomPropertyPrefix + propName
	if propValue1 != getResourcePropValueByName(resource, kubernetesPropName) {
		t.Errorf("TestResourceAnnotations_NilCustomProperties - failed to set resource property %s", kubernetesPropName)
	}

	kubernetesPropName2 := constants.AnnotationCustomPropertyPrefix + prop2Name
	if propValue2 != getResourcePropValueByName(resource, kubernetesPropName2) {
		t.Errorf("TestResourceAnnotations_NilCustomProperties - failed to set resource property %s", kubernetesPropName2)
	}
}

func TestResourceAnnotations_ExistingCustomProperties(t *testing.T) {
	t.Parallel()
	resource := &models.Device{
		CustomProperties: []*models.NameAndValue{
			{
				Name:  &propName,
				Value: &propValue1,
			}, {
				Name:  &prop2Name,
				Value: &propValue2,
			},
		},
	}

	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	b := builder.Builder{} // nolint: exhaustivestruct
	resourceAnnotation := b.ResourceAnnotations(properties)
	resourceAnnotation(resource)
	if propValue1 != getResourcePropValueByName(resource, propName) {
		t.Errorf("TestResourceAnnotations_ExistingCustomProperties - failed to set resource property %s", propName)
	}
	if propValue2 != getResourcePropValueByName(resource, prop2Name) {
		t.Errorf("TestResourceAnnotations_ExistingCustomProperties- failed to set resource property %s", prop2Name)
	}
}

func getResourcePropValueByName(d *models.Device, name string) string {
	if d == nil || d.CustomProperties == nil {
		return ""
	}
	for _, prop := range d.CustomProperties {
		if *prop.Name == name {
			return *prop.Value
		}
	}

	return ""
}
