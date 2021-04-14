package service

import (
	"testing"

	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetServicesMap(t *testing.T) {
	t.Parallel()
	serviceTestCases := []struct {
		name                  string
		clientSet             kubernetes.Interface
		inputNamespace        string
		expectedServicesCount int
		err                   string
	}{
		{
			name:                  "No Service",
			clientSet:             fake.NewSimpleClientset(),
			expectedServicesCount: 0,
		},
		{
			name: "2 Services",
			clientSet: fake.NewSimpleClientset(&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "node-svc",
					Namespace: "default",
				},
				Spec:   v1.ServiceSpec{},
				Status: v1.ServiceStatus{},
			}, &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "golang-svc",
					Namespace: "default",
				},
				Spec:   v1.ServiceSpec{},
				Status: v1.ServiceStatus{},
			}),
			expectedServicesCount: 2,
		},
	}
	assert := assert.New(t)
	// nolint: dupl
	for _, testCase := range serviceTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": "get_svc_map_test"}))
			servicesMap, err := GetServicesMap(lctx, testCase.clientSet, testCase.inputNamespace, "cluster1")

			// check if err not nil
			if err != nil {
				assert.EqualError(err, testCase.err, "TestCase: \"%s\" \nResult: Expected error \"%s\" but got \"%s\"", testCase.name, testCase.err, err.Error())
			}

			// check expected services count
			assert.Equal(testCase.expectedServicesCount, len(servicesMap), "TestCase: \"%s\" \nResult: Expected Service count \"%d\" but got \"%d\"", testCase.name, testCase.expectedServicesCount, len(servicesMap))
		})
	}
}
