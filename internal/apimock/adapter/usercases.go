package adapter

import "github.com/apirator/apirator/internal/apimock/usecase"

type UserCases struct {
	*usecase.ConfigMap
	*usecase.Deployment
	*usecase.DeploymentAvailability
	*usecase.Ingress
	*usecase.Status
	*usecase.OpenAPIDefinition
	*usecase.Service
}
