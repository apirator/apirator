package resources

import (
	"fmt"
	"path/filepath"

	"github.com/apirator/apirator/api/v1alpha1"
	yu "github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	yamlConfigPath = "/etc/oas/oas.yaml"
	jsonConfigPath = "/etc/oas/oas.json"
)

func (b *Builder) ConfigMapFor(resource *v1alpha1.APIMock) (*corev1.ConfigMap, error) {
	bJson, err := yu.YAMLToJSON([]byte(resource.Spec.Definition))
	if err != nil {
		return nil, fmt.Errorf("failed to convert openapi definition to JSON: %w", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    v1alpha1.APIMockLabels,
		},
		Data: map[string]string{
			filepath.Base(yamlConfigPath): resource.Spec.Definition,
			filepath.Base(jsonConfigPath): string(bJson),
		},
	}

	if err := controllerutil.SetControllerReference(resource, cm, b.scheme); err != nil {
		return nil, fmt.Errorf("failed to set ConfigMap %q owner reference: %v", cm.GetName(), err)
	}

	return cm, nil
}
