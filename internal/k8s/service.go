package k8s

import (
	"context"
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/tracing"
	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct{ client client.Client }

func NewService(client client.Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetAPIMock(ctx context.Context, key client.ObjectKey) (*v1alpha1.APIMock, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	r := new(v1alpha1.APIMock)
	err := s.client.Get(ctx, key, r)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return r, nil
}

func (s *Service) Apply(ctx context.Context, inv inventory.Object) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := span.Logger()
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
		if err := s.client.Delete(ctx, obj); err != nil {
			span.SetError(err)
			return fmt.Errorf("failed to delete %T %q: %w", obj, obj.GetName(), err)
		}
		log.Info(fmt.Sprintf("%T %q deleted", obj, obj.GetName()))
	}

	return nil
}

func (s *Service) ListConfigMaps(resource *v1alpha1.APIMock) (*core.ConfigMapList, error) {
	opts := []client.ListOption{
		client.InNamespace(resource.GetNamespace()),
		client.MatchingLabels(resource.MatchLabels()),
	}
	list := new(core.ConfigMapList)
	if err := s.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}
	return list, nil
}

func (s *Service) ListDeployments(resource *v1alpha1.APIMock) (*appsv1.DeploymentList, error) {
	opts := []client.ListOption{
		client.InNamespace(resource.Namespace),
		client.MatchingLabels(resource.MatchLabels()),
	}
	list := new(appsv1.DeploymentList)
	if err := s.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Deployments: %w", err)
	}
	return list, nil
}

func (s *Service) ListServices(resource *v1alpha1.APIMock) (*core.ServiceList, error) {
	opts := []client.ListOption{
		client.InNamespace(resource.GetNamespace()),
		client.MatchingLabels(resource.MatchLabels()),
	}
	list := new(core.ServiceList)
	if err := s.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Services: %w", err)
	}
	return list, nil
}
