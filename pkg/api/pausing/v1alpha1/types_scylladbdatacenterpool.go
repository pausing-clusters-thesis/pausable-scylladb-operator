package v1alpha1

import (
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ScyllaDBDatacenterPoolControllerGVK = GroupVersion.WithKind("ScyllaDBDatacenterPool")
)

type ScyllaDBDatacenterTemplate struct {
	// spec defines the specification of ScyllaDBDatacenter instances in the pool.
	Spec scyllav1alpha1.ScyllaDBDatacenterSpec `json:"spec"`
}

type ScyllaDBDatacenterPoolSpec struct {
	// template describes the ScyllaDBDatacenter that will be created if the number of instances is lower than the desired capacity.
	Template ScyllaDBDatacenterTemplate `json:"template"`

	// proxyStorageClassName is the name of the storage class used to replace the storage class requested by created instances.
	ProxyStorageClassName string `json:"proxyStorageClassName"`

	// capacity is the number of desired ScyllaDBDatacenter instances to keep running at all times.
	Capacity int32 `json:"capacity"`

	// limit describes the maximum number of ScyllaDBDatacenter instances that should be in the pool.
	// Its purpose is to limit the oscillation of the number of instances.
	Limit int32 `json:"limit"`
}

type ScyllaDBDatacenterPoolStatus struct {
	// observedGeneration reflects the most recently observed generation of the ScyllaDBDatacenterPool.
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty"`

	// conditions reflect the latest observed conditions of the ScyllaDBDatacenterPool's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// instances reflect the number of actual instances in the pool.
	// +optional
	Instances *int32 `json:"instances,omitempty"`

	// instances reflect the number of ready instances in the pool.
	// +optional
	ReadyInstances *int32 `json:"readyInstances,omitempty"`
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

type ScyllaDBDatacenterPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec defines the specification of the ScyllaDBDatacenterPool.
	Spec ScyllaDBDatacenterPoolSpec `json:"spec,omitempty"`

	// status reflects the observed state of the ScyllaDBDatacenterPool.
	Status ScyllaDBDatacenterPoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ScyllaDBDatacenterPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScyllaDBDatacenterPool `json:"items"`
}
