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
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/api/v1alpha1/phase"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
)

type Status struct {
	*k8s.Service
}

func (v *Status) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	if apimock.Status.Phase == "" {
		apimock.Status.Phase = phase.Pending
		log.Info(fmt.Sprintf("Updating status to %q", phase.Pending))
		return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
	}

	if apimock.IsValidatedConditionFalse() {
		if apimock.Status.Phase != phase.Error {
			apimock.Status.Phase = phase.Error
			log.Info(fmt.Sprintf("Updating status to %q", phase.Error))
			return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.ContinueProcessing()
	}

	if apimock.IsAvailableConditionFalse() {
		if apimock.Status.Phase != phase.Pending {
			apimock.Status.Phase = phase.Pending
			log.Info(fmt.Sprintf("Updating status to %q", phase.Pending))
			return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.ContinueProcessing()
	}

	if apimock.IsValidatedConditionTrue() && apimock.IsAvailableConditionTrue() {
		if apimock.Status.Phase != phase.Running {
			apimock.Status.Phase = phase.Running
			log.Info(fmt.Sprintf("Updating status to %q", phase.Running))
			return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.ContinueProcessing()
	}

	return operation.ContinueProcessing()
}
