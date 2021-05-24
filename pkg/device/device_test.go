package device_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

var deviceName = "test-device"

func TestBuildeviceWithExistingDeviceInput(t *testing.T) {
	t.Parallel()
	manager := device.Manager{} // nolint: exhaustivestruct
	conf := &config.Config{     // nolint: exhaustivestruct
		Address:         "address",
		ClusterCategory: "category",
		ClusterName:     "clusterName",
		DeleteDevices:   false,
		DisableAlerting: true,
		ClusterGroupID:  123,
		ProxyURL:        "url",
	}

	options := []types.DeviceOption{
		manager.Name("Name"),
		manager.DisplayName("DisplayName"),
		manager.SystemCategory("category", enums.Add),
	}

	inputdevice := getSampleDevice()
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": "build_device_test"}))
	deviceObj, err := util.BuildDevice(lctx, conf, inputdevice, options...)
	if err != nil {
		t.Fail()
		return
	}

	if inputdevice.Name != deviceObj.Name {
		t.Errorf("TestBuildeviceWithExistingDeviceInput - Error building device %v", deviceObj.Name)
	}
}

func getSampleDevice() *models.Device {
	return &models.Device{ // nolint: exhaustivestruct
		Name:        &deviceName,
		DisplayName: &deviceName,
	}
}
