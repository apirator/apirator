package adapter

import (
	"context"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/operation"
)

type Adapter struct {
	*UserCases

	customresource *v1alpha1.APIMock
}

func newAdapter(userCases *UserCases, customresource *v1alpha1.APIMock) *Adapter {
	return &Adapter{UserCases: userCases, customresource: customresource}
}

func (a *Adapter) EnsureConfigMap(ctx context.Context) (*operation.Result, error) {
	return a.ConfigMap.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureDefinitionIsValid(ctx context.Context) (*operation.Result, error) {
	return a.OpenAPIDefinition.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureDeployment(ctx context.Context) (*operation.Result, error) {
	return a.Deployment.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureIngress(ctx context.Context) (*operation.Result, error) {
	return a.Ingress.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureService(ctx context.Context) (*operation.Result, error) {
	return a.Service.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureIsInitialized(ctx context.Context) (*operation.Result, error) {
	return a.InitializedStatus.Ensure(ctx, a.customresource)
}

func (a *Adapter) EnsureFinalizer(ctx context.Context) (*operation.Result, error) {
	panic("implement me")
}

func (a *Adapter) EnsureProvisionedStatus(ctx context.Context) (*operation.Result, error) {
	panic("implement me")
}
