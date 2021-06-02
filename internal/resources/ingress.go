package resources

import (
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) IngressFor(resource *v1alpha1.APIMock) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       networkingv1.IngressSpec{},
	}

	if err := controllerutil.SetControllerReference(resource, ing, b.scheme); err != nil {
		return nil, fmt.Errorf("failed to set Service %q owner reference: %v", resource.GetName(), err)
	}

	return ing, nil
}
