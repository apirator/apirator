package apimock

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/operation"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Adapter struct {
	logger logr.Logger
	scheme *runtime.Scheme
	svc    *Service

	resource *v1alpha1.APIMock
}

func newAdapter(scheme *runtime.Scheme, svc *Service, resource *v1alpha1.APIMock) *Adapter {
	return &Adapter{
		logger:   ctrl.Log.WithName("adapters").WithName("APIMock"),
		scheme:   scheme,
		svc:      svc,
		resource: resource,
	}
}

func (a *Adapter) EnsureIsInitialized(ctx context.Context) (*operation.Result, error) {
	panic("implement me")
}

func (a *Adapter) EnsureFinalizer(ctx context.Context) (*operation.Result, error) {
	panic("implement me")
}

func (a *Adapter) EnsureIngress(ctx context.Context) (*operation.Result, error) {
	panic("implement me")
}

func (a *Adapter) EnsureProvisionedStatus(ctx context.Context) (*operation.Result, error) {
	panic("implement me")
}
