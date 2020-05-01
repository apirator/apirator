package apimock

import (
	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
)

// Manage the clean up logic
func (r *ReconcileAPIMock) manageCleanUpLogic(mock *apirator.APIMock) error {
	return r.RemoveEntryFromIngress(mock)
}
