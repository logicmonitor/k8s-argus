package namespace

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetReversedDeviceGroups(t *testing.T) {
	t.Parallel()
	namespaceTestCases := []struct {
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
	for _, testCase := range namespaceTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := getReversedDeviceGroups(testCase.input)

			// check expected reverse device groups
			assert.Equal(testCase.expectedOutput, output, "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedOutput, output)

			// check expected device groups count
			assert.Equal(testCase.expectedDeviceGroupsCount, len(output), "TestCase: \"%s\" \nResult: Expected map len \"%d\" but got \"%d\"", testCase.name, testCase.expectedDeviceGroupsCount, len(output))
		})
	}
}

func TestGetNamespaceList(t *testing.T) {
	t.Parallel()
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"res": "TestGetNamespaceList"}))
	namespaceListTestCases := []struct {
		name                       string
		lctx                       *lmctx.LMContext
		clientSet                  kubernetes.Interface
		expectedNamespaceList      []string
		expectedNamespaceListCount int
	}{
		{
			name:                       "No namespaces",
			lctx:                       lctx,
			clientSet:                  fake.NewSimpleClientset(),
			expectedNamespaceList:      []string{},
			expectedNamespaceListCount: 0,
		},
		{
			name: "2 namespaces",
			lctx: lctx,
			clientSet: fake.NewSimpleClientset(&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}, &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "kube-system",
				},
			}),
			expectedNamespaceList:      []string{"default", "kube-system"},
			expectedNamespaceListCount: 2,
		},
	}

	assert := assert.New(t)
	for _, testCase := range namespaceListTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			namespaceListOutput := GetNamespaceList(testCase.lctx, testCase.clientSet)

			// check expected namespace list
			assert.Equal(testCase.expectedNamespaceList, namespaceListOutput, "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedNamespaceList, namespaceListOutput)

			// check expected namespace list count
			assert.Equal(testCase.expectedNamespaceListCount, len(namespaceListOutput), "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedNamespaceListCount, len(namespaceListOutput))
		})
	}
}
