package apimock

import (
	"context"

	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/k8s/deployments"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	appsv1 "k8s.io/api/apps/v1"
)

func (a *Adapter) EnsureDeployment(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := deployments.FromAPIMock(a.scheme, a.resource)
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := a.svc.ListDeployments(a.resource)
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForDeployments(list.Items, []appsv1.Deployment{*desired})
	err = a.svc.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}
