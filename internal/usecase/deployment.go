package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/resources"
	"github.com/apirator/apirator/internal/tracing"
	appsv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	*resources.Builder
	*k8s.Service
}

func (d *Deployment) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := d.DeploymentFor(apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := d.ListDeployments(apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForDeployments(list.Items, []appsv1.Deployment{*desired})
	err = d.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}
