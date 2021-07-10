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

package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/reconcile"
	corev1 "k8s.io/api/core/v1"
)

type ConfigMap struct {
	ConfigMapBuilder
	ConfigMapReader
	GenericObjectWriter
}

type ConfigMapBuilder interface {
	ConfigMapFor(apimock *v1alpha1.APIMock) (*corev1.ConfigMap, error)
}

type ConfigMapReader interface {
	ListConfigMaps(ctx context.Context, apimock *v1alpha1.APIMock) (*corev1.ConfigMapList, error)
}

func (c *ConfigMap) EnsureConfigMap(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error) {
	desired, err := c.ConfigMapFor(apimock)
	if err != nil {
		return nil, err
	}

	list, err := c.ListConfigMaps(ctx, apimock)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForConfigMaps(list.Items, []corev1.ConfigMap{*desired})
	err = c.Apply(ctx, inv)
	if err != nil {
		return nil, err
	}

	return reconcile.ContinueProcessing()
}
