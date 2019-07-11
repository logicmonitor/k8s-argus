package device

import (
	"testing"

	"github.com/logicmonitor/lm-sdk-go/models"
)

func TestGetPropertyValue(t *testing.T) {
	deviceName := "test-device"

	customPropertiesName1 := "name1"
	customPropertiesValue1 := "value1"
	customPropertiesName2 := "name2"
	customPropertiesValue2 := "value2"

	systemPropertiesName1 := "system-name1"
	systemPropertiesValue1 := "system-value1"
	systemPropertiesName2 := "system-name2"
	systemPropertiesValue2 := "system-value2"

	device := &models.Device{
		Name:        &deviceName,
		DisplayName: &deviceName,
		CustomProperties: []*models.NameAndValue{
			{
				Name:  &customPropertiesName1,
				Value: &customPropertiesValue1,
			}, {
				Name:  &customPropertiesName2,
				Value: &customPropertiesValue2,
			},
		},
		SystemProperties: []*models.NameAndValue{
			{
				Name:  &systemPropertiesName1,
				Value: &systemPropertiesValue1,
			}, {
				Name:  &systemPropertiesName2,
				Value: &systemPropertiesValue2,
			},
		},
	}

	manage := Manager{}
	value := manage.GetPropertyValue(device, customPropertiesName1)
	t.Logf("name=%s, value=%s", customPropertiesName1, value)
	value = manage.GetPropertyValue(device, systemPropertiesName2)
	t.Logf("name=%s, value=%s", systemPropertiesName2, value)
	value = manage.GetPropertyValue(device, "non-exist-name")
	t.Logf("name=%s, value=%s", "non-exist-name", value)
}
