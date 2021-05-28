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
	APIMock *api.APIMock
	logger  logr.Logger
	svc     *Service
}

func newAdapter(APIMock *api.APIMock, svc *Service) *Adapter {
	return &Adapter{
		APIMock: APIMock,
		logger:  ctrl.Log.WithName("adapters").WithName("APIMock"),
		svc:     svc,
	}
}
