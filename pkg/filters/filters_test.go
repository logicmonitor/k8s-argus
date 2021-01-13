package filters

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func init() {
	initFilterExprMap()
}

func initFilterExprMap() {
	expressionMap = make(map[string]string)
	expressionMap[constants.Pods] = "p1 =~ 'v1' || p3 =~ 'v*'"
	expressionMap[constants.Deployments] = "d1 =~ 'v1' || d4 =~ 'v4'"
	expressionMap[constants.Nodes] = "*"
	expressionMap[constants.Services] = "s1 =~ 'dev' || s1 =~ 'qa'"
	expressionMap["TestDot"] = "kubernetes.node =~ 'abc*'"
	expressionMap["TestDash"] = "node-app =~ 'TestNode'"
	expressionMap["TestSlash"] = "kubernetes/hostname =~ 'host-One'"
	expressionMap["TestAllChars"] = "kubernetes.io/pod-name =~ 'pod-test-01'"
	expressionMap["TestPodName"] = "name =~ 'pod-device'"
}

func TestEvaluate(t *testing.T) {
	t.Parallel()
	filterTestCases := []struct {
		name           string
		resource       string
		evalParams     map[string]interface{}
		expectedResult bool
	}{
		{
			name:           "Invalid resouce name",
			resource:       "XYZ",
			evalParams:     getSampleEvaluationParamsForPod1(),
			expectedResult: false,
		},
		{
			name:           "Pods, check equality",
			resource:       constants.Pods,
			evalParams:     getSampleEvaluationParamsForPod1(),
			expectedResult: true,
		},
		{
			name:           "Pods- different evaluation params",
			resource:       constants.Pods,
			evalParams:     getSampleEvaluationParamsForPod2(),
			expectedResult: true,
		},
		{
			name:           "Nodes- check all",
			resource:       constants.Nodes,
			evalParams:     getSampleEvaluationParamsForNode1(),
			expectedResult: true,
		},
		{
			name:           "Deployment- check valid expession",
			resource:       constants.Deployments,
			evalParams:     getSampleEvaluationParamsForDep1(),
			expectedResult: true,
		},
		{
			name:           "Deployment- check invalid key in expession",
			resource:       constants.Deployments,
			evalParams:     getSampleEvaluationParamsForDep2(),
			expectedResult: false,
		},
		{
			name:           "Services- label with dev value",
			resource:       constants.Services,
			evalParams:     getSampleEvaluationParamsForSvc1(),
			expectedResult: true,
		},
		{
			name:           "Services- label with qa value",
			resource:       constants.Services,
			evalParams:     getSampleEvaluationParamsForSvc2(),
			expectedResult: true,
		},
		{
			name:           "Test dots in key",
			resource:       "TestDot",
			evalParams:     getSampleEvaluationParamsForDotInKey(),
			expectedResult: true,
		},
		{
			name:           "Test dash in key",
			resource:       "TestDash",
			evalParams:     getSampleEvaluationParamsForDashInKey(),
			expectedResult: true,
		},
		{
			name:           "Test slash in key",
			resource:       "TestSlash",
			evalParams:     getSampleEvaluationParamsForSlashInKey(),
			expectedResult: true,
		},
		{
			name:           "Test all supported chars in key and value",
			resource:       "TestAllChars",
			evalParams:     getSampleEvaluationParamsForAllSuppCharsInKeyAndValue(),
			expectedResult: true,
		},
		{
			name:           "Test special chars in resource name",
			resource:       "TestPodName",
			evalParams:     getSampleEvaluationParamsForPod1(),
			expectedResult: true,
		},
	}

	assert := assert.New(t)
	// nolint: dupl
	for _, testCase := range filterTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := Eval(testCase.resource, testCase.evalParams)

			// check expected evaluation result
			assert.Equal(testCase.expectedResult, result, "TestCase: \"%s\" \nResult: Expected evaluate \"%s\" but got \"%s\"", testCase.name, testCase.expectedResult, result)
		})
	}
}

func getSampleEvaluationParamsForPod1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["p1"] = "v1"
	labels["name"] = "pod_device"
	return labels
}

func getSampleEvaluationParamsForPod2() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["p1"] = "v2"
	labels["p3"] = "v4"
	labels["name"] = "pod-device2"
	return labels
}

func getSampleEvaluationParamsForDep1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["d1"] = "v1"
	labels["name"] = "depl-device"
	return labels
}

func getSampleEvaluationParamsForDep2() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["d2"] = "v1"
	labels["d3"] = "v4"
	labels["name"] = "depl-device2"
	return labels
}

func getSampleEvaluationParamsForSvc1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["s1"] = "dev"
	labels["name"] = "svc-device"
	return labels
}

func getSampleEvaluationParamsForSvc2() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["s1"] = "qa"
	labels["s3"] = "v3"
	labels["name"] = "svc-device2"
	return labels
}

func getSampleEvaluationParamsForNode1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["testqa"] = "abc"
	labels["name"] = "node-device"
	return labels
}

func getSampleEvaluationParamsForDotInKey() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["kubernetes_node"] = "abc"
	labels["name"] = "node-device"
	return labels
}

func getSampleEvaluationParamsForDashInKey() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["node_app"] = "TestNode"
	labels["name"] = "node-device"
	return labels
}

func getSampleEvaluationParamsForSlashInKey() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["kubernetes/hostname"] = "host_One"
	labels["name"] = "node-device"
	return labels
}

func getSampleEvaluationParamsForAllSuppCharsInKeyAndValue() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["kubernetes_io/pod_name"] = "pod_test_01"
	labels["name"] = "node-device"
	return labels
}
