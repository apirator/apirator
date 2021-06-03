package usecase

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	wire.Struct(new(ConfigMap), "*"),
	wire.Struct(new(Deployment), "*"),
	wire.Struct(new(Ingress), "*"),
	wire.Struct(new(InitializedStatus), "*"),
	wire.Struct(new(OpenAPIDefinition), "*"),
	wire.Struct(new(Service), "*"),
)
