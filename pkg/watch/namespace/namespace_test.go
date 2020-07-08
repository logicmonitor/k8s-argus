package namespace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReversedDeviceGroups(t *testing.T) {
	t.Parallel()
	deviceGroupTestCases := []struct {
		name                      string
		input                     map[string]int32
		expectedOutput            map[int32]string
		expectedDeviceGroupsCount int
	}{
		{
			name:                      "No key-value pair in map",
			input:                     map[string]int32{},
			expectedOutput:            map[int32]string{},
			expectedDeviceGroupsCount: 0,
		},
		{
			name:                      "2 key-value pair in map",
			input:                     map[string]int32{"foo": 1, "bar": 2},
			expectedOutput:            map[int32]string{1: "foo", 2: "bar"},
			expectedDeviceGroupsCount: 2,
		},
	}

	assert := assert.New(t)
	for _, testCase := range deviceGroupTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := getReversedDeviceGroups(testCase.input)

			// check expected reverse device groups
			assert.Equal(testCase.expectedOutput, output, "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedOutput, output)

			// check expected device groups count
			assert.Equal(testCase.expectedDeviceGroupsCount, len(output), "TestCase: \"%s\" \nResult: Expected map len \"%d\" but got \"%d\"", testCase.name, testCase.expectedDeviceGroupsCount, len(output))
		})
	}
}
