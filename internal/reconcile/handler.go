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

package reconcile

import (
	"context"
	"time"

	"github.com/apirator/apirator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Handler struct {
	operations []Operation
}

func NewHandler(operations ...Operation) *Handler {
	return &Handler{operations: operations}
}

func (h *Handler) Handle(ctx context.Context, apimock *v1alpha1.APIMock) (reconcile.Result, error) {
	for _, op := range h.operations {
		result, err := op(ctx, apimock)
		if err != nil {
			return RequeueOnErr(err)
		}
		if result == nil || result.CancelRequest {
			return DoNotRequeue()
		}
		if result.RequeueRequest {
			return RequeueOnErrAfter(err, result.RequeueDelay)
		}
	}
	return DoNotRequeue()
}

func DoNotRequeue() (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func RequeueOnErr(err error) (reconcile.Result, error) {
	return reconcile.Result{}, err
}

func RequeueOnErrAfter(err error, duration time.Duration) (reconcile.Result, error) {
	return reconcile.Result{RequeueAfter: duration}, err
}
