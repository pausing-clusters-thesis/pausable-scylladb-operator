package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ScyllaDBDatacenterClaimControllerGVK = GroupVersion.WithKind("ScyllaDBDatacenterClaim")
)

type ScyllaDBDatacenterClaimSpec struct {
	// scyllaDBDatacenterPoolName specifies the name of the ScyllaDBDatacenterPool provisioning the ScyllaDBDatacenter.
	ScyllaDBDatacenterPoolName string `json:"scyllaDBDatacenterPoolName"`

	// dnsDomains is a list of DNS domains this cluster is reachable by.
	// These domains are used when setting up the infrastructure, like certificates.
	// EXPERIMENTAL. Do not rely on any particular behaviour controlled by this field.
	// +optional
	DNSDomains []string `json:"dnsDomains,omitempty"`

	// certificateOptions specify parameters related to customizing certificate configuration for the bound ScyllaDBDatacenter.
	// +optional
	CertificateOptions *CertificateOptions `json:"certificateOptions,omitempty"`

	// exposeOptions specifies parameters related to exposing backends of the bound ScyllaDBDatacenter.
	// +optional
	ExposeOptions *ExposeOptions `json:"exposeOptions,omitempty"`
}

type CertificateOptions struct {
	// servingCASecretName references a kubernetes.io/tls type secret containing the TLS cert and key of the serving CA.
	ServingCASecretName string `json:"servingCASecretName"`

	// clientCASecretName references a kubernetes.io/tls type secret containing the TLS cert and key of the client CA.
	ClientCASecretName string `json:"clientCASecretName"`
}

type ScyllaDBDatacenterClaimStatus struct {
	// observedGeneration reflects the most recently observed generation of the ScyllaDBDatacenterClaim.
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty"`

	// conditions reflect the latest observed conditions of the ScyllaDBDatacenterClaim's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// scyllaDBDatacenterName reflects the name of the bound ScyllaDBDatacenter.
	// +optional
	ScyllaDBDatacenterName *string `json:"scyllaDBDatacenterName,omitempty"`
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

type ScyllaDBDatacenterClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec defines the specification of the ScyllaDBDatacenterClaim.
	Spec ScyllaDBDatacenterClaimSpec `json:"spec,omitempty"`

	// status reflects the observed state of the ScyllaDBDatacenterClaim.
	Status ScyllaDBDatacenterClaimStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ScyllaDBDatacenterClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScyllaDBDatacenterClaim `json:"items"`
}
