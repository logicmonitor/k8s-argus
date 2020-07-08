package pod

import (
	// "log"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	// "k8s.io/client-go/tools/cache"
)

func TestGetPodsMap(t *testing.T) {
	t.Parallel()
	podTestCases := []struct {
		name              string
		clientSet         kubernetes.Interface
		inputNamespace    string
		isHostNetwork     bool
		expectedPodName   string
		expectedPodsCount int
		err               string
	}{
		{
			name:              "No Pod",
			clientSet:         fake.NewSimpleClientset(),
			expectedPodsCount: 0,
		},
		{
			name: "2 Pods with HostNetwork disabled",
			clientSet: fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "node-app",
					Namespace: "default",
				},
				Spec: v1.PodSpec{
					HostNetwork: false,
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					PodIP: "10.96.90.1",
				},
			}, &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "golang-app",
					Namespace: "default",
				},
				Spec: v1.PodSpec{
					HostNetwork: false,
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					PodIP: "10.96.90.2",
				},
			}),
			expectedPodsCount: 2,
		},
		{
			name: "1 Pod with HostNetwork enabled",
			clientSet: fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "node-app",
					Namespace: "default",
				},
				Spec: v1.PodSpec{
					HostNetwork: true,
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					PodIP: "10.96.90.1",
				},
			}),
			isHostNetwork:     true,
			expectedPodName:   "node-app",
			expectedPodsCount: 1,
		},
	}

	assert := assert.New(t)
	for _, testCase := range podTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			podsMap, err := GetPodsMap(testCase.clientSet, testCase.inputNamespace)

			// check if err not nil
			if err != nil {
				assert.EqualError(err, testCase.err, "TestCase: \"%s\" \nResult: Expected error \"%s\" but got \"%s\"", testCase.name, testCase.err, err.Error())
			}

			// check expected pods count
			assert.Equal(testCase.expectedPodsCount, len(podsMap), "TestCase: \"%s\" \nResult: Expected pod count \"%d\" but got \"%d\"", testCase.name, testCase.expectedPodsCount, len(podsMap))

			// check if hostNetwork is enabled then pod name will be the IP/DNS name of the pod device
			if testCase.isHostNetwork {
				assert.Equal(testCase.expectedPodName, podsMap[testCase.expectedPodName], "TestCase: \"%s\" \nResult: Expected pod name \"%s\" but got \"%s\"", testCase.name, testCase.expectedPodName, podsMap[testCase.expectedPodName])
			}
		})
	}
}

// type PodListWatch struct {
// 	kubernetes.Interface
// }

// func (p PodListWatch) List(options metav1.ListOptions) (runtime.Object, error) {
// 	return p.Interface.CoreV1().Pods("default").List(options)
// }

// func (p PodListWatch) Watch(options metav1.ListOptions) (watch.Interface, error) {
// 	return p.Interface.CoreV1().Pods("default").Watch(options)
// }

// func TestAddFunc(t *testing.T) {
// 	clientSet := fake.NewSimpleClientset()
// 	clientSet.CoreV1().Pods("default").Create(&v1.Pod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "node-app",
// 			Namespace: "default",
// 		},
// 		Status: v1.PodStatus{
// 			Phase: v1.PodRunning,
// 			PodIP: "10.96.90.1",
// 		},
// 	})

// 	lw := PodListWatch{clientSet}
// 	watch := Watcher{}
// 	_, controller := cache.NewInformer(lw, watch.ObjType(), 2,
// 		cache.ResourceEventHandlerFuncs{
// 			AddFunc:    watch.AddFunc(),
// 			DeleteFunc: func(obj interface{}) {},
// 			UpdateFunc: func(oldObj, newObj interface{}) {},
// 		})
// 	stop := make(chan struct{})
// 	controller.Run(stop)
// 	log.Println("controller: ", controller)
// }
