package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExternalBackendSpec defines the desired state of ExternalBackend
type ExternalBackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ExternalBackend. Edit ExternalBackend_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ExternalBackendStatus defines the observed state of ExternalBackend
type ExternalBackendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
