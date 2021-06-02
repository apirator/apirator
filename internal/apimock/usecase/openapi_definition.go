package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/apimock/status"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	"github.com/getkin/kin-openapi/openapi3"
	corev1 "k8s.io/api/core/v1"
)

type OpenAPIDefinition struct {
	*status.Manager
	*k8s.Service
}

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
		log.Info("invalid OpenAPI definition", "cause", err)
		v.SetCondition(apimock, v1alpha1.APIMockValidOpenAPIDefinition, corev1.ConditionFalse, "invalid OpenAPI definition", err.Error())
		return operation.RequeueOnErrorOrStop(v.UpdateAPIMockStatus(ctx, apimock))
	}

	log.Info("valid OpenAPI definition")
	v.SetCondition(apimock, v1alpha1.APIMockValidOpenAPIDefinition, corev1.ConditionTrue, "valid OpenAPI definition", "")
	return operation.RequeueOnErrorOrStop(v.UpdateAPIMockStatus(ctx, apimock))
}
