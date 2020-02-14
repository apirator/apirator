package apimock

import (
	"context"

	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
)

func (r *ReconcileAPIMock) markAsSuccessful(obj *apirator.APIMock) error {
	return r.updateStatus(obj, apirator.PROVISIONED)
}

func (r *ReconcileAPIMock) markAsFailure(obj *apirator.APIMock) error {
	return r.updateStatus(obj, apirator.ERROR)
}

func (r *ReconcileAPIMock) markAsInvalidOAS(obj *apirator.APIMock) error {
	return r.updateStatus(obj, apirator.INVALID_OAS)
}

func (r *ReconcileAPIMock) updateStatus(obj *apirator.APIMock, status string) error {
	if obj.Status.Phase != status {
		obj.Status.Phase = status
		err := r.client.Update(context.TODO(), obj)
		if err != nil {
			return err
		}
	}
	return nil
}
