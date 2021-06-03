package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
)

type InitializedStatus struct {
	*k8s.Service
}

func (v *InitializedStatus) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if apimock.UpdateStatus() {
		return operation.RequeueOnErrorOrStop(v.UpdateAPIMockStatus(ctx, apimock))
	}
	return operation.ContinueProcessing()
}
