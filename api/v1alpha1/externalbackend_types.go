package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExternalBackendSpec defines the desired state of ExternalBackend
type ExternalBackendSpec struct {
	// Specifies the backend type.
	// Valid values are:
	// - "Consul";
	// - "Vault";
	Type BackendType `json:"type"`

	// +kubebuilder:validation:MinLength=1

	// Secret name that hold backend configuration.
	SecretName string `json:"secretName"`
}

// BackendType.
// +kubebuilder:validation:Enum=Consul;Vault
type BackendType string

const (
	// Consul.
	Consul BackendType = "Consul"

	// Vault.
	Vault BackendType = "Vault"
)

// ExternalBackendStatus defines the observed state of ExternalBackend
type ExternalBackendStatus struct {
	// Information about the backend connection status.
	Connected *bool `json:"connected"`
}

// +kubebuilder:object:root=true

// ExternalBackend is the Schema for the externalbackends API
type ExternalBackend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalBackendSpec   `json:"spec,omitempty"`
	Status ExternalBackendStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExternalBackendList contains a list of ExternalBackend
type ExternalBackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalBackend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExternalBackend{}, &ExternalBackendList{})
}
