// Copyright 2020 apirator.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"context"
	"encoding/json"
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

	return mergeWithDefaults(r)
}

func (s *Service) ListAPIMocks(ctx context.Context, namespace string) (*v1alpha1.APIMockList, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	opts := []client.ListOption{
		client.InNamespace(namespace),
	}
	list := new(v1alpha1.APIMockList)
	if err := s.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list APIMocks: %w", err)
	}
	for i, item := range list.Items {
		if defaults, err := mergeWithDefaults(&item); err != nil {
			return nil, fmt.Errorf("failed to merge APIMocks with default values: %w", err)
		} else {
			list.Items[i] = *defaults
		}
	}
	return list, nil
}

func mergeWithDefaults(horus *v1alpha1.APIMock) (*v1alpha1.APIMock, error) {
	merged := v1alpha1.DefaultAPIMock()
	jb, err := json.Marshal(horus)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jb, &merged)
	if err != nil {
		return nil, err
	}

	return merged, nil
}

func (s *Service) UpdateAPIMockStatus(ctx context.Context, apimock *v1alpha1.APIMock) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if err := s.client.Status().Update(ctx, apimock); err != nil {
		return fmt.Errorf("failed to update APIMock status: %w", err)
	}
	log := span.Logger()
	log.Info("APIMock status updated", "phase", apimock.Status.Phase)
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
