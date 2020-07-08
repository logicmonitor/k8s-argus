package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLabelByPrefix(t *testing.T) {

	testCases := []struct {
		name          string
		inputPrefix   string
		inputLabels   map[string]string
		expectedKey   string
		expectedValue string
	}{
		{
			name:          "Empty input labels",
			inputPrefix:   "node-role.kubernetes.io/",
			inputLabels:   map[string]string{},
			expectedKey:   "",
			expectedValue: "",
		},
		{
			name:        "Passing input labels",
			inputPrefix: "node-role.kubernetes.io/",
			inputLabels: map[string]string{
				"kubernetes.io/hostname":         "master-node",
				"kubernetes.io/os":               "linux",
				"node-role.kubernetes.io/master": "master-node",
			},
			expectedKey:   "node-role.kubernetes.io/master",
			expectedValue: "master-node",
		},
		{
			name:        "Passing input labels but empty value",
			inputPrefix: "node-role.kubernetes.io/",
			inputLabels: map[string]string{
				"kubernetes.io/hostname":         "master-node",
				"kubernetes.io/os":               "linux",
				"node-role.kubernetes.io/master": "",
			},
			expectedKey:   "node-role.kubernetes.io/master",
			expectedValue: "",
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			key, value := GetLabelByPrefix(testCase.inputPrefix, testCase.inputLabels)

			// check expected key
			assert.Equal(testCase.expectedKey, key, "TestCase: \"%s\" \nResult: Expected key \"%s\" but got \"%s\"", testCase.name, testCase.expectedKey, key)

			// check expected value
			assert.Equal(testCase.expectedValue, value, "TestCase: \"%s\" \nResult: Expected value \"%s\" but got \"%s\"", testCase.name, testCase.expectedValue, value)
		})
	}
}
