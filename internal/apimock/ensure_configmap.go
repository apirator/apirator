package apimock

import (
	"context"

	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/k8s/configmaps"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	corev1 "k8s.io/api/core/v1"
)

func (a *Adapter) EnsureConfigMap(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := configmaps.FromAPIMock(a.scheme, a.resource)
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := a.svc.ListConfigMaps(a.resource)
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForConfigMaps(list.Items, []corev1.ConfigMap{*desired})
	err = a.svc.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}
