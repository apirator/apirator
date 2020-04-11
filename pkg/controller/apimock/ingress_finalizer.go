package apimock

import (
	"context"
	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/redhat-cop/operator-utils/pkg/util"
)

const (
	IngressFinalizerName = "ingress.finalizers.apirator.io"
)

// add ingress finalizer
func (r *ReconcileAPIMock) addIngressFinalizer(mock *apirator.APIMock) error {
	log.Info("Adding Ingress Finalizer for the APIMock")
	mock.AddFinalizer(IngressFinalizerName)
	err := r.client.Update(context.TODO(), mock)
	if err != nil {
		log.Error(err, "[Adding] - Failed to update APIMock with finalizer")
		return err
	}
	return nil
}

// remove ingress finalizer
func (r *ReconcileAPIMock) removeIngressFinalizer(mock *apirator.APIMock) error {
	log.Info("Remove Ingress Finalizer for the APIMock")
	mock.RemoveFinalizer(IngressFinalizerName)
	err := r.client.Update(context.TODO(), mock)
	if err != nil {
		log.Error(err, "[Removing] - Failed to update APIMock with finalizer")
		return err
	}
	return nil
}

func (r *ReconcileAPIMock) hasIngressFinalizer(mock *apirator.APIMock) bool {
	return util.HasFinalizer(mock, IngressFinalizerName)
}
