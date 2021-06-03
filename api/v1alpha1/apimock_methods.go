package v1alpha1

import (
	"github.com/apirator/apirator/api/v1alpha1/condition"
	"github.com/apirator/apirator/api/v1alpha1/phase"
	"github.com/apirator/apirator/internal/openapi"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *APIMock) MatchLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       in.GetName(),
		"app.kubernetes.io/managed-by": "apirator",
	}
}

func (in *APIMock) SetStatusConditionForError(err error) (updated bool) {
	if _, ok := err.(*openapi.InvalidDefinitionError); ok {
		message := "OpenAPI definition has " + err.Error()
		existingCondition := meta.FindStatusCondition(in.Status.Conditions, condition.Validated)
		updated = existingCondition == nil || existingCondition.Message != message
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewValidOpenAPIDefinition(metav1.ConditionFalse, message))
	} else {
		message := err.Error()
		existingCondition := meta.FindStatusCondition(in.Status.Conditions, condition.Error)
		updated = existingCondition == nil || existingCondition.Message != message
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewError(metav1.ConditionTrue, message))
	}
	updatedStatus := in.UpdateStatus()
	return updated || updatedStatus
}

func (in *APIMock) SetStatusConditionForValidOpenAPI() (updated bool) {
	if meta.IsStatusConditionFalse(in.Status.Conditions, condition.Validated) {
		message := "OpenAPI definition has been successfully validated."
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewValidOpenAPIDefinition(metav1.ConditionTrue, message))
		updated = true
	}
	updatedStatus := in.UpdateStatus()
	return updated || updatedStatus
}

func (in *APIMock) UpdateStatus() bool {
	if meta.IsStatusConditionFalse(in.Status.Conditions, condition.Validated) && in.Status.Phase != phase.ConfigError {
		in.Status.Phase = phase.ConfigError
		return true
	}
	if meta.IsStatusConditionTrue(in.Status.Conditions, condition.Validated) && in.Status.Phase == phase.ConfigError {
		in.Status.Phase = phase.Pending
		return true
	}
	return false
}
