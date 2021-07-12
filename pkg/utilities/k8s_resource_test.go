package utilities_test

import (
	"testing"

	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestSelfLinkCoreGlobalResource tests core api resource which is not namespaced
func TestSelfLinkCoreGlobalResource(t *testing.T) {
	t.Parallel()
	expectedResult := "/api/v1/nodes/node1"
	objectMeta := &metav1.ObjectMeta{ // nolint: exhaustivestruct
		Name:      "node1",
		Namespace: "",
	}
	result := util.SelfLink(false, "v1", "nodes", meta.AsPartialObjectMetadata(objectMeta))
	if result != expectedResult {
		t.Errorf("Result selflink: %s does not match with expected %s", result, expectedResult)
	}
}

// TestSelfLinkCoreNamespacedResource tests core api resource which is namespaced
func TestSelfLinkCoreNamespacedResource(t *testing.T) {
	t.Parallel()
	expectedResult := "/api/v1/namespaces/ns1/pods/pod1"
	objectMeta := &metav1.ObjectMeta{ // nolint: exhaustivestruct
		Name:      "pod1",
		Namespace: "ns1",
	}
	result := util.SelfLink(true, "v1", "pods", meta.AsPartialObjectMetadata(objectMeta))
	if result != expectedResult {
		t.Errorf("Result selflink: %s does not match with expected %s", result, expectedResult)
	}
}

// TestSelfLinkCoreGlobalResource tests core api resource which is not namespaced
func TestSelfLinkAPIGroupNamespacedResource(t *testing.T) {
	t.Parallel()
	expectedResult := "/apis/apps/v1/namespaces/ns1/deployments/deploy1"
	objectMeta := &metav1.ObjectMeta{ // nolint: exhaustivestruct
		Name:      "deploy1",
		Namespace: "ns1",
	}
	result := util.SelfLink(true, "apps/v1", "deployments", meta.AsPartialObjectMetadata(objectMeta))
	if result != expectedResult {
		t.Errorf("Result selflink: %s does not match with expected %s", result, expectedResult)
	}
}
