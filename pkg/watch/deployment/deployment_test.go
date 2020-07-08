package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetDeploymentsMap(t *testing.T) {
	t.Parallel()
	podTestCases := []struct {
		name                     string
		clientSet                kubernetes.Interface
		inputNamespace           string
		expectedDeploymentsCount int
		err                      string
	}{
		{
			name:                     "No Deployment",
			clientSet:                fake.NewSimpleClientset(),
			expectedDeploymentsCount: 0,
		},
		{
			name: "2 Deployments",
			clientSet: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "node-app",
					Namespace: "default",
					UID:       "0b760df4-f746-4034-86bd-30e10fae5521",
				},
			}, &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "golang-app",
					Namespace: "default",
					UID:       "c48021ae-b68b-4ba1-befe-5e1b659212d3",
				},
			}),
			expectedDeploymentsCount: 2,
		},
	}
	assert := assert.New(t)
	for _, testCase := range podTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			deploymentsMap, err := GetDeploymentsMap(testCase.clientSet, testCase.inputNamespace)

			// check if err not nil
			if err != nil {
				assert.EqualError(err, testCase.err, "TestCase: \"%s\" \nResult: Expected error \"%s\" but got \"%s\"", testCase.name, testCase.err, err.Error())
			}

			// check expected deployments count
			assert.Equal(testCase.expectedDeploymentsCount, len(deploymentsMap), "TestCase: \"%s\" \nResult: Expected Deployment count \"%d\" but got \"%d\"", testCase.name, testCase.expectedDeploymentsCount, len(deploymentsMap))
		})
	}
}
