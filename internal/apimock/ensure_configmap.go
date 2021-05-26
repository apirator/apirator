package apimock

import (
	"context"
	"fmt"
	api "github.com/apirator/apirator/api/v1alpha1"
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

	desired, err := newDesiredConfigMap(a.APIMock)
	if err != nil {
		span.SetError(err)
		return nil, err
	}
	if err := controllerutil.SetControllerReference(a.APIMock, desired, a.svc.scheme); err != nil {
		span.SetError(err)
		return nil, fmt.Errorf("failed to set ConfigMap %q owner reference: %v", desired.GetName(), err)
	}

	opts := []client.ListOption{
		client.InNamespace(a.APIMock.Namespace),
		client.MatchingLabels(map[string]string{"app.kubernetes.io/managed-by": "apirator"}),
	}
	list := &core.ConfigMapList{}
	if err := a.svc.client.List(ctx, list, opts...); err != nil {
		span.SetError(err)
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}

	inv := inventory.ForConfigMaps(list.Items, []core.ConfigMap{*desired})
	err = a.createConfigMaps(ctx, inv)
	if err != nil {
		return nil, err
	}

	err = a.updateConfigMap(ctx, inv)
	if err != nil {
		return nil, err
	}

	err = a.deleteConfigMap(ctx, inv)
	if err != nil {
		return nil, err
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) createConfigMaps(ctx context.Context, inv inventory.ConfigMap) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := a.Log.WithValues("trace", span.String())
	defer span.Finish()

	for _, cm := range inv.Create {
		log.Info(fmt.Sprintf("creating ConfigMap %q", cm.GetName()))
		if err := a.svc.client.Create(ctx, &cm); err != nil {
			span.SetError(err)
			return fmt.Errorf("failed to create ConfigMap %q: %w", cm.GetName(), err)
		}
	}

	return nil
}

func (a *Adapter) updateConfigMap(ctx context.Context, inv inventory.ConfigMap) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := a.Log.WithValues("trace", span.String())
	defer span.Finish()

	for _, cm := range inv.Update {
		log.Info(fmt.Sprintf("updating ConfigMap %q", cm.GetName()))
		if err := a.svc.client.Update(ctx, &cm); err != nil {
			span.SetError(err)
			return fmt.Errorf("failed to update ConfigMap %q: %w", cm.GetName(), err)
		}
	}

	return nil
}

func (a *Adapter) deleteConfigMap(ctx context.Context, inv inventory.ConfigMap) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	log := a.Log.WithValues("trace", span.String())
	defer span.Finish()

	for _, cm := range inv.Delete {
		log.Info(fmt.Sprintf("deleting ConfigMap %q", cm.GetName()))
		if err := a.svc.client.Delete(ctx, &cm); err != nil {
			span.SetError(err)
			return fmt.Errorf("failed to delete ConfigMap %q: %w", cm.GetName(), err)
		}
	}

	return nil
}

func newDesiredConfigMap(apimock *api.APIMock) (*core.ConfigMap, error) {
	bJson, err := yu.YAMLToJSON([]byte(apimock.Spec.Definition))
	if err != nil {
		return nil, fmt.Errorf("failed to convert openapi definition to JSON: %w", err)
	}

	return &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apimock.Name,
			Namespace: apimock.Namespace,
			Labels:    map[string]string{"app.kubernetes.io/managed-by": "apirator"},
		},
		Data: map[string]string{
			filepath.Base(yamlConfigPath): apimock.Spec.Definition,
			filepath.Base(jsonConfigPath): string(bJson),
		},
	}, nil
}
