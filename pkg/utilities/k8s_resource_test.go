package utilities

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestSelfLinkCoreGlobalResource tests core api resource which is not namespaced
func TestSelfLinkCoreGlobalResource(t *testing.T) {
	expectedResult := "/api/v1/nodes/node1"
	objectMeta := metav1.ObjectMeta{
		Name:      "node1",
		Namespace: "",
	}
	result := SelfLink(false, "v1", "nodes", objectMeta)
	if result != expectedResult {
		t.Errorf("Result selflink: %s does not match with expected %s", result, expectedResult)
	}
}

// TestSelfLinkCoreNamespacedResource tests core api resource which is namespaced
func TestSelfLinkCoreNamespacedResource(t *testing.T) {
	expectedResult := "/api/v1/namespaces/ns1/pods/pod1"
	objectMeta := metav1.ObjectMeta{
		Name:      "pod1",
		Namespace: "ns1",
	}
	result := SelfLink(true, "v1", "pods", objectMeta)
	if result != expectedResult {
		t.Errorf("Result selflink: %s does not match with expected %s", result, expectedResult)
	}
}

// TestSelfLinkCoreGlobalResource tests core api resource which is not namespaced
func TestSelfLinkAPIGroupNamespacedResource(t *testing.T) {
	expectedResult := "/apis/apps/v1/namespaces/ns1/deployments/deploy1"
	objectMeta := metav1.ObjectMeta{
		Name:      "deploy1",
		Namespace: "ns1",
	}
	result := SelfLink(true, "apps/v1", "deployments", objectMeta)
	if result != expectedResult {
		t.Errorf("Result selflink: %s does not match with expected %s", result, expectedResult)
	}
}
