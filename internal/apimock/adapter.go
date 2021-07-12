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

package apimock

import (
	"github.com/apirator/apirator/internal/usecase"
)

type Adapter struct {
	*usecase.ConfigMap
	*usecase.Deployment
	*usecase.DeploymentAvailability
	*usecase.Ingress
	*usecase.IngressFinalizer
	*usecase.OpenAPIDefinition
	*usecase.Service
	*usecase.Status
}

func NewAdapter(
	cmBuilder usecase.ConfigMapBuilder, cmReader usecase.ConfigMapReader, goWriter usecase.GenericObjectWriter,
	depBuilder usecase.DeploymentBuilder, depReader usecase.DeploymentReader, depStatusReader usecase.DeploymentStatusReader,
	apiStatusWriter usecase.APIMockStatusWriter, ingBuilder usecase.IngressBuilder, ingReader usecase.IngressReader,
	apiReader usecase.APIMockReader, apiWriter usecase.APIMockWriter, apiValidator usecase.OpenAPIValidator,
	svcBuilder usecase.ServiceBuilder, svcReader usecase.ServiceReader) *Adapter {
	return &Adapter{
		ConfigMap:              usecase.NewConfigMap(cmBuilder, cmReader, goWriter),
		Deployment:             usecase.NewDeployment(depBuilder, depReader, goWriter),
		DeploymentAvailability: usecase.NewDeploymentAvailability(depStatusReader, apiStatusWriter),
		Ingress:                usecase.NewIngress(ingBuilder, ingReader, apiReader, goWriter),
		IngressFinalizer:       usecase.NewIngressFinalizer(apiWriter),
		OpenAPIDefinition:      usecase.NewOpenAPIDefinition(apiValidator, apiStatusWriter),
		Service:                usecase.NewService(svcBuilder, svcReader, goWriter),
		Status:                 usecase.NewStatus(apiStatusWriter),
	}
}
