//+build wireinject

package main

import (
	"github.com/apirator/apirator/controllers"
	"github.com/google/wire"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func newAPIMockReconciler(mgr manager.Manager) (*controllers.APIMockReconciler, error) {
	wire.Build(providers)
	return nil, nil
}
