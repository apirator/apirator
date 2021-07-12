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

package openapi

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

type Validator struct {
	loader *openapi3.Loader
}

func NewValidator(loader *openapi3.Loader) *Validator {
	return &Validator{loader: loader}
}

func (v *Validator) Validate(definition string) error {
	data, err := v.loader.LoadFromData([]byte(definition))
	if err != nil {
		return err
	}
	err = data.Validate(context.Background())
	if err != nil {
		return &InvalidDefinitionError{message: err.Error()}
	}
	return nil
}
