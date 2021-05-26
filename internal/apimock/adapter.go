package apimock

import (
	"fmt"
	"path/filepath"

	api "github.com/apirator/apirator/api/v1alpha1"
	yu "github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

func (a *Adapter) DesiredConfigMap() (*core.ConfigMap, error) {
	bJson, err := yu.YAMLToJSON([]byte(a.APIMock.Spec.Definition))
	if err != nil {
		return nil, fmt.Errorf("failed to convert openapi definition to JSON: %w", err)
	}

	cm := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: a.APIMock.Name, Namespace: a.APIMock.Namespace},
		Data: map[string]string{
			filepath.Base(yamlConfigPath): a.APIMock.Spec.Definition,
			filepath.Base(jsonConfigPath): string(bJson),
		},
	}
	if err := controllerutil.SetControllerReference(a.APIMock, cm, a.svc.scheme); err != nil {
		return nil, fmt.Errorf("failed to set ConfigMap %q owner reference: %v", cm.GetName(), err)
	}

	return cm, nil
}
