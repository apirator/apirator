package usecase

import (
	"context"
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/api/v1alpha1/phase"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
)

type Status struct {
	*k8s.Service
}

func (v *Status) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	if apimock.Status.Phase == "" {
		apimock.Status.Phase = phase.Pending
		log.Info(fmt.Sprintf("Updating status to %q", phase.Pending))
		return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
	}

	if apimock.IsValidatedConditionFalse() {
		if apimock.Status.Phase != phase.Error {
			apimock.Status.Phase = phase.Error
			log.Info(fmt.Sprintf("Updating status to %q", phase.Error))
			return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.ContinueProcessing()
	}

	if apimock.IsAvailableConditionFalse() {
		if apimock.Status.Phase != phase.Pending {
			apimock.Status.Phase = phase.Pending
			log.Info(fmt.Sprintf("Updating status to %q", phase.Pending))
			return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.ContinueProcessing()
	}

	if apimock.IsValidatedConditionTrue() && apimock.IsAvailableConditionTrue() {
		if apimock.Status.Phase != phase.Running {
			apimock.Status.Phase = phase.Running
			log.Info(fmt.Sprintf("Updating status to %q", phase.Running))
			return operation.RequeueWithError(v.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.ContinueProcessing()
	}

	return operation.ContinueProcessing()
}
