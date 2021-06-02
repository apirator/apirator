package adapter

import "github.com/apirator/apirator/internal/apimock/usecase"

type UserCases struct {
	*usecase.ConfigMap
	*usecase.OpenAPIDefinition
	*usecase.Deployment
	*usecase.Ingress
	*usecase.Service
}
