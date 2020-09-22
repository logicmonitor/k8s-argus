package cronjob

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUpdatedHistoryValue(t *testing.T) {
	t.Parallel()
	historyValueTestCases := []struct {
		name               string
		inputExistingValue string
		inputNewValue      string
		expectedOutput     string
	}{
		{
			name:               "Empty existing value",
			inputExistingValue: "",
			inputNewValue:      "argus-0.14.0",
			expectedOutput:     "argus-0.14.0",
		},
		{
			name:               "Existing value with new value",
			inputExistingValue: "argus-0.13.0",
			inputNewValue:      "argus-0.14.0",
			expectedOutput:     "argus-0.13.0, argus-0.14.0",
		},
		{
			name:               "Existing value with 10 comma separated value",
			inputExistingValue: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10",
			inputNewValue:      "11",
			expectedOutput:     "2, 3, 4, 5, 6, 7, 8, 9, 10, 11",
		},
		{
			name:               "Existing value with more than 10 comma separated value",
			inputExistingValue: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13",
			inputNewValue:      "14",
			expectedOutput:     "5, 6, 7, 8, 9, 10, 11, 12, 13, 14",
		},
	}

	assert := assert.New(t)
	// nolint: dupl
	for _, testCase := range historyValueTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := getUpdatedHistoryValue(testCase.inputExistingValue, testCase.inputNewValue)
			assert.Equal(testCase.expectedOutput, output, "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedOutput, output)
		})
	}
}
