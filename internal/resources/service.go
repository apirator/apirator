package resources

import (
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	svcPortName  = "http"
	docsPortName = "docs"
)

func (b *Builder) ServiceFor(resource *v1alpha1.APIMock) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Labels:    v1alpha1.APIMockLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: v1alpha1.APIMockLabels,
			Type:     resource.Spec.ServiceDefinition.ServiceType,
			Ports: []corev1.ServicePort{
				{
					Name:       svcPortName,
					Protocol:   "TCP",
					Port:       int32(resource.Spec.ServiceDefinition.Port),
					TargetPort: intstr.FromInt(8000),
				},
				{
					Name:       docsPortName,
					Protocol:   "TCP",
					Port:       int32(8080),
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(resource, svc, b.scheme); err != nil {
		return nil, fmt.Errorf("failed to set Service %q owner reference: %v", resource.GetName(), err)
	}

	return svc, nil
}
