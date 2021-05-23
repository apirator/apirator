package controllers

import (
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var Providers = wire.NewSet(
	NewAPIMockReconciler,
)

func NewAPIMockReconciler(client client.Client, scheme *runtime.Scheme) *APIMockReconciler {
	return &APIMockReconciler{
		Client: client,
		Log:    ctrl.Log.WithName("controllers").WithName("APIMock"),
		Scheme: scheme,
	}
}
