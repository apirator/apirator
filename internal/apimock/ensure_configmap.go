package apimock

import (
	"context"
	"fmt"

	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
)

func (a *Adapter) EnsureConfigMap(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := a.Log.WithValues("trace", span.String())
	defer span.Finish()

	cm, err := a.svc.SearchConfigMap(ctx, a.APIMock)
	if err != nil {
		return nil, err
	}
	if cm == nil {
		cm, err = a.svc.CreateConfigMap(ctx, a.APIMock)
		if err != nil {
			return nil, err
		}
		log.Info(fmt.Sprintf("the ConfigMap %q was created", cm.GetName()))
	}

	return operation.ContinueProcessing()
}
