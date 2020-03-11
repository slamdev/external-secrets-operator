package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExternalConfigMapSpec defines the desired state of ExternalConfigMap
type ExternalConfigMapSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ExternalConfigMap. Edit ExternalConfigMap_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ExternalConfigMapStatus defines the observed state of ExternalConfigMap
type ExternalConfigMapStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ExternalConfigMap is the Schema for the externalconfigmaps API
type ExternalConfigMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalConfigMapSpec   `json:"spec,omitempty"`
	Status ExternalConfigMapStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExternalConfigMapList contains a list of ExternalConfigMap
type ExternalConfigMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalConfigMap `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExternalConfigMap{}, &ExternalConfigMapList{})
}
