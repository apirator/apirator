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
	"fmt"
	"github.com/redhat-cop/operator-utils/pkg/util"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	PROVISIONED          = "PROVISIONED"
	ERROR                = "ERROR"
	INVALID_OAS          = "INVALID_OAS"
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
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string
type APIMockSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Definition        string            `json:"definition,omitempty"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition,omitempty"`
	Watch             bool              `json:"watch,omitempty"`
	Selector          map[string]string `json:"selector,omitempty"`
	Host              string            `json:"host,omitempty"`
	Initialized       bool              `json:"initialized,omitempty"`
}

// APIMockStatus defines the observed state of APIMock
type APIMockStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:validation:Enum=PROVISIONED;ERROR;INVALID_OAS;
	Phase string `json:"phase,omitempty"`
	Steps []Step `json:"steps"`
}

// Service Definition it will "link" the mock with created service
type ServiceDefinition struct {
	Port        int            `json:"port,omitempty"`
	ServiceType v1.ServiceType `json:"serviceType,omitempty"`
}

// It indicates if apimock will be exposed in ingress-controller
func (in *APIMock) ExposeInIngress() bool {
	it := in.Spec.Selector[IngressTag]
	ns := in.Spec.Selector[NamespaceTag]
	return len(it) > 0 && len(ns) > 0
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIMock is the Schema for the apimocks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=apimocks,scope=Namespaced
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string
// +kubebuilder:printcolumn:JSONPath=".metadata.annotations['apirator\\.io/cluster-ip']",name=CLUSTER-IP,type=string
// +kubebuilder:printcolumn:JSONPath=".metadata.annotations['apirator\\.io/ports']",name=PORT(S),type=string
type APIMock struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIMockSpec   `json:"spec,omitempty"`
	Status APIMockStatus `json:"status,omitempty"`
}

func (in *APIMock) AnnotatePorts(ports []v1.ServicePort) (updated bool) {
	var p string
	for _, port := range ports {
		if len(p) > 0 {
			p = fmt.Sprintf("%s,%d/%s", p, port.Port, port.Protocol)
		} else {
			p = fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		}
	}
	updated = in.Annotations["apirator.io/ports"] != p
	in.Annotations["apirator.io/ports"] = p
	return updated
}

func (in *APIMock) AnnotateClusterIP(ip string) (updated bool) {
	updated = in.Annotations["apirator.io/cluster-ip"] != ip
	in.Annotations["apirator.io/cluster-ip"] = ip
	return updated
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

// ================================================================
//                        Domain functions
// ================================================================

// Add the desired finalizer
func (in *APIMock) AddFinalizer(finalizerName string) {
	util.AddFinalizer(in, finalizerName)
}

// Remove the desired finalizer
func (in *APIMock) RemoveFinalizer(finalizerName string) {
	util.RemoveFinalizer(in, finalizerName)
}

// check if CR was initialized
func (in *APIMock) IsInitialized() bool {
	if in.Spec.Initialized {
		return true
	} else {
		if in.ExposeInIngress() {
			in.AddFinalizer(IngressFinalizerName)
		}
	}
	in.Spec.Initialized = true
	return false
}

// It describes if the instance has the desired finalizer
func (in *APIMock) HasFinalizer(finalizerName string) bool {
	return util.HasFinalizer(in, finalizerName)
}
