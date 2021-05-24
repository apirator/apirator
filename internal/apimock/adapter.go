package apimock

import (
	"context"
	api "github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Adapter struct {
	*api.APIMock
	svc *Service
	Log logr.Logger
}

func newAdapter(APIMock *api.APIMock, svc *Service) *Adapter {
	return &Adapter{
		APIMock: APIMock,
		Log:     ctrl.Log.WithName("adapters").WithName("APIMock"),
		svc:     svc,
	}
}

func (a *Adapter) EnsureConfigMap(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	_ = a.Log.WithValues("trace", span.String())
	defer span.Finish()

	return operation.ContinueProcessing()
}
