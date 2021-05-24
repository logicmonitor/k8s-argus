package permission_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestHasDeploymentPermissions(t *testing.T) {
	t.Parallel()
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
			clientSet: fake.NewSimpleClientset(&appsv1.Deployment{ // nolint: exhaustivestruct
				ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
					Name:      "node-app",
					Namespace: "default",
				},
			}, &appsv1.Deployment{ // nolint: exhaustivestruct
				ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
					Name:      "golang-app",
					Namespace: "default",
				},
			}),
			hasPermission: true,
		},
	}

	assertObj := assert.New(t)
	for _, testCase := range deploymentTestCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			permission.Init(tc.clientSet)
			hasPermission := permission.HasDeploymentPermissions()
			assertObj.Equal(tc.hasPermission, hasPermission, "TestCase: \"%s\" \nResult: Expected \"%t\" but got \"%t\"", tc.name, tc.hasPermission, hasPermission)
		})
	}
}
