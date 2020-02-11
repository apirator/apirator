package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	PROVISIONED = "PROVISIONED"
	ERROR       = "ERROR"
)

// APIMockSpec defines the desired state of APIMock
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string
type APIMockSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Definition    string `json:"definition,omitempty"`
	Port          int    `port:"definition,omitempty"`
	ContainerPort int    `json:"containerPort,omitempty"`
}

// APIMockStatus defines the observed state of APIMock
type APIMockStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:validation:Enum=PROVISIONED;ERROR;
	Phase string `json:"phase,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIMock is the Schema for the apimocks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=apimocks,scope=Namespaced
type APIMock struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIMockSpec   `json:"spec,omitempty"`
	Status APIMockStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIMockList contains a list of APIMock
type APIMockList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIMock `json:"items"`
}

func init() {
	SchemeBuilder.Register(&APIMock{}, &APIMockList{})
}
