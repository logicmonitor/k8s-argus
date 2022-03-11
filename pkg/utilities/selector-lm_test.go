package utilities_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCoalesceMapAppliesTo(t *testing.T) {
	input := map[string]string{
		"abc": "xyz",
		"mno": "pqr",
	}
	assertObj := assert.New(t)
	result := utilities.GenerateSelectorAppliesTo(input)
	assertObj.Equal("kubernetes.label.abc == \"xyz\" && kubernetes.label.mno == \"pqr\"", result)
}

// nolint: dupl
func TestGenerateSelectorAppliesTo(t *testing.T) {
	testCases := []struct {
		name           string
		input          metav1.LabelSelector
		expectedOutput string
	}{
		{
			name:           "Nil Input test",
			input:          metav1.LabelSelector{},
			expectedOutput: "false()",
		},
		{
			name: "In Operator test",
			input: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "testKeyIn",
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{"Val1", "Val2"},
					},
				},
			},
			expectedOutput: "kubernetes.label.testKeyIn =~ \"(?=^Val1$)|(?=^Val2$)\"",
		},
		{
			name: "NotIn Operator test",
			input: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "notInTestKey",
						Operator: metav1.LabelSelectorOpNotIn,
						Values:   []string{"FirstNotIn", "SecondNotIn"},
					},
				},
			},
			expectedOutput: "kubernetes.label.notInTestKey !~ \"(?=^FirstNotIn$)|(?=^SecondNotIn$)\"",
		},
		{
			name: "Exists Operator test",
			input: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "existsKey",
						Operator: metav1.LabelSelectorOpExists,
					},
				},
			},
			expectedOutput: "exists(\"kubernetes.label.existsKey\")",
		},
		{
			name: "DoesNotExists Operator test",
			input: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "notExistsKey",
						Operator: metav1.LabelSelectorOpDoesNotExist,
					},
				},
			},
			expectedOutput: "!exists(\"kubernetes.label.notExistsKey\")",
		},
		{
			name: "And while coalesce test",
			input: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "existsKey",
						Operator: metav1.LabelSelectorOpExists,
					},
					{
						Key:      "notExistsKey",
						Operator: metav1.LabelSelectorOpDoesNotExist,
					},
				},
			},
			expectedOutput: "exists(\"kubernetes.label.existsKey\") && !exists(\"kubernetes.label.notExistsKey\")",
		},
	}
	assertObj := assert.New(t)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := utilities.GenerateSelectorAppliesTo(testCase.input)
			assertObj.Equal(testCase.expectedOutput, result, "TestCase: \"%s\" \nResult: Expected key \"%s\" but got \"%s\"", testCase.name, testCase.expectedOutput, result)
		})
	}
}
