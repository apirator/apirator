// Copyright 2020 apirator.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apimock

import (
	"context"
	"github.com/apirator/apirator/pkg/controller/oas"
	resources "github.com/apirator/apirator/pkg/controller/predicates"

	apiratorv1alpha1 "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/redhat-cop/operator-utils/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_apimock")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAPIMock{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("apimock-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &apiratorv1alpha1.APIMock{}}, &handler.EnqueueRequestForObject{}, resources.StatusChangedPredicate{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &apiratorv1alpha1.APIMock{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAPIMock{}

type ReconcileAPIMock struct {
	client client.Client
	scheme *runtime.Scheme
	util.ReconcilerBase
}

func (r *ReconcileAPIMock) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling APIMock")

	// Get the updated instance
	instance := &apiratorv1alpha1.APIMock{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	reqLogger.Info("Starting Reconcile Rules...", "APIMock.IsInitialized", instance.Spec.Initialized)

	// initializing instance
	if ok := instance.IsInitialized(); !ok {
		reqLogger.Info("Initializing APIMock...", "APIMock.Name", instance.GetName())
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "unable to update APIMock", "mock.name", instance.GetName())
			return r.ManageError(instance, err)
		}
		reqLogger.Info("APIMock initialized successfully...", "APIMock.IsInitialized", instance.Spec.Initialized)
		return reconcile.Result{}, nil
	}

	// In deletion process (handling exclusion)
	if util.IsBeingDeleted(instance) {
		reqLogger.Info("Deleting APIMock...", "APIMock.Name", instance.GetName())
		// Everything is ok, there is nothing else to do
		if instance.HasFinalizer(apiratorv1alpha1.IngressFinalizerName) && instance.ExposeInIngress() {
			reqLogger.Info("Executing cleanup logic...", "APIMock.Name", instance.GetName())
			err := r.manageCleanUpLogic(instance)
			if err != nil {
				log.Error(err, "unable to delete apimock", "mock.name", instance.GetName())
				return r.ManageError(instance, err)
			}
			reqLogger.Info("Cleanup executed successfully", "APIMock.Name", instance.GetName())
			reqLogger.Info("Removing finalizers...", "APIMock.Name", instance.GetName())
			instance.RemoveFinalizer(apiratorv1alpha1.IngressFinalizerName)
			updErr := r.client.Update(context.TODO(), instance)
			if updErr != nil {
				reqLogger.Error(err, "unable to update mock", "mock.name", instance.GetName())
				return r.ManageError(instance, err)
			}
			reqLogger.Info("Finalizers removed successfully", "APIMock.Name", instance.GetName())
			return reconcile.Result{}, nil
		} else {
			reqLogger.Info("There is nothing to do. Success", "APIMock.Name", instance.GetName())
			return reconcile.Result{}, nil
		}
	}

	reqLogger.Info("Starting configuration logic...", "APIMock.IsInitialized", instance.Spec.Initialized)
	doc, errOas := oas.Validate(instance.Spec.Definition)
	if errOas != nil {
		reqLogger.Error(errOas, "Open API Specification is invalid")
		if err := r.markAsInvalidOAS(instance); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, errOas
	}

	updatedCfgMap, cfmErr := r.EnsureConfigMap(instance)
	if cfmErr != nil {
		if err := r.markAsFailure(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	depErr := r.EnsureDeployment(instance)
	if depErr != nil {
		if err := r.markAsFailure(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	updatedSvc, svcErr := r.EnsureService(instance)
	if svcErr != nil {
		if err := r.markAsFailure(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	ingErr := r.EnsureIngress(instance, doc)
	if ingErr != nil {
		if err := r.markAsFailure(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Should update by cfgmap", "ConfigMap", updatedCfgMap)
	reqLogger.Info("Should update by SVC", "Service", updatedSvc)

	// update the final status
	if err := r.markAsSuccessful(instance); err != nil {
		return reconcile.Result{}, err
	}

	// updated the instance
	if updatedCfgMap || updatedSvc {
		reqLogger.Info("Updating mock. Something is different...", "APIMock.Name", instance.GetName())
		if err := r.client.Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "Error on update mock")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
