package main

import (
	ctrl "github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/manager"
	"github.com/google/wire"
)

var providers = wire.NewSet(
	ctrl.Providers,
	k8s.Providers,
	manager.Providers,
)
