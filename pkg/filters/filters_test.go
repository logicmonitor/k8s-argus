package filters

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func init() {
	initFilterExprMap()
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
			name:           "Nodes- check back slash in key",
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
			name:           "Services- label with qa value ",
			resource:       constants.Services,
			evalParams:     getSampleEvaluationParamsForSvc2(),
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

func initFilterExprMap() {
	expressionMap = make(map[string]string)
	expressionMap[constants.Pods] = "p1 =~ 'v1' || p3 =~ 'v*'"
	expressionMap[constants.Deployments] = "d1 =~ 'v1' || d4 =~ 'v4'"
	expressionMap[constants.Nodes] = "test\\/qa == 'abc'"
	expressionMap[constants.Services] = "s1 =~ 'dev' || s1 =~ 'qa'"
}

func getSampleEvaluationParamsForPod1() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["p1"] = "v1"
	labels["name"] = "pod-device"
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
	labels["test/qa"] = "abc"
	labels["name"] = "node-device"
	return labels
}
