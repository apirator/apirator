package main

import (
	ctrl "github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/apimock"
	"github.com/apirator/apirator/internal/k8s"
	"github.com/apirator/apirator/internal/manager"
	"github.com/google/wire"
)

var providers = wire.NewSet(
	apimock.Providers,
	ctrl.Providers,
	k8s.Providers,
	manager.Providers,
)
