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

package operation

import (
	"context"
	"time"
)

type Func func(ctx context.Context) (*Result, error)

type Result struct {
	RequeueDelay   time.Duration
	RequeueRequest bool
	CancelRequest  bool
}

func (r *Result) RequeueOrCancel() bool {
	return r.RequeueRequest || r.CancelRequest
}

func ContinueResult() *Result {
	return &Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  false,
	}
}

func StopResult() *Result {
	return &Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  true,
	}
}

func StopProcessing() (result *Result, err error) {
	result = StopResult()
	return
}

func Requeue() (result *Result, err error) {
	result = &Result{
		RequeueDelay:   0,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	return
}

func RequeueWithError(errIn error) (result *Result, err error) {
	result = &Result{
		RequeueDelay:   0,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	err = errIn
	return
}

func RequeueOnErrorOrStop(errIn error) (result *Result, err error) {
	result = &Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  true,
	}
	err = errIn
	return
}

func RequeueOnErrorOrContinue(errIn error) (result *Result, err error) {
	result = &Result{
		RequeueDelay:   0,
		RequeueRequest: false,
		CancelRequest:  false,
	}
	err = errIn
	return
}

func RequeueAfter(delay time.Duration, errIn error) (result *Result, err error) {
	result = &Result{
		RequeueDelay:   delay,
		RequeueRequest: true,
		CancelRequest:  false,
	}
	err = errIn
	return
}

func ContinueProcessing() (result *Result, err error) {
	result = ContinueResult()
	return
}
