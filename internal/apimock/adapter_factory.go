package apimock

import (
	"context"

	"github.com/apirator/apirator/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AdapterFactory struct {
	scheme *runtime.Scheme
	svc    *Service
}

func NewAdapterFactory(scheme *runtime.Scheme, svc *Service) *AdapterFactory {
	return &AdapterFactory{scheme: scheme, svc: svc}
}

func (a *AdapterFactory) CreateAPIMockAdapter(ctx context.Context, key client.ObjectKey) (controllers.APIMockAdapter, error) {
	svc := a.svc
	scheme := a.scheme
	resource, err := svc.LookupAPIMock(ctx, key)
	if err != nil {
		return nil, err
	}
	return newAdapter(scheme, svc, resource), err
}
