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
	"time"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/reconcile"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type DeploymentAvailability struct {
	APIMockStatusWriter
	DeploymentStatusReader
}

type DeploymentStatusReader interface {
	GetDeploymentStatus(ctx context.Context, apimock *v1alpha1.APIMock) (*appsv1.DeploymentStatus, error)
}

func (d *DeploymentAvailability) EnsureDeploymentAvailability(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error) {
	status, err := d.statusOf(ctx, apimock)
	if err != nil {
		return nil, err
	}

	if status.HasAvailableCondition() {
		if apimock.SetAvailableConditionTrue() {
			return reconcile.RequeueOnErrorOrStop(d.UpdateAPIMockStatus(ctx, apimock))
		}
		return reconcile.ContinueProcessing()
	}

	if apimock.SetAvailableConditionFalse() {
		return reconcile.RequeueOnErrorOrStop(d.UpdateAPIMockStatus(ctx, apimock))
	}

	return reconcile.RequeueAfter(10 * time.Second)
}

func (d *DeploymentAvailability) statusOf(ctx context.Context, apimock *v1alpha1.APIMock) (*deployStatus, error) {
	status, err := d.GetDeploymentStatus(ctx, apimock)
	if err != nil {
		return nil, err
	}
	return &deployStatus{status}, nil
}

type deployStatus struct{ *appsv1.DeploymentStatus }

func (d *deployStatus) HasAvailableCondition() bool {
	for _, condition := range d.Conditions {
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
