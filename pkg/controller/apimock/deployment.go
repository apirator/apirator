package apimock

import (
	"context"
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/apirator/apirator/pkg/controller/k8s/util/labels"
	"github.com/apirator/apirator/pkg/controller/k8s/util/owner"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *ReconcileAPIMock) EnsureDeployment(mock *v1alpha1.APIMock) error {
	svcK8s := &v1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      mock.GetName(),
		Namespace: mock.Namespace,
	}, svcK8s)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Deployment not found. Starting creation...", "Deployment.Namespace", mock.Namespace, "Deployment.Name", mock.Name)
		var reps int32
		reps = int32(1)
		pt := BuildPodTemplate(mock)
		d := &v1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      mock.GetName(),
				Namespace: mock.GetNamespace(),
				Labels:    labels.LabelForAPIMock(mock),
			},
			Spec: v1.DeploymentSpec{
				Replicas: &reps,
				Selector: &metav1.LabelSelector{MatchLabels: labels.LabelForAPIMock(mock)},
				Template: pt,
				Strategy: v1.DeploymentStrategy{
					Type: v1.RollingUpdateDeploymentStrategyType,
					RollingUpdate: &v1.RollingUpdateDeployment{
						MaxUnavailable: func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
						MaxSurge:       func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
					},
				},
			},
		}
		ref := owner.AsOwner(&mock.ObjectMeta)
		d.SetOwnerReferences([]metav1.OwnerReference{ref})
		err := r.client.Create(context.TODO(), d)
		if err != nil {
			log.Error(err, "Failed to create Deployment")
			return err
		}
		log.Info("Deployment created successfully", "Deployment.Namespace", d.Namespace, "Deployment.Name", d.Name)
		return nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return err
	}
	return nil
}
