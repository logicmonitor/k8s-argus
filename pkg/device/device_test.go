package device

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vkumbhar94/lm-sdk-go/models"
)

var (
	deviceName             = "test-device"
	customPropertiesName1  = "name1"
	customPropertiesValue1 = "value1"
	customPropertiesName2  = "name2"
	customPropertiesValue2 = "value2"

	systemPropertiesName1  = "system-name1"
	systemPropertiesValue1 = "system-value1"
	systemPropertiesName2  = "system-name2"
	systemPropertiesValue2 = "system-value2"
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

	inputdevice := getSampleDevice()
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": "build_device_test"}))
	device := buildDevice(lctx, config, inputdevice, options...)

	if inputdevice.Name != device.Name {
		t.Errorf("TestBuildeviceWithExistingDeviceInput - Error building device %v", device.Name)
	}
}

func TestFindByDisplayNamesWithEmptyDisplayNames(t *testing.T) {
	displayNames := []string{}
	manager := Manager{}

	expectedDevice := []*models.Device{}
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": "build_device_test"}))
	actualDevice, err := manager.FindByDisplayNames(lctx, "pods", displayNames...)
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
