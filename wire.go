//+build wireinject

package main

import (
	"github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/apimock"
	"github.com/apirator/apirator/internal/apimock/adapter"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/openapi"
	"github.com/apirator/apirator/internal/resources"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func newAPIMockReconciler(mgr manager.Manager) (*controllers.APIMockReconciler, error) {
	wire.Build(
		extractClient,
		extractScheme,

		apimock.Providers,
		controllers.NewAPIMockReconciler,
		k8s.NewService,
		openapi.NewValidator,
		openapi3.NewLoader,
		resources.NewBuilder,
		wire.Bind(new(controllers.AdapterFactory), new(*adapter.Factory)))
	return nil, nil
}

func extractScheme(mgr manager.Manager) *runtime.Scheme {
	return mgr.GetScheme()
}

func extractClient(mgr manager.Manager) client.Client {
	return mgr.GetClient()
}
