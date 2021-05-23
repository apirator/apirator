package mock

import (
	"context"
	"fmt"

	api "github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/tracing"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (s *Service) LookupResource(ctx context.Context, key client.ObjectKey) (*api.APIMock, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	hm := &api.APIMock{}
	err := s.Client.Get(ctx, key, hm)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		span.SetError(err)
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}

	return hm, nil
}
