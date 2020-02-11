package apimock

import (
	"context"
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/apirator/apirator/pkg/controller/k8s/util/labels"
	"github.com/apirator/apirator/pkg/controller/k8s/util/owner"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	svcPortName = "http"
)

func (r *ReconcileAPIMock) EnsureService(mock *v1alpha1.APIMock) error {
	svcPresent := mock.Spec.Port != 0
	if svcPresent {
		svcK8s := &v1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      mock.GetName(),
			Namespace: mock.Namespace,
		}, svcK8s)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Service not found. Starting creation...", "Service.Namespace", mock.Namespace, "Service.Name", mock.Name)
			svcPort := v1.ServicePort{
				Name:       svcPortName,
				Protocol:   "tcp",
				Port:       int32(mock.Spec.Port),
				TargetPort: intstr.FromInt(mock.Spec.ContainerPort),
			}
			svc := &v1.Service{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      mock.GetName(),
					Namespace: mock.GetNamespace(),
				},
				Spec: v1.ServiceSpec{
					Selector: labels.LabelForAPIMock(mock),
					Type:     v1.ServiceTypeLoadBalancer,
					Ports:    []v1.ServicePort{svcPort},
				},
			}
			owner.AddOwnerRefToObject(svc, owner.AsOwner(&mock.ObjectMeta))
			err := r.client.Create(context.TODO(), svc)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
				return err
			}
			log.Info("Service created successfully", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return err
		}
	} else {
		log.Info("Service is not necessary")
	}
	return nil
}
