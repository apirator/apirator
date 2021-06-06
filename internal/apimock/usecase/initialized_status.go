package usecase

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/api/v1alpha1/phase"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InitializedStatus struct {
	*k8s.Service
}

func (v *InitializedStatus) Ensure(ctx context.Context, apimock *v1alpha1.APIMock) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if apimock.Status.Conditions == nil {
		apimock.Status.Conditions = make([]metav1.Condition, 0, 1)
		apimock.Status.Phase = phase.Pending
		return operation.RequeueOnErrorOrStop(v.UpdateAPIMockStatus(ctx, apimock))
	}
	return operation.ContinueProcessing()
}
