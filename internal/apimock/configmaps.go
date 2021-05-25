package apimock

import (
	"context"
	"fmt"
	"path/filepath"

	api "github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/tracing"
	yu "github.com/ghodss/yaml"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	yamlConfigPath = "/etc/oas/oas.yaml"
	jsonConfigPath = "/etc/oas/oas.json"
)

func (s *Service) SearchConfigMap(ctx context.Context, apimock *api.APIMock) (*core.ConfigMap, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	cm := &core.ConfigMap{}
	err := s.client.Get(ctx, types.NamespacedName{Name: apimock.Name, Namespace: apimock.Namespace}, cm)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		span.SetError(err)
		return nil, fmt.Errorf("failed to lookup ConfigMap: %w", err)
	}
	isControlledBy := metav1.IsControlledBy(cm, apimock)
	if !isControlledBy {
		return nil, fmt.Errorf("conflicting ConfigMap")
	}
	return cm, nil
}

func (s *Service) CreateConfigMap(ctx context.Context, apimock *api.APIMock) (*core.ConfigMap, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	bJson, err := yu.YAMLToJSON([]byte(apimock.Spec.Definition))
	if err != nil {
		span.SetError(err)
		return nil, fmt.Errorf("failed to convert openapi definition to JSON: %w", err)
	}

	cm := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: apimock.Name, Namespace: apimock.Namespace},
		Data: map[string]string{
			filepath.Base(yamlConfigPath): apimock.Spec.Definition,
			filepath.Base(jsonConfigPath): string(bJson),
		},
	}
	if err := controllerutil.SetControllerReference(apimock, cm, s.scheme); err != nil {
		span.SetError(err)
		return nil, fmt.Errorf("failed to set ConfigMap %q owner reference: %v", cm.GetName(), err)
	}

	err = s.client.Create(ctx, cm)
	if err != nil {
		span.SetError(err)
		return nil, fmt.Errorf("failed to create ConfigMap: %w", err)
	}

	return cm, nil
}
