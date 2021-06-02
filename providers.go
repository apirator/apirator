package main

import (
	"github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/apimock/adapter"
	"github.com/apirator/apirator/internal/apimock/status"
	"github.com/apirator/apirator/internal/apimock/usecase"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/resources"
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

	adapter.NewFactory,
	controllers.NewAPIMockReconciler,
	k8s.NewService,
	resources.NewBuilder,
	usecase.Providers,
	wire.Bind(new(controllers.AdapterFactory), new(*adapter.Factory)),
	wire.Struct(new(adapter.UserCases), "*"),
	wire.Struct(new(status.Manager), "*"),
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
