package adapter

import (
	"context"

	"github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Factory struct {
	userCases *UserCases
	svc       *k8s.Service
}

func NewFactory(userCases *UserCases, svc *k8s.Service) *Factory {
	return &Factory{userCases: userCases, svc: svc}
}

func (a *Factory) CreateAPIMockAdapter(ctx context.Context, key client.ObjectKey) (controllers.APIMockAdapter, error) {
	resource, err := a.svc.GetAPIMock(ctx, key)
	if err != nil {
		return nil, err
	}
	return newAdapter(a.userCases, resource), err
}
