package namespace_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/watch/namespace"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

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
			clientSet: fake.NewSimpleClientset(&corev1.Namespace{ // nolint: exhaustivestruct
				ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
					Name: "default",
				},
			}, &corev1.Namespace{ // nolint: exhaustivestruct
				ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
					Name: "kube-system",
				},
			}),
			expectedNamespaceList:      []string{"default", "kube-system"},
			expectedNamespaceListCount: 2,
		},
	}

	assertObj := assert.New(t)
	for _, testCase := range namespaceListTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			namespaceListOutput := namespace.GetNamespaceList(testCase.lctx, testCase.clientSet)

			// check expected namespace list
			assertObj.Equal(testCase.expectedNamespaceList, namespaceListOutput, "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedNamespaceList, namespaceListOutput)

			// check expected namespace list count
			assertObj.Equal(testCase.expectedNamespaceListCount, len(namespaceListOutput), "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedNamespaceListCount, len(namespaceListOutput))
		})
	}
}
