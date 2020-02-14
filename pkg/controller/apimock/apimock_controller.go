package apimock

import (
	"context"
	"github.com/apirator/apirator/pkg/controller/oas"

	apiratorv1alpha1 "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
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

	err = c.Watch(&source.Kind{Type: &apiratorv1alpha1.APIMock{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
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
}

func (r *ReconcileAPIMock) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling APIMock")

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

	errOas := oas.Validate(instance.Spec.Definition)
	if errOas != nil {
		if err := r.markAsInvalidOAS(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	cfmErr := r.EnsureConfigMap(instance)
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

	svcErr := r.EnsureService(instance)
	if svcErr != nil {
		if err := r.markAsFailure(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err := r.markAsSuccessful(instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
