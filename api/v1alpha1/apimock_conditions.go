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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *APIMock) IsAvailableConditionFalse() bool {
	return meta.IsStatusConditionFalse(in.Status.Conditions, APIMockAvailable)
}

func (in *APIMock) IsAvailableConditionTrue() bool {
	return meta.IsStatusConditionTrue(in.Status.Conditions, APIMockAvailable)
}

func (in *APIMock) IsValidatedConditionFalse() bool {
	return meta.IsStatusConditionFalse(in.Status.Conditions, APIMockValidated)
}

func (in *APIMock) IsValidatedConditionTrue() bool {
	return meta.IsStatusConditionTrue(in.Status.Conditions, APIMockValidated)
}

func (in *APIMock) SetAvailableConditionFalse() bool {
	if !in.IsAvailableConditionFalse() {
		meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
			Type:    APIMockAvailable,
			Status:  metav1.ConditionFalse,
			Reason:  "MinimumReplicasAvailable",
			Message: "Deployment has no minimum availability.",
		})
		return true
	}
	return false
}

func (in *APIMock) SetAvailableConditionTrue() bool {
	if !in.IsAvailableConditionTrue() {
		meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
			Type:    APIMockAvailable,
			Status:  metav1.ConditionTrue,
			Reason:  "MinimumReplicasAvailable",
			Message: "Deployment has minimum availability.",
		})
		return true
	}
	return false
}

func (in *APIMock) SetValidatedConditionFalse(err error) bool {
	if !in.IsValidatedConditionFalse() {
		meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
			Type:    APIMockValidated,
			Status:  metav1.ConditionFalse,
			Reason:  "ValidOpenAPIDefinition",
			Message: err.Error(),
		})
		return true
	}
	return false
}

func (in *APIMock) SetValidatedConditionTrue() bool {
	if !in.IsValidatedConditionTrue() {
		meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
			Type:    APIMockValidated,
			Status:  metav1.ConditionTrue,
			Reason:  "ValidOpenAPIDefinition",
			Message: "OpenAPI definition has been successfully validated.",
		})
		return true
	}
	return false
}
