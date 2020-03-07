// Copyright 2020 apirator.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	PROVISIONED = "PROVISIONED"
	ERROR       = "ERROR"
	INVALID_OAS = "INVALID_OAS"
)

// APIMockSpec defines the desired state of APIMock
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string
type APIMockSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Definition        string            `json:"definition,omitempty"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition,omitempty"`
	Watch             bool              `json:"watch,omitempty"`
}

// APIMockStatus defines the observed state of APIMock
type APIMockStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:validation:Enum=PROVISIONED;ERROR;INVALID_OAS;
	Phase string `json:"phase,omitempty"`
}

type ServiceDefinition struct {
	Port        int            `json:"port,omitempty"`
	ServiceType v1.ServiceType `json:"serviceType,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIMock is the Schema for the apimocks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=apimocks,scope=Namespaced
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string
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
