//go:build wireinject
// +build wireinject

package controllers

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/apirator/apirator/internal/apimock"
	internalclient "github.com/apirator/apirator/internal/client"
	"github.com/apirator/apirator/internal/objects"
	"github.com/apirator/apirator/internal/openapi"
	"github.com/apirator/apirator/internal/usecase"
)

func newAPIMockReconciler(clientClient client.Client, scheme *runtime.Scheme) (*APIMockReconciler, error) {
	wire.Build(
		objects.NewBuilder,
		openapi.NewValidator,
		openapi3.NewLoader,
		usecase.NewConfigMap,
		usecase.NewDeployment,
		usecase.NewDeploymentAvailability,
		usecase.NewIngress,
		usecase.NewIngressFinalizer,
		usecase.NewOpenAPIDefinition,
		usecase.NewService,
		usecase.NewStatus,
		wire.Bind(new(APIMockAdapter), new(*apimock.Adapter)),
		wire.Bind(new(APIMockReader), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.APIMockReader), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.APIMockStatusWriter), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.APIMockWriter), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.ConfigMapBuilder), new(*objects.Builder)),
		wire.Bind(new(usecase.ConfigMapReader), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.DeploymentBuilder), new(*objects.Builder)),
		wire.Bind(new(usecase.DeploymentReader), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.DeploymentStatusReader), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.GenericObjectWriter), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.IngressBuilder), new(*objects.Builder)),
		wire.Bind(new(usecase.IngressReader), new(*internalclient.Kubernetes)),
		wire.Bind(new(usecase.OpenAPIValidator), new(*openapi.Validator)),
		wire.Bind(new(usecase.ServiceBuilder), new(*objects.Builder)),
		wire.Bind(new(usecase.ServiceReader), new(*internalclient.Kubernetes)),
		wire.Struct(new(apimock.Adapter), "*"),
		wire.Struct(new(APIMockReconciler), "*"),
		wire.Struct(new(internalclient.Kubernetes), "*"),
	)
	return nil, nil
}
