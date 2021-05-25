package apimock

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewService,
)
