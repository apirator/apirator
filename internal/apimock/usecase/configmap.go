package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/resources"
	"github.com/apirator/apirator/internal/tracing"
	corev1 "k8s.io/api/core/v1"
)

type ConfigMap struct {
	*resources.Builder
	*k8s.Service
}

func (c *ConfigMap) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := c.ConfigMapFor(apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := c.ListConfigMaps(ctx, apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForConfigMaps(list.Items, []corev1.ConfigMap{*desired})
	err = c.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}
