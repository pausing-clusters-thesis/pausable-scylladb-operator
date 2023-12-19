package v1alpha1

import (
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ExposeOptions struct {
	// cql specifies expose options for CQL SSL backend.
	// +optional
	CQL *scyllav1alpha1.CQLExposeOptions `json:"cql,omitempty"`
}

type PausableScyllaDBDatacenterSpec struct {
	// scyllaDBDatacenterPoolName specifies the name of the ScyllaDBDatacenterPool provisioning the underlying ScyllaDBDatacenter.
	ScyllaDBDatacenterPoolName string `json:"scyllaDBDatacenterPoolName"`

	// paused specifies whether the resource should be paused.
	// When paused, the underlying ScyllaDBDatacenter is released.
	// +kubebuilder:default:=false
	// +optional
	Paused *bool `json:"paused,omitempty"`

	// exposeOptions specifies parameters related to exposing PausableScyllaDBDatacenter backends.
	// +optional
	ExposeOptions *ExposeOptions `json:"exposeOptions,omitempty"`
}

type PausableScyllaDBDatacenterStatus struct {
	// observedGeneration reflects the most recently observed generation of the PausableScyllaDBDatacenter.
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty"`

	// conditions reflect the latest observed conditions of the PausableScyllaDBDatacenter's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="AVAILABLE",type=string,JSONPath=".status.conditions[?(@.type=='Available')].status"
// +kubebuilder:printcolumn:name="PROGRESSING",type=string,JSONPath=".status.conditions[?(@.type=='Progressing')].status"
// +kubebuilder:printcolumn:name="DEGRADED",type=string,JSONPath=".status.conditions[?(@.type=='Degraded')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

type PausableScyllaDBDatacenter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec defines the specification of the PausableScyllaDBDatacenter.
	Spec PausableScyllaDBDatacenterSpec `json:"spec,omitempty"`

	// status reflects the observed state of the PausableScyllaDBDatacenter.
	Status PausableScyllaDBDatacenterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PausableScyllaDBDatacenterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PausableScyllaDBDatacenter `json:"items"`
}
