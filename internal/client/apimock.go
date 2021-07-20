package client

import (
	"context"
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (k *Kubernetes) GetAPIMock(ctx context.Context, key client.ObjectKey) (*v1alpha1.APIMock, error) {
	r := new(v1alpha1.APIMock)
	err := k.Get(ctx, key, r)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return r, nil
}

func (k *Kubernetes) ListAPIMocks(ctx context.Context, namespace string) (*v1alpha1.APIMockList, error) {
	opts := []client.ListOption{
		client.InNamespace(namespace),
	}
	list := new(v1alpha1.APIMockList)
	if err := k.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list APIMocks: %w", err)
	}
	return list, nil
}

func (k *Kubernetes) UpdateAPIMock(ctx context.Context, apimock *v1alpha1.APIMock) error {
	if err := k.Update(ctx, apimock); err != nil {
		return fmt.Errorf("failed to update APIMock: %w", err)
	}
	return nil
}

func (k *Kubernetes) UpdateAPIMockStatus(ctx context.Context, apimock *v1alpha1.APIMock) error {
	if err := k.Status().Update(ctx, apimock); err != nil {
		return fmt.Errorf("failed to update APIMock status: %w", err)
	}
	return nil
}
