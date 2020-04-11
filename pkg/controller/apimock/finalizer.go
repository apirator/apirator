package apimock

import (
	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/redhat-cop/operator-utils/pkg/util"
)

// Manage the clean up logic
func (r *ReconcileAPIMock) ManageCleanUpLogic(mock *apirator.APIMock) error {
	return r.RemoveEntryFromIngress(mock)
}

// Set up initial instance
func (r *ReconcileAPIMock) SetUp(mock *apirator.APIMock) error {
	return r.addIngressFinalizer(mock)
}

// ensure the apimock finalizers
func (r *ReconcileAPIMock) EnsureFinalizer(mock *apirator.APIMock) error {
	// Not marked for deletion
	if util.IsBeingDeleted(mock) {

		// there is nothing to do. Work completed
		if r.hasIngressFinalizer(mock) {
			return nil
		}

		// Remove necessary elements
		err := r.ManageCleanUpLogic(mock)
		if err != nil {
			log.Error(err, "unable to delete mock", "mock", mock)
			return err
		}

		// Remove finalizers
		err = r.removeIngressFinalizer(mock)
		if err != nil {
			log.Error(err, "unable to update instance", "instance", mock)
			return err
		}

	}
	return nil
}
