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