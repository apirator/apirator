package client

import (
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/apirator/apirator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func (k *Kubernetes) ListConfigMaps(ctx context.Context, apimock *v1alpha1.APIMock) (*corev1.ConfigMapList, error) {
	opts := []client.ListOption{
		client.InNamespace(apimock.GetNamespace()),
		client.MatchingLabels(apimock.MatchLabels()),
	}
	list := new(corev1.ConfigMapList)
	if err := k.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}
	return list, nil
}
