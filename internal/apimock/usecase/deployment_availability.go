package usecase

import (
	"context"
	"time"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
)

type DeploymentAvailability struct {
	*k8s.Service
}

func (d *DeploymentAvailability) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := span.Logger()

	status, err := d.GetDeploymentStatus(ctx, apimock)
	if err != nil {
		return nil, err
	}

	if status.HasAvailableCondition() {
		if apimock.SetAvailableConditionTrue() {
			log.Info("mock deployment has minimum availability")
			return operation.RequeueOnErrorOrStop(d.UpdateAPIMockStatus(ctx, apimock))
		}
		return operation.ContinueProcessing()
	}

	if apimock.SetAvailableConditionFalse() {
		log.Info("mock deployment has no minimum availability")
		return operation.RequeueOnErrorOrStop(d.UpdateAPIMockStatus(ctx, apimock))
	}

	log.Info("waiting for mock deployment availability")
	return operation.RequeueAfter(10*time.Second, nil)
}
