package deployment

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetDeploymentsMap(t *testing.T) {
	t.Parallel()
	deploymentTestCases := []struct {
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
	// nolint: dupl
	for _, testCase := range deploymentTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": "get_deploys_map_test"}))
			deploymentsMap, err := GetDeploymentsMap(lctx, testCase.clientSet, testCase.inputNamespace)

			// check if err not nil
			if err != nil {
				assert.EqualError(err, testCase.err, "TestCase: \"%s\" \nResult: Expected error \"%s\" but got \"%s\"", testCase.name, testCase.err, err.Error())
			}

			// check expected deployments count
			assert.Equal(testCase.expectedDeploymentsCount, len(deploymentsMap), "TestCase: \"%s\" \nResult: Expected Deployment count \"%d\" but got \"%d\"", testCase.name, testCase.expectedDeploymentsCount, len(deploymentsMap))
		})
	}
}

func TestGetHelmChartDetailsFromDeployments(t *testing.T) {
	t.Parallel()
	argus := "argus"
	argusHelmChartKey := "argus.helm-chart"
	argusHelmRevisionKey := "argus.helm-revision"
	argusHelmChartValue := "argus-0.14.0"
	argusHelmRevisionValue := "2"

	collectorsetController := "collectorset-controller"
	cscHelmChartKey := "collectorset-controller.helm-chart"
	cscHelmRevisionKey := "collectorset-controller.helm-revision"
	cscHelmChartValue := "collectorset-controller-0.9.0"
	cscHelmRevisionValue := "1"
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"res": "TestGetHelmChartDetailsFromDeployments"}))
	helmChartTestCases := []struct {
		name                    string
		lctx                    *lmctx.LMContext
		clientSet               kubernetes.Interface
		customProperties        map[string]string
		expectedProperties      map[string]string
		expectedPropertiesCount int
	}{
		{
			name: "No deployment with LabelSelector 'chart in (\"argus\", \"collectorset-controller\")",
			lctx: lctx,
			clientSet: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "sample-golang-app",
					Labels: map[string]string{},
				},
			}),
			customProperties:        make(map[string]string),
			expectedProperties:      make(map[string]string),
			expectedPropertiesCount: 0,
		},
		{
			name: "2 deployments with LabelSelector 'chart in (\"argus\", \"collectorset-controller\")",
			lctx: lctx,
			clientSet: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: argus,
					Labels: map[string]string{
						"chart":         argus,
						"helm-chart":    argusHelmChartValue,
						"helm-revision": argusHelmRevisionValue,
					},
					Namespace: "default",
				},
			}, &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "collectorset-controller",
					Labels: map[string]string{
						"chart":         collectorsetController,
						"helm-chart":    cscHelmChartValue,
						"helm-revision": cscHelmRevisionValue,
					},
					Namespace: "default",
				},
			},
				&v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
				}, &v1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "kube-system",
					},
				}),
			customProperties: make(map[string]string),
			expectedProperties: map[string]string{
				argusHelmChartKey:    argusHelmChartValue,
				argusHelmRevisionKey: argusHelmRevisionValue,
				cscHelmChartKey:      cscHelmChartValue,
				cscHelmRevisionKey:   cscHelmRevisionValue,
			},
			expectedPropertiesCount: 4,
		},
	}

	assert := assert.New(t)
	for _, testCase := range helmChartTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			customPropertiesOutput := GetHelmChartDetailsFromDeployments(testCase.lctx, testCase.customProperties, testCase.clientSet)

			// check expected customProperties list
			assert.Equal(testCase.expectedProperties, customPropertiesOutput, "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedProperties, customPropertiesOutput)

			// check expected customProperties list count
			assert.Equal(testCase.expectedPropertiesCount, len(customPropertiesOutput), "TestCase: \"%s\" \nResult: Expected output \"%v\" but got \"%v\"", testCase.name, testCase.expectedPropertiesCount, len(customPropertiesOutput))
		})
	}

}
