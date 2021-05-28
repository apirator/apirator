/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

var APIMockLabels = map[string]string{"app.kubernetes.io/managed-by": "apirator"}

const (
	Provisioned          = "Provisioned"
	Error                = "Error"
	InvalidOas           = "InvalidOAS"
	WaitingAnnotations   = "WaitingAnnotations"
	IngressTag           = "ingress"
	NamespaceTag         = "namespace"
	IngressFinalizerName = "ingress.finalizers.apirator.io"
)

type Step struct {
	Action      string      `json:"action,omitempty"`
	LastUpdate  metav1.Time `json:"lastUpdate,omitempty"`
	Description string      `json:"description,omitempty"`
}

// APIMockSpec defines the desired state of APIMock
type APIMockSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Definition        string            `json:"definition,omitempty"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition,omitempty"`
	Watch             bool              `json:"watch,omitempty"`
	Selector          map[string]string `json:"selector,omitempty"`
	Host              string            `json:"host,omitempty"`
	Initialized       bool              `json:"initialized,omitempty"`
}

// APIMockStatus defines the observed state of APIMock
type APIMockStatus struct { // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase string `json:"phase,omitempty"`
	Steps []Step `json:"steps"`
}

// Service Definition it will "link" the mock with created service
type ServiceDefinition struct {
	Port        int                `json:"port,omitempty"`
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string

// APIMock is the Schema for the apimocks API
type APIMock struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIMockSpec   `json:"spec,omitempty"`
	Status APIMockStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// APIMockList contains a list of APIMock
type APIMockList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIMock `json:"items"`
}

func init() {
	SchemeBuilder.Register(&APIMock{}, &APIMockList{})
}
