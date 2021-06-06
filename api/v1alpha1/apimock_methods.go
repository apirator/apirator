package v1alpha1

import (
	"github.com/apirator/apirator/api/v1alpha1/condition"
	"github.com/apirator/apirator/internal/openapi"
	"k8s.io/apimachinery/pkg/api/meta"
)

func (in *APIMock) MatchLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       in.GetName(),
		"app.kubernetes.io/managed-by": "apirator",
	}
}

func (in *APIMock) SetConditionForError(err error) (updated bool) {
	if _, ok := err.(*openapi.InvalidDefinitionError); ok {
		message := "OpenAPI definition has " + err.Error()
		existingCondition := meta.FindStatusCondition(in.Status.Conditions, condition.Validated)
		updated = existingCondition == nil || existingCondition.Message != message
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewOpenAPIValidation(false, message))
	} else {
		message := err.Error()
		existingCondition := meta.FindStatusCondition(in.Status.Conditions, condition.Error)
		updated = existingCondition == nil || existingCondition.Message != message
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewError(true, message))
	}
	return updated
}

func (in *APIMock) SetConditionForValidOpenAPI() (updated bool) {
	if meta.IsStatusConditionFalse(in.Status.Conditions, condition.Validated) {
		message := "OpenAPI definition has been successfully validated."
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewOpenAPIValidation(true, message))
		updated = true
	}
	return updated
}

func (in *APIMock) SetConditionForAvailability(available bool) (updated bool) {
	if available && !meta.IsStatusConditionTrue(in.Status.Conditions, condition.Available) {
		message := "Deployment has minimum availability."
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewAvailable(true, message))
		updated = true
	}
	if !available && !meta.IsStatusConditionFalse(in.Status.Conditions, condition.Available) {
		message := "Deployment has no minimum availability."
		meta.SetStatusCondition(&in.Status.Conditions, condition.NewAvailable(false, message))
		updated = true
	}
	return updated
}
