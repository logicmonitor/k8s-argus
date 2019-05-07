package v1alpha1

import (
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/distributor"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CollectorSetState is the CollectorSet controller's state string.
type CollectorSetState string

const (
	// CollectorSetStateCreated is a status string.
	CollectorSetStateCreated CollectorSetState = "Created"
	// CollectorSetStateRegistered is a status string.
	CollectorSetStateRegistered CollectorSetState = "Registered"
	// CollectorSetResourcePlural is the plural for the CRD.
	CollectorSetResourcePlural = "collectorsets"
)

// CollectorSet represents the collectorset in Kubernetes.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CollectorSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CollectorSetSpec   `json:"spec"`
	Status            CollectorSetStatus `json:"status,omitempty"`
}

// CollectorSetSpec represents the collectorset controller's spec.
type CollectorSetSpec struct {
	Replicas          *int32              `json:"replicas"`
	Size              string              `json:"size,omitempty"`
	ClusterName       string              `json:"clusterName"`
	GroupID           int32               `json:"groupID,omitempty"`           //default value is 0, it means no group is offered
	EscalationChainID int32               `json:"escalationChainID,omitempty"` //default value is 0, it means disable notification
	CollectorVersion  int32               `json:"collectorVersion,omitempty"`  //default value is 0, it means the latest GD version
	UseEA             bool                `json:"useEA,omitempty"`             //default value is false, it means the latest GD version
	Policy            *CollectorSetPolicy `json:"policy"`
}

// CollectorSetStatus is the CollectorSet controller's status.
type CollectorSetStatus struct {
	State CollectorSetState `json:"state,omitempty"`
	IDs   []int32           `json:"ids,omitempty"`
}

// CollectorSetPolicy is the CollectorSet controller's status.
type CollectorSetPolicy struct {
	Orchestrator        string            `json:"orchestrator,omitempty"`
	DistibutionStrategy *distributor.Type `json:"distributionStrategy"`
}

// CollectorSetList represents a list of collectorsets.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CollectorSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CollectorSet `json:"items"`
}
