package client

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/apirator/apirator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
)

func (k *Kubernetes) GetDeploymentStatus(ctx context.Context, resource *v1alpha1.APIMock) (*appsv1.DeploymentStatus, error) {
	key := types.NamespacedName{Namespace: resource.GetNamespace(), Name: resource.GetName()}
	dep := &appsv1.Deployment{}
	if err := k.Get(ctx, key, dep); err != nil {
		return nil, fmt.Errorf("failed to lookup Deployment: %w", err)
	}
	return &dep.Status, nil
}

func (k *Kubernetes) ListDeployments(ctx context.Context, resource *v1alpha1.APIMock) (*appsv1.DeploymentList, error) {
	opts := []client.ListOption{
		client.InNamespace(resource.Namespace),
		client.MatchingLabels(resource.MatchLabels()),
	}
	list := new(appsv1.DeploymentList)
	if err := k.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Deployments: %w", err)
	}
	return list, nil
}
