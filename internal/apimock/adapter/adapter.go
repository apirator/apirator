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

package adapter

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/operation"
)

type Adapter struct {
	*UserCases

	customresource *v1alpha1.APIMock
}

func newAdapter(userCases *UserCases, customresource *v1alpha1.APIMock) *Adapter {
	return &Adapter{UserCases: userCases, customresource: customresource}
}

func (a *Adapter) EnsureConfigMap(ctx context.Context) (*operation.Result, error) {
	return a.ConfigMap.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureDefinitionIsValid(ctx context.Context) (*operation.Result, error) {
	return a.OpenAPIDefinition.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureDeployment(ctx context.Context) (*operation.Result, error) {
	return a.Deployment.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureIngress(ctx context.Context) (*operation.Result, error) {
	return a.Ingress.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureService(ctx context.Context) (*operation.Result, error) {
	return a.Service.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureStatus(ctx context.Context) (*operation.Result, error) {
	return a.Status.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureIngressFinalizer(ctx context.Context) (*operation.Result, error) {
	return a.IngressFinalizer.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureDeploymentAvailability(ctx context.Context) (*operation.Result, error) {
	return a.DeploymentAvailability.Ensure(ctx, a.customresource)
}
