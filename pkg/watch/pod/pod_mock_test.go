package pod_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAddFunc(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	podObject := &corev1.Pod{ // nolint: exhaustivestruct
		ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
			Name:      "dummy-pod-name",
			Namespace: "dummy-namespace",
			Labels: map[string]string{
				"abc": "xyz",
			},
		},
		Spec: corev1.PodSpec{ // nolint: exhaustivestruct
			HostNetwork: true,
		},
		Status: corev1.PodStatus{ // nolint: exhaustivestruct
			Phase: corev1.PodRunning,
			PodIP: "127.0.0.1",
		},
	}
	t.Log(podObject)

	// mDevMgr := mocks.NewMockDeviceManager(ctrl)
	// pWatcher := Watcher{
	// }
	// dummyDeviceOption := func(d *models.Device) {}
	// mDevMgr.EXPECT().Name(gomock.Eq("dummy-pod-name")).Return(dummyDeviceOption)
	// mDevMgr.EXPECT().ResourceLabels(gomock.Eq(podObject.Labels)).Return(dummyDeviceOption)
	// mDevMgr.EXPECT().DisplayName(gomock.Eq("dummy-desired-disp-name")).Return(dummyDeviceOption)
	// mDevMgr.EXPECT().SystemCategory("KubernetesPod").Return(dummyDeviceOption)
	// mDevMgr.EXPECT().Auto(gomock.Any(), gomock.Any()).Return(dummyDeviceOption).Times(5)
	// mDevMgr.EXPECT().System(gomock.Any(), gomock.Any()).Return(dummyDeviceOption).Times(1)
	// mDevMgr.EXPECT().Custom(gomock.Any(), gomock.Any()).Return(dummyDeviceOption).Times(3)
	// mDevMgr.EXPECT().GetDisplayName(gomock.Eq("dummy-pod-name"), gomock.Eq("dummy-namespace"), gomock.Eq("pods")).Return("dummy-desired-disp-name").Times(4)
	//
	// mDevMgr.EXPECT().Add(gomock.Any(), gomock.Eq("pods"), gomock.Eq(podObject.Labels), gomock.Any()).DoAndReturn(
	//	func(a *lmctx.LMContext, b string, c map[string]string, d ...types.DeviceOption) (interface{}, string, map[string]string, []types.DeviceOption) {
	//		t.Logf("Add called with Arguments: %v %v %v %v", a, b, c, d)
	//		if a.Extract(constants.IsPingDevice).(bool) == false {
	//			ctrl.T.Errorf("Context param %v is false, Expected is true")
	//		}
	//
	// return a, b, c, d
	//	},
	// )

	// f := pWatcher.AddFuncOptions()
	// f(podObject)
}

func TestAddFuncPodWithoutIP(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	podObject := &corev1.Pod{ // nolint: exhaustivestruct
		ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
			Name:      "dummy-pod-name",
			Namespace: "dummy-namespace",
			Labels: map[string]string{
				"abc": "xyz",
			},
		},
		Spec: corev1.PodSpec{ // nolint: exhaustivestruct
			HostNetwork: true,
		},
		Status: corev1.PodStatus{ // nolint: exhaustivestruct
			Phase: corev1.PodPending,
		},
	}
	t.Logf("%v", podObject)

	// mDevMgr := mocks.NewMockDeviceManager(ctrl)
	// pWatcher := Watcher{
	// }
	// mDevMgr.EXPECT().GetDesiredDisplayName(gomock.Eq("dummy-pod-name"), gomock.Eq("dummy-namespace"), gomock.Eq("pods")).Return("dummy-desired-disp-name")

	// f := pWatcher.AddFunc()
	// f(podObject)
}
