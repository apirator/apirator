package apimock

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/apirator/apirator/internal/inventory"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	yu "github.com/ghodss/yaml"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (a *Adapter) EnsureConfigMap(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := a.newDesiredConfigMap()
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := a.listConfigMaps()
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForConfigMaps(list.Items, []core.ConfigMap{*desired})
	err = a.svc.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) listConfigMaps() (*core.ConfigMapList, error) {
	opts := []client.ListOption{
		client.InNamespace(a.resource.Namespace),
		client.MatchingLabels(Labels),
	}
	list := new(core.ConfigMapList)
	if err := a.svc.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}
	return list, nil
}

func (a *Adapter) newDesiredConfigMap() (*core.ConfigMap, error) {
	bJson, err := yu.YAMLToJSON([]byte(a.resource.Spec.Definition))
	if err != nil {
		return nil, fmt.Errorf("failed to convert openapi definition to JSON: %w", err)
	}

	cm := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.resource.Name,
			Namespace: a.resource.Namespace,
			Labels:    Labels,
		},
		Data: map[string]string{
			filepath.Base(yamlConfigPath): a.resource.Spec.Definition,
			filepath.Base(jsonConfigPath): string(bJson),
		},
	}

	if err := controllerutil.SetControllerReference(a.resource, cm, a.scheme); err != nil {
		return nil, fmt.Errorf("failed to set ConfigMap %q owner reference: %v", cm.GetName(), err)
	}

	return cm, nil
}
