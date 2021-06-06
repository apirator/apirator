package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/api/v1alpha1/phase"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/openapi"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
)

type OpenAPIDefinition struct {
	*openapi.Validator
	*k8s.Service
}

func (v *OpenAPIDefinition) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	err := v.Validate(apimock.Spec.Definition)
	if err != nil {
		if apimock.SetConditionForError(err) {
			log.Info("invalid OpenAPI definition", "cause", err)
			apimock.Status.Phase = phase.Error
			return operation.RequeueOnErrorOrStop(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.StopProcessing()
	}

	if apimock.SetConditionForValidOpenAPI() {
		log.Info("valid OpenAPI definition")
		apimock.Status.Phase = phase.Pending
		return operation.RequeueOnErrorOrStop(v.UpdateAPIMockStatus(ctx, apimock))
	}

	return operation.ContinueProcessing()
}
