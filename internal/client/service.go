package client

import (
	"context"
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (k *Kubernetes) ListServices(ctx context.Context, resource *v1alpha1.APIMock) (*corev1.ServiceList, error) {
	opts := []client.ListOption{
		client.InNamespace(resource.GetNamespace()),
		client.MatchingLabels(resource.MatchLabels()),
	}
	list := new(corev1.ServiceList)
	if err := k.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Services: %w", err)
	}
	return list, nil
}
