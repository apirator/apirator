// Copyright 2020 apirator.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
