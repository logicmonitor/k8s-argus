package builder

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/vkumbhar94/lm-sdk-go/models"
)

var (
	propName   = "name"
	prop2Name  = "name2"
	propValue1 = "value1"
	propValue2 = "value2"
)

func TestResourceLabels_NilDevice(t *testing.T) {
	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	builder := Builder{}

	resourceLabel := builder.ResourceLabels(properties)
	resourceLabel(nil)
	if resourceLabel == nil {
		t.Errorf("TestResourceLabels_NilDevice - invalid inputs")
	}
}

func TestResourceLabels_NilCustomProperties(t *testing.T) {
	device := &models.Device{
		CustomProperties: []*models.NameAndValue{},
	}

	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	builder := Builder{}

	resourceLabel := builder.ResourceLabels(properties)
	resourceLabel(device)

	kubernetesPropName := constants.LabelCustomPropertyPrefix + propName
	if propValue1 != getDevicePropValueByName(device, kubernetesPropName) {
		t.Errorf("TestResourceLabels_NilCustomProperties - failed to set device property %s", kubernetesPropName)
	}

	kubernetesPropName2 := constants.LabelCustomPropertyPrefix + prop2Name
	if propValue2 != getDevicePropValueByName(device, kubernetesPropName2) {
		t.Errorf("TestResourceLabels_NilCustomProperties - failed to set device property %s", kubernetesPropName2)
	}

}

func TestResourceLabels_ExistingCustomProperties(t *testing.T) {
	device := &models.Device{
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

	builder := Builder{}

	resourceLabel := builder.ResourceLabels(properties)
	resourceLabel(device)

	if propValue1 != getDevicePropValueByName(device, propName) {
		t.Errorf("TestResourceLabels_ExistingCustomProperties - failed to set device property %s", propName)
	}

	if propValue2 != getDevicePropValueByName(device, prop2Name) {
		t.Errorf("TestResourceLabels_ExistingCustomProperties- failed to set device property %s", prop2Name)
	}

}

func TestBuilder_SetProperty(t *testing.T) {
	device := &models.Device{
		CustomProperties: []*models.NameAndValue{},
	}

	setProp := setProperty(propName, propValue1)
	setProp(device)
	if propValue1 != getDevicePropValueByName(device, propName) {
		t.Errorf("TestBuilder_SetProperty - failed to set prop %s:%s to the device", propName, propValue1)
	}
	setProp = setProperty(propName, propValue2)
	setProp(device)
	if propValue2 != getDevicePropValueByName(device, propName) {
		t.Errorf("TestBuilder_SetProperty - failed to set prop %s:%s to the device", propName, propValue2)
	}

	sysPropValue1 := "k1=v1,k2=v2"
	sysPropValue2 := constants.PodCategory

	setProp = setProperty(constants.K8sSystemCategoriesPropertyKey, sysPropValue1)
	setProp(device)
	if propValue2 != getDevicePropValueByName(device, propName) {
		t.Errorf("TestBuilder_SetProperty - failed to set prop %s:%s to the device", propName, propValue2)
	}
	if sysPropValue1 != getDevicePropValueByName(device, constants.K8sSystemCategoriesPropertyKey) {
		t.Errorf("TestBuilder_SetProperty - failed to set prop %s:%s to the device", constants.K8sSystemCategoriesPropertyKey, sysPropValue1)
	}
	setProp = setProperty(constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
	setProp(device)
	if sysPropValue1+","+constants.PodCategory != getDevicePropValueByName(device, constants.K8sSystemCategoriesPropertyKey) {
		t.Errorf("TestBuilder_SetProperty - failed to set prop %s:%s to the device", constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
	}
	setProp = setProperty(constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
	setProp(device)
	if sysPropValue1+","+constants.PodCategory != getDevicePropValueByName(device, constants.K8sSystemCategoriesPropertyKey) {
		t.Errorf("TestBuilder_SetProperty - failed to set prop %s:%s to the device", constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
	}
}

func getDevicePropValueByName(d *models.Device, name string) string {
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
