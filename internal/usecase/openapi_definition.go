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

type OpenAPIValidator interface {
	Validate(definition string) error
}

type OpenAPIDefinition struct {
	validator OpenAPIValidator
	writer    APIMockStatusWriter
}

func NewOpenAPIDefinition(validator OpenAPIValidator, writer APIMockStatusWriter) *OpenAPIDefinition {
	return &OpenAPIDefinition{validator: validator, writer: writer}
}

func (v *OpenAPIDefinition) EnsureDefinitionIsValid(ctx context.Context, apimock *v1alpha1.APIMock) (*reconcile.OperationResult, error) {
	err := v.validator.Validate(apimock.Spec.Definition)
	if err != nil {
		if apimock.SetValidatedConditionFalse(err) {
			return reconcile.RequeueOnErrorOrStop(v.writer.UpdateAPIMockStatus(ctx, apimock))
		}
		return reconcile.StopProcessing()
	}

	if apimock.SetValidatedConditionTrue() {
		return reconcile.RequeueOnErrorOrStop(v.writer.UpdateAPIMockStatus(ctx, apimock))
	}

	return reconcile.ContinueProcessing()
}
