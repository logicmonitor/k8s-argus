package device

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	deviceName = "test-device"
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

func getSampleDevice() *models.Device {
	return &models.Device{
		Name:        &deviceName,
		DisplayName: &deviceName,
	}
}
