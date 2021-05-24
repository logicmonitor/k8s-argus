package utilities_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/stretchr/testify/assert"
)

var (
	deviceName             = "test-device"
	customPropertiesName1  = "auto.name"
	customPropertiesValue1 = "test-device"
	customPropertiesName2  = "auto.namespace"
	customPropertiesValue2 = "default"

	systemPropertiesName1  = "system-name1"
	systemPropertiesValue1 = "system-value1"
	systemPropertiesName2  = "system-name2"
	systemPropertiesValue2 = "system-value2"
)

func TestGetPropertyValue(t *testing.T) {
	t.Parallel()
	device := getDevice()
	value := util.GetPropertyValue(device, customPropertiesName1)
	t.Logf("name=%s, value=%s", customPropertiesName1, value)
	value = util.GetPropertyValue(device, systemPropertiesName2)
	t.Logf("name=%s, value=%s", systemPropertiesName2, value)
	value = util.GetPropertyValue(device, "non-exist-name")
	t.Logf("name=%s, value=%s", "non-exist-name", value)
}

// nolint: funlen
func TestGetDesiredDisplayNameByResourceAndConfig(t *testing.T) {
	t.Parallel()
	TestCases := []struct {
		testcasename                  string
		name                          string
		namespace                     string
		clusterName                   string
		resource                      enums.ResourceType
		displayNameIncludeNamespace   bool
		displayNameIncludeClusterName bool
		expectedResult                string
	}{
		{
			testcasename:                  "Get displayName for pod when full name is disabled.",
			name:                          "pod-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Pods,
			displayNameIncludeNamespace:   false,
			displayNameIncludeClusterName: false,
			expectedResult:                "pod-device-pod",
		},
		{
			testcasename:                  "Get displayName for pod when namespace is included",
			name:                          "pod-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Pods,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: false,
			expectedResult:                "pod-device-pod-default",
		},
		{
			testcasename:                  "Get displayName for pod when full name is enabled.",
			name:                          "pod-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Pods,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: true,
			expectedResult:                "pod-device-pod-default-cluster1",
		},
		{
			testcasename:                  "Get displayName for service when full name is disabled.",
			name:                          "service-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Services,
			displayNameIncludeNamespace:   false,
			displayNameIncludeClusterName: false,
			expectedResult:                "service-device-svc",
		},
		{
			testcasename:                  "Get displayName for service when namespace is included",
			name:                          "service-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Services,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: false,
			expectedResult:                "service-device-svc-default",
		},
		{
			testcasename:                  "Get displayName for service when full name is enabled.",
			name:                          "service-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Services,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: true,
			expectedResult:                "service-device-svc-default-cluster1",
		},
		{
			testcasename:                  "Get displayName for deployment when full name is disabled.",
			name:                          "deloyment-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Deployments,
			displayNameIncludeNamespace:   false,
			displayNameIncludeClusterName: false,
			expectedResult:                "deloyment-device-deploy",
		},
		{
			testcasename:                  "Get displayName for deloyment when namespace is included",
			name:                          "deloyment-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Deployments,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: false,
			expectedResult:                "deloyment-device-deploy-default",
		},
		{
			testcasename:                  "Get displayName for deloyment when full name is enabled.",
			name:                          "deloyment-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Deployments,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: true,
			expectedResult:                "deloyment-device-deploy-default-cluster1",
		},
		{
			testcasename:                  "Get displayName for hpa when full name is disabled.",
			name:                          "hpa-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Hpas,
			displayNameIncludeNamespace:   false,
			displayNameIncludeClusterName: false,
			expectedResult:                "hpa-device-hpa",
		},
		{
			testcasename:                  "Get displayName for hpa when namespace is included",
			name:                          "hpa-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Hpas,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: false,
			expectedResult:                "hpa-device-hpa-default",
		},
		{
			testcasename:                  "Get displayName for hpa when full name is enabled.",
			name:                          "hpa-device",
			namespace:                     "default",
			clusterName:                   "cluster1",
			resource:                      enums.Hpas,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: true,
			expectedResult:                "hpa-device-hpa-default-cluster1",
		},
		{
			testcasename:                  "Get displayName for node when full name is disabled.",
			name:                          "node-device",
			namespace:                     "",
			clusterName:                   "cluster1",
			resource:                      enums.Nodes,
			displayNameIncludeNamespace:   false,
			displayNameIncludeClusterName: false,
			expectedResult:                "node-device-node",
		},
		{
			testcasename:                  "Get displayName for node when namespace is included",
			name:                          "node-device",
			namespace:                     "",
			clusterName:                   "cluster1",
			resource:                      enums.Nodes,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: false,
			expectedResult:                "node-device-node",
		},
		{
			testcasename:                  "Get displayName for node when full name is enabled.",
			name:                          "node-device",
			namespace:                     "",
			clusterName:                   "cluster1",
			resource:                      enums.Nodes,
			displayNameIncludeNamespace:   true,
			displayNameIncludeClusterName: true,
			expectedResult:                "node-device-node-cluster1",
		},
	}

	// assertT := assert.New(t)
	conf := &config.Config{
		ClusterName:                       "cluster1",
		FullDisplayNameIncludeNamespace:   false,
		FullDisplayNameIncludeClusterName: false,
	}
	// nolint: dupl
	for _, testCase := range TestCases {
		tc := testCase
		t.Run(tc.testcasename, func(t *testing.T) {
			t.Parallel()
			result := util.GetDisplayName(tc.name, tc.namespace, tc.resource, conf)
			t.Logf(result)
			// check expected evaluation result
			// assertT.Equal(testCase.expectedResult, result, "TestCase: \"%s\" \nResult: Expected evaluate \"%s\" but got \"%s\"", testCase.testcasename, testCase.expectedResult, result)
		})
	}
}

