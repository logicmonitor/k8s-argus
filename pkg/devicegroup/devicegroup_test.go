package devicegroup

import (
	"fmt"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"testing"
)

var testAndString = `hasCategory("foo") && auto.bar == "baz"`
var testOrString = `hasCategory("foo") || auto.bar == "baz"`

func TestAppliesToBuilder(t *testing.T) {
	builder := NewAppliesToBuilder().HasCategory("foo").And().Auto("bar").Equals("baz")
	if builder.String() != testAndString {
		t.Errorf("appliesTo string is invalid: %s", builder.String())
	}

	builder = NewAppliesToBuilder().HasCategory("foo").Or().Auto("bar").Equals("baz")
	if builder.String() != testOrString {
		t.Errorf("appliesTo string is invalid: %s", builder.String())
	}
}

func TestAddDeviceGroup(t *testing.T) {
	// define the client
	config := client.NewConfig()
	accessID := "Jz2EuQQKMJw9n5x3Xg9r"
	config.SetAccessID(&accessID)
	accessKey := "SyV52F6+D76({^q(]bcjj+R2XkV_HHVX36+)7t2i"
	config.SetAccessKey(&accessKey)
	domain := "jeremy.logicmonitor.com:8080"
	config.SetAccountDomain(&domain)
	config.TransportCfg.WithSchemes([]string{"http"})

	// this is only used for QA environment where the https is invalid
	// config.TransportCfg.WithSchemes([]string{"http"})
	// this is only used for QA environment where the https is invalid

	lmSdk := client.New(config)

	fmt.Println(addDeviceGroup(lmSdk))
}
func addDeviceGroup(lmSdk *client.LMSdkGo) (string, error) {
	params := lm.NewAddDeviceGroupParams()
	name := "Kubernetes Cluster: aaa12"
	params.SetBody(&models.DeviceGroup{
		Name:            &name,
		Description:     "A dynamic device group for Kubernetes.",
		ParentID:        0,
		AppliesTo:       "hasCategory(\"KubernetesCluster\") && auto.clustername == \"cluster-jeremy\"",
		DisableAlerting: true,
	})
	restResponse, err := lmSdk.LM.AddDeviceGroup(params)
	if err != nil {
		return "", fmt.Errorf("failed to add device group %q: %v", name, err)
	}

	deviceGroup := restResponse.Payload

	return deviceGroup.Description, nil
}
