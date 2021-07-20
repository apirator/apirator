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

	"github.com/apirator/apirator/api/v1alpha1"
	"github.com/apirator/apirator/internal/reconcile"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// APIMockReconciler reconciles a APIMock object
type APIMockReconciler struct {
	APIMockReader
	APIMockAdapter
}

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apirator.io,resources=apimocks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apirator.io,resources=apimocks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apirator.io,resources=apimocks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *APIMockReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	apimock, err := r.GetAPIMock(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.DoNotRequeue()
		}
		return reconcile.RequeueOnErr(err)
	}

	return reconcile.NewHandler(
		r.EnsureStatus,
		r.EnsureDefinitionIsValid,
		r.EnsureConfigMap,
		r.EnsureDeployment,
		r.EnsureService,
		r.EnsureIngress,
		r.EnsureIngressFinalizer,
		r.EnsureDeploymentAvailability,
	).Handle(ctx, apimock)
}

// SetupWithManager sets up the controller with the Manager.
func (r *APIMockReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.APIMock{}).
		Complete(r)
}
