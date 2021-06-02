package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/resources"
	"github.com/apirator/apirator/internal/tracing"
)

type Ingress struct {
	*resources.Builder
	*k8s.Service
}

func (i *Ingress) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := i.IngressFor(apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	_ = desired
	// TODO

	return operation.ContinueProcessing()
}
