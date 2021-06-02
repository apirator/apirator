package apimock

import (
	"context"

	"github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AdapterFactory struct {
	userCases *UserCases
	svc       *k8s.Service
}

func NewAdapterFactory(userCases *UserCases, svc *k8s.Service) *AdapterFactory {
	return &AdapterFactory{userCases: userCases, svc: svc}
}

func (a *AdapterFactory) CreateAPIMockAdapter(ctx context.Context, key client.ObjectKey) (controllers.APIMockAdapter, error) {
	resource, err := a.svc.GetAPIMock(ctx, key)
	if err != nil {
		return nil, err
	}
	return newAdapter(a.userCases, resource), err
}
