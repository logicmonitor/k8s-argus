package builder_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
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
	device := &models.Device{
		CustomProperties: []*models.NameAndValue{},
	}

	properties := map[string]string{
		propName:  propValue1,
		prop2Name: propValue2,
	}

	b := builder.Builder{} // nolint: exhaustivestruct

	resourceLabel := b.ResourceLabels(properties)
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
	t.Parallel()
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

	b := builder.Builder{} // nolint: exhaustivestruct
	resourceLabel := b.ResourceLabels(properties)
	resourceLabel(device)
	if propValue1 != getDevicePropValueByName(device, propName) {
		t.Errorf("TestResourceLabels_ExistingCustomProperties - failed to set device property %s", propName)
	}
	if propValue2 != getDevicePropValueByName(device, prop2Name) {
		t.Errorf("TestResourceLabels_ExistingCustomProperties- failed to set device property %s", prop2Name)
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
