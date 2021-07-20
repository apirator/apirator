/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"github.com/apirator/apirator/internal/apimock"
	internalclient "github.com/apirator/apirator/internal/client"
	"github.com/apirator/apirator/internal/objects"
	"github.com/apirator/apirator/internal/openapi"
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SetupWithManager(mgr ctrl.Manager) error {
	r, err := newAPIMockReconciler(mgr.GetClient(), mgr.GetScheme())
	if err != nil {
		return err
	}
	return r.SetupWithManager(mgr)
}

func newAPIMockReconciler(clientClient client.Client, scheme *runtime.Scheme) (*APIMockReconciler, error) {
	k8s := &internalclient.Kubernetes{Client: clientClient}
	builder := objects.NewBuilder(scheme)
	loader := openapi3.NewLoader()
	validator := openapi.NewValidator(loader)
	adapter := apimock.NewAdapter(builder, k8s, k8s, builder, k8s, k8s, k8s, builder, k8s, k8s, k8s, validator, builder, k8s)
	return &APIMockReconciler{APIMockReader: k8s, APIMockAdapter: adapter}, nil
}
