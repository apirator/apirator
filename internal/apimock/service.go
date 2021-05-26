package apimock

import (
	"context"
	"fmt"

	api "github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/tracing"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewService(client client.Client, scheme *runtime.Scheme) *Service {
	return &Service{client: client, scheme: scheme}
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
