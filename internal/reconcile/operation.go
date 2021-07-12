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
)

type (
	Operation       func(ctx context.Context, apimock *v1alpha1.APIMock) (*OperationResult, error)
	OperationResult struct {
		RequeueDelay   time.Duration
		RequeueRequest bool
		CancelRequest  bool
	}
)

func (r *OperationResult) RequeueOrCancel() bool {
	return r.RequeueRequest || r.CancelRequest
}

func ContinueResult() *OperationResult {
	return &OperationResult{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  false,
	}
}

func StopResult() *OperationResult {
	return &OperationResult{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  true,
	}
}

func StopProcessing() (result *OperationResult, err error) {
	result = StopResult()
	return
}

func Requeue() (result *OperationResult, err error) {
	result = &OperationResult{
		RequeueDelay:   0,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	return
}

func RequeueWithError(errIn error) (result *OperationResult, err error) {
	result = &OperationResult{
		RequeueDelay:   0,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	err = errIn
	return
}

func RequeueOnErrorOrStop(errIn error) (result *OperationResult, err error) {
	result = &OperationResult{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  true,
	}
	err = errIn
	return
}

func RequeueOnErrorOrContinue(errIn error) (result *OperationResult, err error) {
	result = &OperationResult{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  false,
	}
	err = errIn
	return
}

func RequeueAfter(delay time.Duration) (result *OperationResult, err error) {
	result = &OperationResult{
		RequeueDelay:   delay,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	return
}

func ContinueProcessing() (result *OperationResult, err error) {
	result = ContinueResult()
	return
}
