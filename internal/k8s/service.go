package k8s

import (
	"context"
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/tracing"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
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

func (s *Service) UpdateAPIMockStatus(ctx context.Context, apimock *v1alpha1.APIMock) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if err := s.client.Status().Update(ctx, apimock); err != nil {
		return fmt.Errorf("failed to update APIMock status: %w", err)
	}
	log := span.Logger()
	log.Info("APIMock status updated")
	return nil
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

func (s *Service) ListConfigMaps(ctx context.Context, resource *v1alpha1.APIMock) (*corev1.ConfigMapList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(resource.GetNamespace()),
		client.MatchingLabels(resource.MatchLabels()),
	}
	list := new(corev1.ConfigMapList)
	if err := s.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}
	return list, nil
}

func (s *Service) ListDeployments(ctx context.Context, resource *v1alpha1.APIMock) (*appsv1.DeploymentList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

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

func (s *Service) ListServices(ctx context.Context, resource *v1alpha1.APIMock) (*corev1.ServiceList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(resource.GetNamespace()),
		client.MatchingLabels(resource.MatchLabels()),
	}
	list := new(corev1.ServiceList)
	if err := s.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Services: %w", err)
	}
	return list, nil
}

func (s *Service) GetDeploymentStatus(ctx context.Context, resource *v1alpha1.APIMock) (*DeploymentStatus, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	key := types.NamespacedName{Namespace: resource.GetNamespace(), Name: resource.GetName()}
	dep := &appsv1.Deployment{}
	if err := s.client.Get(ctx, key, dep); err != nil {
		return nil, fmt.Errorf("failed to lookup Deployment: %w", err)
	}
	return &DeploymentStatus{DeploymentStatus: &dep.Status}, nil
}
