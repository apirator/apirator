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

type APIMockPhase string

// These are the valid statuses of APIMock.
const (
	PodPending APIMockPhase = "Pending"
	PodRunning APIMockPhase = "Running"
	PodFailed  APIMockPhase = "Failed"
	PodUnknown APIMockPhase = "Unknown"
)

type APIMockConditionType string

// These are valid conditions of a deployment APIMock.
const (
	APIMockProvisioned        APIMockConditionType = "Provisioned"
	APIMockError              APIMockConditionType = "Error"
	APIMockInvalidOAS         APIMockConditionType = "InvalidOAS"
	APIMockWaitingAnnotations APIMockConditionType = "WaitingAnnotations"
)

// APIMockCondition contains details for the current condition of this APIMock.
type APIMockCondition struct {
	// Type of APIMock condition.
	Type APIMockConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`

	// The last time this condition was probed.
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

// APIMockSpec defines the desired state of APIMock
type APIMockSpec struct {
	Definition        string            `json:"definition,omitempty"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition,omitempty"`
	Watch             bool              `json:"watch,omitempty"`
	Selector          map[string]string `json:"selector,omitempty"`
	Host              string            `json:"host,omitempty"`
	Initialized       bool              `json:"initialized,omitempty"`
}

// APIMockStatus defines the observed state of APIMock
type APIMockStatus struct {
	// The phase of a APIMock is a simple, high-level summary of where the APIMock is in its lifecycle.
	// +optional
	Phase APIMockPhase `json:"phase,omitempty"`

	// Represents the latest available observations of a APIMock's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []APIMockCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// ServiceDefinition it will "link" the mock with created service
type ServiceDefinition struct {
	Port        int                `json:"port,omitempty"`
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
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
