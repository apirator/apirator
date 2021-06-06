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

type Service struct {
	*resources.Builder
	*k8s.Service
}

func (s *Service) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := s.ServiceFor(apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := s.ListServices(ctx, apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForServices(list.Items, []corev1.Service{*desired})
	err = s.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}
