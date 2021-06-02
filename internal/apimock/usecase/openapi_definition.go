package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIDefinition struct{}

func (v *OpenAPIDefinition) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	loader := &openapi3.Loader{Context: ctx}
	doc, err := loader.LoadFromData([]byte(apimock.Spec.Definition))
	if err != nil {
		return nil, span.HandleError(err)
	}

	err = doc.Validate(ctx)
	if err != nil {
		log.Info("openapi definition is invalid", "cause", err)

		// TODO: mark as invalid oas
		return operation.StopProcessing()
	}

	log.Info("openapi definition is valid")
	return operation.ContinueProcessing()
}
