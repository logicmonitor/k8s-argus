package filters

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func generateFilterExprMap() {
	expressionMap = make(map[string]string)
	expressionMap[constants.Pods] = "l1 =~ 'v1'"
	expressionMap[constants.Deployments] = ""
	expressionMap[constants.Nodes] = "*"
	expressionMap[constants.Services] = "name =~ 'test*'"
}

func init() {
	generateFilterExprMap()
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
			name:           "Pods expression, check equality",
			resource:       constants.Pods,
			evalParams:     getSampleEvaluationParams(),
			expectedResult: true,
		},
		{
			name:           "Nodes expression- check filter all",
			resource:       constants.Nodes,
			evalParams:     getSampleEvaluationParams(),
			expectedResult: true,
		},
		{
			name:           "Deployment expression- check empty",
			resource:       constants.Deployments,
			evalParams:     getSampleEvaluationParams(),
			expectedResult: false,
		},
		{
			name:           "Services expression- check blob expression ",
			resource:       constants.Services,
			evalParams:     getSampleEvaluationParams(),
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

func getSampleEvaluationParams() map[string]interface{} {
	labels := make(map[string]interface{})
	labels["l1"] = "v1"
	labels["l2"] = "v2"
	labels["name"] = "test-device"
	return labels
}
