package apimock

import (
	api "github.com/apirator/apirator/api/v1alpha1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	yamlConfigPath = "/etc/oas/oas.yaml"
	jsonConfigPath = "/etc/oas/oas.json"
)

type Adapter struct {
	resource *api.APIMock
	logger   logr.Logger
	svc      *Service
}

func newAdapter(resource *api.APIMock, svc *Service) *Adapter {
	return &Adapter{
		resource: resource,
		logger:   ctrl.Log.WithName("adapters").WithName("APIMock"),
		svc:      svc,
	}
}
