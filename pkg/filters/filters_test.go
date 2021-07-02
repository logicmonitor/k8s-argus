package filters_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/enums"
)

/*func init() {
	initFilterExprMap()
}*/

/*func initFilterExprMap() {

	var str = `
filters:
  pod: []
  service: []
  node: []
  deployment:
    - 'name == "collectorset-controller"'
  hpa: []`
	filters.Init()
}*/

const nodeName = "node-resource"

func TestEvaluate(t *testing.T) {
	t.Parallel()
	filterTestCases := getEvalTestData()
	// assert2 := assert.New(t)
	// nolint: dupl
	for _, testCase := range filterTestCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// result, err := filters.Eval(testCase.resource, testCase.evalParams)
			// if err != nil {
			//	assert2.Fail("Error for test %s: %s", testCase.name, err)
			// }
			// t.Log(result)

			// check expected evaluation result
			// assert2.Equal(testCase.expectedResult, result, "TestCase: \"%s\" \nResult: Expected evaluate \"%s\" but got \"%s\"", testCase.name, testCase.expectedResult, result)
		})
	}
}

func getEvalTestData() []struct {
	name           string
	resource       enums.ResourceType
	evalParams     map[string]interface{}
	expectedResult bool
} {
	filterTestCases := []struct {
		name           string
		resource       enums.ResourceType
		evalParams     map[string]interface{}
		expectedResult bool
	}{
		{
			name:           "Invalid resource name",
			resource:       enums.Unknown,
			evalParams:     getSampleEvaluationParamsForPod1(),
			expectedResult: false,
		},
		{
			name:           "Pods, check equality",
			resource:       enums.Pods,
			evalParams:     getSampleEvaluationParamsForPod1(),
			expectedResult: true,
		},
		{
			name:           "Pods- different evaluation params",
			resource:       enums.Pods,
			evalParams:     getSampleEvaluationParamsForPod2(),
			expectedResult: true,
		},
		{
			name:           "Nodes- check all",
			resource:       enums.Nodes,
			evalParams:     getSampleEvaluationParamsForNode1(),
			expectedResult: true,
		},
		{
			name:           "Deployment- check valid expession",
			resource:       enums.Deployments,
			evalParams:     getSampleEvaluationParamsForDep1(),
			expectedResult: true,
		},
		{
			name:           "Deployment- check invalid key in expession",
			resource:       enums.Deployments,
			evalParams:     getSampleEvaluationParamsForDep2(),
			expectedResult: false,
		},
		{
			name:           "Services- label with dev value",
			resource:       enums.Services,
			evalParams:     getSampleEvaluationParamsForSvc1(),
			expectedResult: true,
		},
		{
			name:           "Services- label with qa value",
			resource:       enums.Services,
			evalParams:     getSampleEvaluationParamsForSvc2(),
			expectedResult: true,
		},
		{
			name:           "Test dots in key",
			resource:       enums.Unknown,
			evalParams:     getSampleEvaluationParamsForDotInKey(),
			expectedResult: true,
		},
		{
			name:           "Test dash in key",
			resource:       enums.Unknown,
			evalParams:     getSampleEvaluationParamsForDashInKey(),
			expectedResult: true,
		},
		{
			name:           "Test slash in key",
			resource:       enums.Unknown,
			evalParams:     getSampleEvaluationParamsForSlashInKey(),
			expectedResult: true,
		},
		{
			name:           "Test all supported chars in key and value",
			resource:       enums.Unknown,
			evalParams:     getSampleEvaluationParamsForAllSuppCharsInKeyAndValue(),
			expectedResult: true,
		},
		{
			name:           "Test special chars in resource name",
			resource:       enums.Unknown,
			evalParams:     getSampleEvaluationParamsForPod1(),
			expectedResult: true,
		},
	}
	return filterTestCases
}

func getSampleEvaluationParamsForPod1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["p1"] = "v1"
	labels["name"] = "pod_resource"

	return labels
}

func getSampleEvaluationParamsForPod2() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["p1"] = "v2"
	labels["p3"] = "v4"
	labels["name"] = "pod-resource2"

	return labels
}

func getSampleEvaluationParamsForDep1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["d1"] = "v1"
	labels["name"] = "depl-resource"

	return labels
}

func getSampleEvaluationParamsForDep2() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["d2"] = "v1"
	labels["d3"] = "v4"
	labels["name"] = "depl-resource2"

	return labels
}

func getSampleEvaluationParamsForSvc1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["s1"] = "dev"
	labels["name"] = "svc-resource"

	return labels
}

func getSampleEvaluationParamsForSvc2() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["s1"] = "qa"
	labels["s3"] = "v3"
	labels["name"] = "svc-resource2"

	return labels
}

func getSampleEvaluationParamsForNode1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["testqa"] = "abc"

	labels["name"] = nodeName

	return labels
}

func getSampleEvaluationParamsForDotInKey() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["kubernetes_node"] = "abc"
	labels["name"] = nodeName

	return labels
}

func getSampleEvaluationParamsForDashInKey() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["node_app"] = "TestNode"
	labels["name"] = nodeName

	return labels
}

func getSampleEvaluationParamsForSlashInKey() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["kubernetes/hostname"] = "host_One"
	labels["name"] = nodeName

	return labels
}

func getSampleEvaluationParamsForAllSuppCharsInKeyAndValue() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["kubernetes_io/pod_name"] = "pod_test_01"
	labels["name"] = nodeName

	return labels
}
