package apimock

import (
	"context"

	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/k8s/services"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	corev1 "k8s.io/api/core/v1"
)

func (a *Adapter) EnsureService(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := services.FromAPIMock(a.scheme, a.resource)
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := a.svc.ListServices(a.resource)
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForServices(list.Items, []corev1.Service{*desired})
	err = a.svc.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}
