package device

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/stretchr/testify/assert"
)

var (
	deviceName             string = "test-device"
	customPropertiesName1  string = "name1"
	customPropertiesValue1 string = "value1"
	customPropertiesName2  string = "name2"
	customPropertiesValue2 string = "value2"

	systemPropertiesName1  string = "system-name1"
	systemPropertiesValue1 string = "system-value1"
	systemPropertiesName2  string = "system-name2"
	systemPropertiesValue2 string = "system-value2"
)

func TestBuildeviceWithExistingDeviceInput(t *testing.T) {
	manager := Manager{}
	config := &config.Config{
		Address:         "address",
		ClusterCategory: "category",
		ClusterName:     "clusterName",
		Debug:           false,
		DeleteDevices:   false,
		DisableAlerting: true,
		ClusterGroupID:  123,
		ProxyURL:        "url",
	}

	options := []types.DeviceOption{
		manager.Name("Name"),
		manager.DisplayName("DisplayName"),
		manager.SystemCategories("catgory"),
	}

	controllerClient := &manager.ControllerClient

	inputdevice := getSampleDevice()
	device := buildDevice(config, *controllerClient, inputdevice, options...)

	if inputdevice.Name != device.Name {
		t.Errorf("Error building device %v", device.Name)
	}
}

func TestFindByDisplayNamesWithEmptyDisplayNames(t *testing.T) {
	displayNames := []string{}
	manager := Manager{}

	expectedDevice := []*models.Device{}
	actualDevice, err := manager.FindByDisplayNames(displayNames...)
	assert.Nil(t, err)
	assert.Equal(t, expectedDevice, actualDevice)
}

func TestGetPropertyValue(t *testing.T) {
	device := getSampleDevice()
	manage := Manager{}
	value := manage.GetPropertyValue(device, customPropertiesName1)
	t.Logf("name=%s, value=%s", customPropertiesName1, value)
	value = manage.GetPropertyValue(device, systemPropertiesName2)
	t.Logf("name=%s, value=%s", systemPropertiesName2, value)
	value = manage.GetPropertyValue(device, "non-exist-name")
	t.Logf("name=%s, value=%s", "non-exist-name", value)
}

func getSampleDevice() *models.Device {
	return &models.Device{
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
}
