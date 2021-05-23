package controllers

import (
	"github.com/apirator/apirator/internal/mock"
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

var Providers = wire.NewSet(
	NewAPIMockReconciler,
)

func NewAPIMockReconciler(svc *mock.Service, scheme *runtime.Scheme) *APIMockReconciler {
	return &APIMockReconciler{
		Service: svc,
		Log:     ctrl.Log.WithName("controllers").WithName("APIMock"),
		Scheme:  scheme,
	}
}
