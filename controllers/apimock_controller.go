/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// APIMockReconciler reconciles a APIMock object
type APIMockReconciler struct {
	factory AdapterFactory
	logger  logr.Logger
}

func NewAPIMockReconciler(factory AdapterFactory) *APIMockReconciler {
	return &APIMockReconciler{
		factory: factory,
		logger:  ctrl.Log.WithName("controllers").WithName("APIMock"),
	}
}

//+kubebuilder:rbac:groups=apimocks.apirator.io,resources=apimocks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apimocks.apirator.io,resources=apimocks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apimocks.apirator.io,resources=apimocks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the APIMock object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *APIMockReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, tracing.WithCustomResource(req.NamespacedName))
	defer span.Finish()

	log := r.logger.WithValues("trace", span.String())
	log.Info("reconciling")

	adapter, err := r.factory.CreateAPIMockAdapter(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return r.doNotRequeue()
		}
		return r.requeueOnErr(err)
	}

	result, err := r.handle(ctx,
		adapter.EnsureDefinitionIsValid,
		adapter.EnsureConfigMap,
		adapter.EnsureDeployment,
		adapter.EnsureService,
	)

	log.V(1).
		WithValues("error", err != nil, "requeing", result.Requeue, "delay", result.RequeueAfter).
		Info("finished reconcile")
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *APIMockReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.APIMock{}).
		Complete(r)
}

func (r *APIMockReconciler) handle(ctx context.Context, operations ...operation.Func) (reconcile.Result, error) {
	for _, op := range operations {
		result, err := op(ctx)
		if err != nil && result == nil {
			return r.requeueOnErr(err)
		}
		if err != nil || (result != nil && result.RequeueRequest) {
			return r.requeueAfter(result.RequeueDelay, err)
		}
		if result.CancelRequest {
			return r.doNotRequeue()
		}
	}
	return r.doNotRequeue()
}

func (r *APIMockReconciler) doNotRequeue() (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (r *APIMockReconciler) requeueOnErr(err error) (reconcile.Result, error) {
	return reconcile.Result{}, err
}

func (r *APIMockReconciler) requeueAfter(duration time.Duration, err error) (reconcile.Result, error) {
	return reconcile.Result{RequeueAfter: duration}, err
}
