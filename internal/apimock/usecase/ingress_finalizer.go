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
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const ingressFinalizer string = "finalizer.ingress.apirator.io"

type IngressFinalizer struct {
	*k8s.Service
}

func (i *IngressFinalizer) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := span.Logger()
	defer span.Finish()

	deleted := apimock.GetDeletionTimestamp() != nil
	containsFinalizer := controllerutil.ContainsFinalizer(apimock, ingressFinalizer)
	if !containsFinalizer && !deleted {
		log.Info("adding finalizer to resource", "Finalizer", ingressFinalizer)
		controllerutil.AddFinalizer(apimock, ingressFinalizer)
		return operation.RequeueOnErrorOrStop(i.UpdateAPIMock(ctx, apimock))
	}

	if deleted {
		log.Info("removing finalizer from resource", "Finalizer", ingressFinalizer)
		controllerutil.RemoveFinalizer(apimock, ingressFinalizer)
		return operation.RequeueOnErrorOrStop(i.UpdateAPIMock(ctx, apimock))
	}

	return operation.ContinueProcessing()
}
