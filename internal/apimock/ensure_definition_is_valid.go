package apimock

import (
	"context"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	"github.com/getkin/kin-openapi/openapi3"
)

func (a *Adapter) EnsureDefinitionIsValid(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := a.Log.WithValues("trace", span.String())
	defer span.Finish()

	loader := &openapi3.Loader{Context: ctx}
	doc, err := loader.LoadFromData([]byte(a.Spec.Definition))
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	err = doc.Validate(ctx)
	if err != nil {
		span.SetError(err)
		log.Info("openapi definition is invalid", "cause", err)

		// TODO: mark as invalid oas
		return operation.StopProcessing()
	}

	return operation.ContinueProcessing()
}
