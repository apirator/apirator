package apimock

import (
	"context"
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *ReconcileAPIMock) EnsureIngress(mock *v1alpha1.APIMock, doc *openapi3.Swagger) error {
	if mock.ExposeInIngress() {
		log.Info("Ingress Should be configured.Starting...", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)
		ingressK8s := &v1beta1.Ingress{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      mock.Spec.Selector[v1alpha1.IngressTag],
			Namespace: mock.Spec.Selector[v1alpha1.NamespaceTag],
		}, ingressK8s)
		if err != nil && errors.IsNotFound(err) {
			error := &v1alpha1.IngressControllerNotFound{Name: mock.Spec.Selector[v1alpha1.IngressTag]}
			log.Error(error, "Ingress-Controller not found. ", "Service.Namespace", mock.Spec.Selector[v1alpha1.IngressTag], "Service.Name", mock.Spec.Selector[v1alpha1.NamespaceTag])
			return error
		}

		if checkHostConflict(ingressK8s, mock) {
			error := &v1alpha1.HostAlreadyExist{
				Host:    mock.Spec.Host,
				Ingress: mock.Spec.Selector[v1alpha1.IngressTag],
			}
			log.Error(error, "Host already configured. ", "Service.Namespace", mock.Spec.Selector[v1alpha1.IngressTag], "Service.Name", mock.Spec.Selector[v1alpha1.NamespaceTag])
			return error
		}

		// append new rule in ingress controller. SUCCESS
		ingressK8s.Spec.Rules = append(ingressK8s.Spec.Rules, newRule(mock, doc))

		// update ingress
		updateIngErr := r.client.Update(context.TODO(), ingressK8s)
		if updateIngErr != nil {
			log.Error(err, "Failed to update ingress", "Ingress.Namespace", ingressK8s.GetNamespace(), "Ingress.Name", ingressK8s.GetName())
			return err
		}
		log.Info("Ingress configured successfully", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)
	}
	return nil
}

// check if host was previously already configured
func checkHostConflict(ingress *v1beta1.Ingress, mock *v1alpha1.APIMock) bool {
	for _, r := range ingress.Spec.Rules {
		if mock.Spec.Host == r.Host {
			return true
		}
	}
	return false
}

// create new rule from mock
func newRule(mock *v1alpha1.APIMock, doc *openapi3.Swagger) v1beta1.IngressRule {
	var paths []v1beta1.HTTPIngressPath
	for path := range doc.Paths {
		inP := v1beta1.HTTPIngressPath{
			Path: path,
			Backend: v1beta1.IngressBackend{
				ServiceName: mock.GetName(),
				ServicePort: intstr.FromInt(mock.Spec.ServiceDefinition.Port),
			},
		}
		paths = append(paths, inP)
	}
	igv := &v1beta1.HTTPIngressRuleValue{Paths: paths}
	irv := v1beta1.IngressRuleValue{HTTP: igv}
	return v1beta1.IngressRule{
		Host:             mock.Spec.Host,
		IngressRuleValue: irv,
	}
}
