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

package oas

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func Validate(definition string) error {
	doc := &openapi3.Swagger{}
	err := yaml.Unmarshal([]byte(definition), doc)
	if err != nil {
		log.Log.Error(err, "Error to parse yaml to oas")
	}
	oasErr := doc.Validate(context.TODO())
	if oasErr != nil {
		log.Log.Error(oasErr, "Open API Specification is invalid")
	}
	log.Log.Info("Open API Specification is VALID")
	return nil
}
