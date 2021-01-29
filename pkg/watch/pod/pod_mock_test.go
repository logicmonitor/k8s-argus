package pod

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/mocks"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAddFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	podObject := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dummy-pod-name",
			Namespace: "dummy-namespace",
			Labels: map[string]string{
				"abc": "xyz",
			},
		},
		Spec: v1.PodSpec{
			HostNetwork: true,
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
			PodIP: "127.0.0.1",
		},
	}

	mDevMgr := mocks.NewMockDeviceManager(ctrl)
	pWatcher := Watcher{
		mDevMgr,
		nil,
	}
	dummyDeviceOption := func(d *models.Device) {}
	mDevMgr.EXPECT().Name(gomock.Eq("dummy-pod-name")).Return(dummyDeviceOption)
	mDevMgr.EXPECT().ResourceLabels(gomock.Eq(podObject.Labels)).Return(dummyDeviceOption)
	mDevMgr.EXPECT().DisplayName(gomock.Eq("dummy-desired-disp-name")).Return(dummyDeviceOption)
	mDevMgr.EXPECT().SystemCategories("KubernetesPod").Return(dummyDeviceOption)
	mDevMgr.EXPECT().Auto(gomock.Any(), gomock.Any()).Return(dummyDeviceOption).Times(5)
	mDevMgr.EXPECT().System(gomock.Any(), gomock.Any()).Return(dummyDeviceOption).Times(1)
	mDevMgr.EXPECT().Custom(gomock.Any(), gomock.Any()).Return(dummyDeviceOption).Times(3)
	mDevMgr.EXPECT().GetDesiredDisplayName(gomock.Eq("dummy-pod-name"), gomock.Eq("dummy-namespace"), gomock.Eq("pods")).Return("dummy-desired-disp-name").Times(4)

	mDevMgr.EXPECT().Add(gomock.Any(), gomock.Eq("pods"), gomock.Eq(podObject.Labels), gomock.Any()).DoAndReturn(
		func(a *lmctx.LMContext, b string, c map[string]string, d ...types.DeviceOption) (interface{}, string, map[string]string, []types.DeviceOption) {
			t.Logf("Add called with Arguments: %v %v %v %v", a, b, c, d)
			if a.Extract(constants.IsPingDevice).(bool) == false {
				ctrl.T.Errorf("Context param %v is false, Expected is true")
			}
			return a, b, c, d
		},
	)

	f := pWatcher.AddFunc()
	f(podObject)

}
func TestAddFuncPodWithoutIP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	podObject := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dummy-pod-name",
			Namespace: "dummy-namespace",
			Labels: map[string]string{
				"abc": "xyz",
			},
		},
		Spec: v1.PodSpec{
			HostNetwork: true,
		},
		Status: v1.PodStatus{
			Phase: v1.PodPending,
		},
	}

	mDevMgr := mocks.NewMockDeviceManager(ctrl)
	pWatcher := Watcher{
		mDevMgr,
		nil,
	}
	mDevMgr.EXPECT().GetDesiredDisplayName(gomock.Eq("dummy-pod-name"), gomock.Eq("dummy-namespace"), gomock.Eq("pods")).Return("dummy-desired-disp-name")

	f := pWatcher.AddFunc()
	f(podObject)

}
