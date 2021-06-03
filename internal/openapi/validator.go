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
