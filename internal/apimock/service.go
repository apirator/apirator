package apimock

import (
	"context"
	"fmt"

	api "github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/tracing"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct {
	client client.Client
	logger logr.Logger
	scheme *runtime.Scheme
}

func NewService(client client.Client, scheme *runtime.Scheme) *Service {
	return &Service{
		client: client,
		logger: ctrl.Log.WithName("services").WithName("APIMock"),
		scheme: scheme,
	}
}

func (s *Service) LookupResourceAdapter(ctx context.Context, key client.ObjectKey) (*Adapter, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	hm := &api.APIMock{}
	err := s.client.Get(ctx, key, hm)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		span.SetError(err)
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}

	return newAdapter(hm, s), nil
}

func (s *Service) Apply(ctx context.Context, inv inventory.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := s.logger.WithValues("trace", span.String())
	defer span.Finish()

	for _, obj := range inv.Create {
		if err := s.client.Create(ctx, obj); err != nil {
			span.SetError(err)
			return fmt.Errorf("failed to create %T %q: %w", obj, obj.GetName(), err)
		}
		log.Info(fmt.Sprintf("%T %q created", obj, obj.GetName()))
	}

	for _, obj := range inv.Update {
		if err := s.client.Update(ctx, obj); err != nil {
			span.SetError(err)
			return fmt.Errorf("failed to update %T %q: %w", obj, obj.GetName(), err)
		}
		log.Info(fmt.Sprintf("%T %q updated", obj, obj.GetName()))
	}

	for _, obj := range inv.Delete {
		if err := s.client.Update(ctx, obj); err != nil {
			span.SetError(err)
			return fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err)
		}
		log.Info(fmt.Sprintf("%T %q deleted", obj, obj.GetName()))
	}

	return nil
}
