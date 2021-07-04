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

const (
	APIMockValidated string = "Validated"
	APIMockAvailable        = "Available"
)

// APIMockSpec defines the desired state of APIMock
type APIMockSpec struct {
	//+kubebuilder:validation:Required
	Definition string   `json:"definition,omitempty"`
	Service    Service  `json:"service,omitempty"`
	Ingress    *Ingress `json:"ingress,omitempty"`
	Watch      bool     `json:"watch,omitempty"`
}

type Phase string

// These are the valid statuses of APIMock.
const (
	Pending Phase = "Pending"
	Running Phase = "Running"
	Error   Phase = "Error"
	Failed  Phase = "Failed"
)

// APIMockStatus defines the observed state of APIMock
type APIMockStatus struct {
	// The phase of a APIMock is a simple, high-level summary of where the APIMock is in its lifecycle.
	// +optional
	Phase Phase `json:"phase,omitempty"`

	// Represents the latest available observations of a APIMock's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// Service it will "link" the mock with created service
type Service struct {
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=65535
	Port int `json:"port,omitempty"`
	//+kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	Type corev1.ServiceType `json:"type,omitempty"`
	//+kubebuilder:validation:Enum=Local;Cluster
	ExternalTrafficPolicy corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	Annotations           map[string]string                       `json:"annotations,omitempty"`
}

// Ingress will configure the resource that allows you to access to your mock API
type Ingress struct {
	Hostname    string            `json:"hostname"`
	Path        string            `json:"path"`
	PathType    *string           `json:"pathType,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CertManager *bool             `json:"certManager,omitempty"`
	TLS         *TLS              `json:"tls,omitempty"`
}

// TLS enables configuration for the hostname defined at ingress.hostname parameter
type TLS struct {
	SecretName string `json:"secretName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string

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
