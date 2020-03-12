package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExternalConfigMapSpec defines the desired state of ExternalConfigMap
type ExternalConfigMapSpec struct {
	// +kubebuilder:validation:MinLength=1

	// Name of the ExternalBackend resource that is used to get a secret value.
	BackendName string `json:"backendName"`

	// +kubebuilder:validation:MinLength=1

	// Key in the backend that holds a secret value.
	Key string `json:"key"`
}

// ExternalConfigMapStatus defines the observed state of ExternalConfigMap
type ExternalConfigMapStatus struct {
	// Information when was the last time the secret was successfully synced.
	// +optional
	LastSyncedTime *metav1.Time `json:"lastSyncedTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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

// +kubebuilder:docs-gen:collapse=Root Object Definitions
