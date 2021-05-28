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
	svc     *Service
	Log     logr.Logger
}

func newAdapter(APIMock *api.APIMock, svc *Service) *Adapter {
	return &Adapter{
		APIMock: APIMock,
		Log:     ctrl.Log.WithName("adapters").WithName("APIMock"),
		svc:     svc,
	}
}
