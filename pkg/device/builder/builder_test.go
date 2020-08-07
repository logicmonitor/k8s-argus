package builder

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/lm-sdk-go/models"
)

func TestBuilder_SetProperty(t *testing.T) {
	propName := "name"
	propValue1 := "value1"
	propValue2 := "value2"

	device := &models.Device{
		CustomProperties: []*models.NameAndValue{},
	}

	setProp := setProperty(propName, propValue1)
	setProp(device)
	if propValue1 != getDevicePropValueByName(device, propName) {
		t.Errorf("failed to set prop %s:%s to the device", propName, propValue1)
	}
	setProp = setProperty(propName, propValue2)
	setProp(device)
	if propValue2 != getDevicePropValueByName(device, propName) {
		t.Errorf("failed to set prop %s:%s to the device", propName, propValue2)
	}

	sysPropValue1 := "k1=v1,k2=v2"
	sysPropValue2 := constants.PodCategory

	setProp = setProperty(constants.K8sSystemCategoriesPropertyKey, sysPropValue1)
	setProp(device)
	if propValue2 != getDevicePropValueByName(device, propName) {
		t.Errorf("failed to set prop %s:%s to the device", propName, propValue2)
	}
	if sysPropValue1 != getDevicePropValueByName(device, constants.K8sSystemCategoriesPropertyKey) {
		t.Errorf("failed to set prop %s:%s to the device", constants.K8sSystemCategoriesPropertyKey, sysPropValue1)
	}
	setProp = setProperty(constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
	setProp(device)
	if sysPropValue1+","+constants.PodCategory != getDevicePropValueByName(device, constants.K8sSystemCategoriesPropertyKey) {
		t.Errorf("failed to set prop %s:%s to the device", constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
	}
	setProp = setProperty(constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
	setProp(device)
	if sysPropValue1+","+constants.PodCategory != getDevicePropValueByName(device, constants.K8sSystemCategoriesPropertyKey) {
		t.Errorf("failed to set prop %s:%s to the device", constants.K8sSystemCategoriesPropertyKey, sysPropValue2)
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
