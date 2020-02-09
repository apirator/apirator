package apimock

import (
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/apirator/apirator/pkg/controller/k8s/util/labels"
	"github.com/apirator/apirator/pkg/controller/k8s/util/owner"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	svcPortName = "mock-svc-port"
)

func service(mock *v1alpha1.APIMock) error {

	svcPresent := mock.Spec.Port != 0

	svcPort := v1.ServicePort{
		Name:       svcPortName,
		Protocol:   "http",
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

	if svcPresent {

	}

	return nil
}
