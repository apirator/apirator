package client

import (
	"context"
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (k *Kubernetes) ListIngresses(ctx context.Context, resource *v1alpha1.APIMock) (*networkingv1beta1.IngressList, error) {
	opts := []client.ListOption{
		client.InNamespace(resource.GetNamespace()),
		client.MatchingLabels(map[string]string{"app.kubernetes.io/managed-by": "apirator"}),
	}
	list := new(networkingv1beta1.IngressList)
	if err := k.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Ingresses: %w", err)
	}
	return list, nil
}
