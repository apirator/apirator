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
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/resources"
	"github.com/apirator/apirator/internal/tracing"
	networkingv1 "k8s.io/api/networking/v1"
)

type Ingress struct {
	*resources.Builder
	*k8s.Service
}

func (i *Ingress) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	apimocks, err := i.ListAPIMocks(ctx, apimock.GetNamespace())
	if err != nil {
		return nil, span.HandleError(err)
	}

	desired, err := i.IngressFor(apimocks)
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := i.ListIngresses(ctx, apimock)
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForIngresses(list.Items, []networkingv1.Ingress{*desired})
	err = i.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}
