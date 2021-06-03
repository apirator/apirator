package main

import (
	"github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/apimock/adapter"
	"github.com/apirator/apirator/internal/apimock/usecase"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/openapi"
	"github.com/apirator/apirator/internal/resources"
	"github.com/getkin/kin-openapi/openapi3"
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
	openapi.NewValidator,
	openapi3.NewLoader,
	resources.NewBuilder,
	usecase.Providers,
	wire.Bind(new(controllers.AdapterFactory), new(*adapter.Factory)),
	wire.Struct(new(adapter.UserCases), "*"),
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
