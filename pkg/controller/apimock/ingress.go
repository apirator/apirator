package apimock

import (
	"context"
	"encoding/json"
	"github.com/apirator/apirator/internal/steps"
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
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

		log.Info("Checking ingress entries...", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)

		// only check for the new mock
		if !mock.CheckStep(steps.IngressEntryCreated) {
			log.Info("Start adding new entry in ingress. Checking hosts...", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)
			if checkHostConflict(ingressK8s, mock) {
				error := &v1alpha1.HostAlreadyExist{
					Host:    mock.Spec.Host,
					Ingress: mock.Spec.Selector[v1alpha1.IngressTag],
				}
				log.Error(error, "Host already configured. ", "Service.Namespace", mock.Spec.Selector[v1alpha1.IngressTag], "Service.Name", mock.Spec.Selector[v1alpha1.NamespaceTag])
				return error
			}
			// append new rule in ingress controller. SUCCESS
			log.Info("Adding new entry in ingress. Hosts checked everything is ok.", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)
			ingressK8s.Spec.Rules = append(ingressK8s.Spec.Rules, newRule(mock, doc))
		} else {
			log.Info("Ingress was configured previously", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)
		}

		// update ingress
		updateIngErr := r.client.Update(context.TODO(), ingressK8s)
		if updateIngErr != nil {
			log.Error(err, "[Adding] - Failed to update ingress", "Ingress.Namespace", ingressK8s.GetNamespace(), "Ingress.Name", ingressK8s.GetName())
			return err
		} else {
			log.Info("Ingress configured successfully", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)
		}
		mock.AddStep(steps.NewIngressEntryCreated())
	}
	return nil
}

// remove apimock entry from ingress
func (r *ReconcileAPIMock) RemoveEntryFromIngress(mock *v1alpha1.APIMock) error {
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
	ingressK8s.Spec.Rules = removeEntry(ingressK8s.Spec.Rules, mock.Spec.Host)
	updateIngErr := r.client.Update(context.TODO(), ingressK8s)
	if updateIngErr != nil {
		log.Error(err, "[Removing] - Failed to update ingress", "Ingress.Namespace", ingressK8s.GetNamespace(), "Ingress.Name", ingressK8s.GetName())
		return err
	}
	log.Info("Entry removed from ingress successfully", "Mock.Namespace", mock.Namespace, "Mock.Name", mock.Name)
	return nil
}

// remove entry by host
func removeEntry(rules []v1beta1.IngressRule, host string) []v1beta1.IngressRule {
	var elIdx = 0
	for idx, r := range rules {
		if host == r.Host {
			elIdx = idx
			break
		}
	}
	return append(rules[:elIdx], rules[elIdx+1:]...)
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
	inP := v1beta1.HTTPIngressPath{
		Path: path(doc),
		Backend: v1beta1.IngressBackend{
			ServiceName: mock.GetName(),
			ServicePort: intstr.FromInt(mock.Spec.ServiceDefinition.Port),
		},
	}
	paths = append(paths, inP)
	igv := &v1beta1.HTTPIngressRuleValue{Paths: paths}
	irv := v1beta1.IngressRuleValue{HTTP: igv}
	return v1beta1.IngressRule{
		Host:             mock.Spec.Host,
		IngressRuleValue: irv,
	}
}

// find mock api path
// https://swagger.io/docs/specification/openapi-extensions/
func path(doc *openapi3.Swagger) string {
	i := doc.Info.Extensions["x-apirator-mock-path"]
	json := i.(json.RawMessage)
	path := strings.Trim(string(json), "\"")
	log.Info("Path", "Mock.Path", path)
	return path
}
