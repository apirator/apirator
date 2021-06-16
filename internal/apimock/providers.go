package apimock

import (
	"github.com/apirator/apirator/internal/apimock/adapter"
	"github.com/apirator/apirator/internal/apimock/usecase"
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	adapter.NewFactory,
	wire.Struct(new(adapter.UserCases), "*"),

	wire.Struct(new(usecase.ConfigMap), "*"),
	wire.Struct(new(usecase.Deployment), "*"),
	wire.Struct(new(usecase.DeploymentAvailability), "*"),
	wire.Struct(new(usecase.Ingress), "*"),
	wire.Struct(new(usecase.Status), "*"),
	wire.Struct(new(usecase.OpenAPIDefinition), "*"),
	wire.Struct(new(usecase.Service), "*"),
)
