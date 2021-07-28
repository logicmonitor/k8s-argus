package resource_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resource"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

var resourceName = "test-resource"

func TestBuilresourceWithExistingDeviceInput(t *testing.T) {
	t.Parallel()
	manager := resource.Manager{} // nolint: exhaustivestruct
	conf := &config.Config{       // nolint: exhaustivestruct
		Address:         "address",
		ClusterName:     "clusterName",
		DeleteResources: false,
		DisableAlerting: true,
		ClusterGroupID:  123,
		ProxyURL:        "url",
	}

	options := []types.ResourceOption{
		manager.Name("Name"),
		manager.DisplayName("DisplayName"),
		manager.SystemCategory("category", enums.Add),
	}

	inputresource := getSampleDevice()
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"resource_id": "build_resource_test"}))
	resourceObj, err := util.BuildResource(lctx, conf, inputresource, options...)
	if err != nil {
		t.Fail()
		return
	}

	if inputresource.Name != resourceObj.Name {
		t.Errorf("TestBuilresourceWithExistingDeviceInput - Error building resource %v", resourceObj.Name)
	}
}

func getSampleDevice() *models.Device {
	return &models.Device{ // nolint: exhaustivestruct
		Name:        &resourceName,
		DisplayName: &resourceName,
	}
}
