package adapter

import "github.com/apirator/apirator/internal/apimock/usecase"

type UserCases struct {
	*usecase.ConfigMap
	*usecase.Deployment
	*usecase.Ingress
	*usecase.InitializedStatus
	*usecase.OpenAPIDefinition
	*usecase.Service
}
