/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/apirator/apirator/internal/operation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AdapterFactory interface {
	CreateAPIMockAdapter(ctx context.Context, key client.ObjectKey) (APIMockAdapter, error)
}

type APIMockAdapter interface {
	EnsureDefinitionIsValid(ctx context.Context) (*operation.Result, error)
	EnsureIsInitialized(ctx context.Context) (*operation.Result, error)
	EnsureFinalizer(ctx context.Context) (*operation.Result, error)
	EnsureConfigMap(ctx context.Context) (*operation.Result, error)
	EnsureDeployment(ctx context.Context) (*operation.Result, error)
	EnsureService(ctx context.Context) (*operation.Result, error)
	EnsureIngress(ctx context.Context) (*operation.Result, error)
	EnsureDeploymentAvailability(ctx context.Context) (*operation.Result, error)
}
