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

	"github.com/apirator/apirator/api/v1alpha1"

	"github.com/apirator/apirator/internal/reconcile"
)

// APIMockAdapter have reconcile subroutines that performs actions to keep the desired state of the cluster
type APIMockAdapter interface {
	EnsureStatus(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
	EnsureDefinitionIsValid(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
	EnsureConfigMap(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
	EnsureDeployment(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
	EnsureService(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
	EnsureIngress(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
	EnsureIngressFinalizer(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
	EnsureDeploymentAvailability(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error)
}
