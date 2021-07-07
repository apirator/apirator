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
	"github.com/apirator/apirator/internal/reconcile"
)

type Status struct {
	APIMockStatusWriter
}

func (v *Status) EnsureStatus(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error) {
	if apimock.Status.Phase == "" {
		apimock.Status.Phase = v1alpha1.PendingPhase
		return reconcile.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
	}

	if apimock.IsValidatedConditionFalse() {
		if apimock.Status.Phase != v1alpha1.ErrorPhase {
			apimock.Status.Phase = v1alpha1.ErrorPhase
			return reconcile.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return reconcile.ContinueProcessing()
	}

	if apimock.IsAvailableConditionFalse() {
		if apimock.Status.Phase != v1alpha1.PendingPhase {
			apimock.Status.Phase = v1alpha1.PendingPhase
			return reconcile.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return reconcile.ContinueProcessing()
	}

	if apimock.IsValidatedConditionTrue() && apimock.IsAvailableConditionTrue() {
		if apimock.Status.Phase != v1alpha1.RunningPhase {
			apimock.Status.Phase = v1alpha1.RunningPhase
			return reconcile.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return reconcile.ContinueProcessing()
	}

	return reconcile.ContinueProcessing()
}
