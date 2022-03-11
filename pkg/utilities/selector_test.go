package utilities_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCoalesceMap(t *testing.T) {
	input := map[string]string{
		"abc": "xyz",
		"mno": "pqr",
	}
	assertObj := assert.New(t)
	result := utilities.GenerateSelectorExpression(input)
	assertObj.Equal("abc=xyz, mno=pqr", result)
}

func TestCoalesceMapNil(t *testing.T) {
	assertObj := assert.New(t)
	result := utilities.GenerateSelectorExpression(nil)
	assertObj.Equal("null", result)
}

// nolint: dupl
func TestGenerateSelectorExpression(t *testing.T) {
	testCases := []struct {
		name           string
		input          metav1.LabelSelector
		expectedOutput string
	}{
		{
			name:           "Nil Input test",
			input:          metav1.LabelSelector{},
			expectedOutput: "null",
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
			expectedOutput: "testKeyIn in (Val1, Val2)",
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
			expectedOutput: "notInTestKey notin (FirstNotIn, SecondNotIn)",
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
			expectedOutput: "existsKey",
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
			expectedOutput: "!notExistsKey",
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
			expectedOutput: "existsKey, !notExistsKey",
		},
	}
	assertObj := assert.New(t)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := utilities.GenerateSelectorExpression(testCase.input)
			assertObj.Equal(testCase.expectedOutput, result, "TestCase: \"%s\" \nResult: Expected key \"%s\" but got \"%s\"", testCase.name, testCase.expectedOutput, result)
		})
	}
}
