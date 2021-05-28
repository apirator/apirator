package main

import (
	"github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/apimock"
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var providers = wire.NewSet(
	extractClient,
	extractRestConfig,
	extractScheme,

	apimock.NewAdapterFactory,
	apimock.NewService,
	controllers.NewAPIMockReconciler,
	wire.Bind(new(controllers.AdapterFactory), new(*apimock.AdapterFactory)),
)

func extractScheme(mgr manager.Manager) *runtime.Scheme {
	return mgr.GetScheme()
}

func extractClient(mgr manager.Manager) client.Client {
	return mgr.GetClient()
}

func extractRestConfig(mgr manager.Manager) *rest.Config {
	return mgr.GetConfig()
}
