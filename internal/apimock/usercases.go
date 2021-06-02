package apimock

import "github.com/apirator/apirator/internal/usecase"

type UserCases struct {
	*usecase.ConfigMap
	*usecase.ValidOpenAPI
	*usecase.Deployment
	*usecase.Ingress
	*usecase.Service
}
