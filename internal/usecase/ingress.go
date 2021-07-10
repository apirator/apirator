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

	"github.com/apirator/apirator/internal/reconcile"
	networkingv1 "k8s.io/api/networking/v1"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/inventory"
)

type (
	APIMockReader interface {
		ListAPIMocks(ctx context.Context, namespace string) (*v1alpha1.APIMockList, error)
	}
	IngressBuilder interface {
		IngressesFor(resources *v1alpha1.APIMockList) (*networkingv1.IngressList, error)
	}
	IngressReader interface {
		ListIngresses(ctx context.Context, resource *v1alpha1.APIMock) (*networkingv1.IngressList, error)
	}
)

type Ingress struct {
	builder   IngressBuilder
	ingreader IngressReader
	apireader APIMockReader
	writer    GenericObjectWriter
}

func NewIngress(builder IngressBuilder, ingreader IngressReader, apireader APIMockReader, writer GenericObjectWriter) *Ingress {
	return &Ingress{builder: builder, ingreader: ingreader, apireader: apireader, writer: writer}
}

func (i *Ingress) EnsureIngress(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error) {
	apimocks, err := i.apireader.ListAPIMocks(ctx, apimock.GetNamespace())
	if err != nil {
		return nil, err
	}

	desired, err := i.builder.IngressesFor(apimocks)
	if err != nil {
		return nil, err
	}

	list, err := i.ingreader.ListIngresses(ctx, apimock)
	if err != nil {
		return nil, err
	}

	inv := inventory.ForIngresses(list.Items, desired.Items)
	err = i.writer.Apply(ctx, inv)
	if err != nil {
		return nil, err
	}

	return reconcile.ContinueProcessing()
}
