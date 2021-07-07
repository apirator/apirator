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
	appsv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	DeploymentBuilder
	DeploymentReader
	GenericObjectWriter
}

type DeploymentBuilder interface {
	DeploymentFor(resource *v1alpha1.APIMock) (*appsv1.Deployment, error)
}

type DeploymentReader interface {
	ListDeployments(ctx context.Context, resource *v1alpha1.APIMock) (*appsv1.DeploymentList, error)
}

func (d *Deployment) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error) {
	desired, err := d.DeploymentFor(apimock)
	if err != nil {
		return nil, err
	}

	list, err := d.ListDeployments(ctx, apimock)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForDeployments(list.Items, []appsv1.Deployment{*desired})
	err = d.Apply(ctx, inv)
	if err != nil {
		return nil, err
	}

	return reconcile.ContinueProcessing()
}