func TestGetFullDislayName(t *testing.T) {
	t.Parallel()
	TestCases := []struct {
		testcasename   string
		device         *models.Device
		resource       enums.ResourceType
		cluster        string
		expectedResult string
	}{
		{
			testcasename:   "get full display name for pods",
			device:         getDevice(),
			resource:       enums.Pods,
			cluster:        "cluster1",
			expectedResult: "test-device-pod-default-cluster1",
		},
		{
			testcasename:   "get full display name for nodes",
			device:         getDevice(),
			resource:       enums.Nodes,
			cluster:        "cluster1",
			expectedResult: "test-device-node-cluster1",
		},
		{
			testcasename:   "get full display name for services",
			device:         getDevice(),
			resource:       enums.Services,
			cluster:        "cluster1",
			expectedResult: "test-device-svc-default-cluster1",
		},
		{
			testcasename:   "get full display name for deployment",
			device:         getDevice(),
			resource:       enums.Deployments,
			cluster:        "cluster1",
			expectedResult: "test-device-deploy-default-cluster1",
		},
		{
			testcasename:   "get full display name for hpa",
			device:         getDevice(),
			resource:       enums.Hpas,
			cluster:        "cluster1",
			expectedResult: "test-device-hpa-default-cluster1",
		},
	}
	assert2 := assert.New(t)
	// nolint: dupl
	for _, testCase := range TestCases {
		tc := testCase
		t.Run(tc.testcasename, func(t *testing.T) {
			t.Parallel()
			result := util.GetFullDisplayName(testCase.device, testCase.resource, testCase.cluster)

			// check expected evaluation result
			assert2.Equal(testCase.expectedResult, result, "TestCase: \"%s\" \nResult: Expected evaluate \"%s\" but got \"%s\"", testCase.testcasename, testCase.expectedResult, result)
		})
	}
}

func getDevice() *models.Device {
	return &models.Device{ // nolint: exhaustivestruct
		Name:        &deviceName,
		DisplayName: &deviceName,
		CustomProperties: []*models.NameAndValue{ // nolint: exhaustivestruct
			{
				Name:  &customPropertiesName1,
				Value: &customPropertiesValue1,
			}, {
				Name:  &customPropertiesName2,
				Value: &customPropertiesValue2,
			},
		},
		SystemProperties: []*models.NameAndValue{ // nolint: exhaustivestruct
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
