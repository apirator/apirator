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

package apimock

import (
	"context"

	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
)

func (r *ReconcileAPIMock) markAsSuccessful(obj *apirator.APIMock) error {
	return r.updateStatus(obj, apirator.PROVISIONED)
}

func (r *ReconcileAPIMock) markAsFailure(obj *apirator.APIMock) error {
	return r.updateStatus(obj, apirator.ERROR)
}

func (r *ReconcileAPIMock) markAsInvalidOAS(obj *apirator.APIMock) error {
	return r.updateStatus(obj, apirator.INVALID_OAS)
}

func (r *ReconcileAPIMock) updateStatus(obj *apirator.APIMock, status string) error {
	if obj.Status.Phase != status {
		obj.Status.Phase = status
		err := r.client.Update(context.TODO(), obj)
		if err != nil {
			return err
		}
	}
	return nil
}
