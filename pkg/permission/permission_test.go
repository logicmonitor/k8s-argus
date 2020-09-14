package permission

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestHasDeploymentPermissions(t *testing.T) {

	deploymentTestCases := []struct {
		name          string
		clientSet     kubernetes.Interface
		hasPermission bool
	}{
		{
			name:          "No Deployment",
			clientSet:     fake.NewSimpleClientset(),
			hasPermission: true,
		},
		{
			name: "2 Deployments",
			clientSet: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "node-app",
					Namespace: "default",
				},
			}, &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "golang-app",
					Namespace: "default",
				},
			}),
			hasPermission: true,
		},
	}

	assert := assert.New(t)
	for _, testCase := range deploymentTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			Init(testCase.clientSet)
			hasPermission := HasDeploymentPermissions()
			assert.Equal(testCase.hasPermission, hasPermission, "TestCase: \"%s\" \nResult: Expected \"%t\" but got \"%t\"", testCase.name, testCase.hasPermission, hasPermission)
		})
	}
}
