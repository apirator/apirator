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

package resources

import (
	"fmt"

	"github.com/apirator/apirator/api/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *Builder) IngressFor(resources *v1alpha1.APIMockList) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "apirator",
			Labels: map[string]string{"app.kubernetes.io/managed-by": "apirator"},
		},
		Spec: networkingv1.IngressSpec{
			TLS:   newIngressTLS(resources),
			Rules: newIngressRules(resources),
		},
	}

	for _, r := range resources.Items {
		if r.Spec.Ingress != nil {
			ing.Annotations = merge(r.Spec.Ingress.Annotations, ing.Annotations)
			ing.Namespace = r.GetNamespace()
			if err := controllerutil.SetOwnerReference(&r, ing, b.scheme); err != nil {
				return nil, fmt.Errorf("failed to set Service %q owner reference: %v", r.GetName(), err)
			}
		}
	}

	return ing, nil
}

func newIngressRules(resources *v1alpha1.APIMockList) []networkingv1.IngressRule {
	hosts := mapHosts(resources)
	tls := make([]networkingv1.IngressRule, 0, len(hosts))
	for host, backends := range hosts {
		tls = append(tls, networkingv1.IngressRule{
			Host: host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{Paths: backends},
			},
		})
	}
	return tls
}

func newIngressTLS(resources *v1alpha1.APIMockList) []networkingv1.IngressTLS {
	secrets := mapTLSSecrets(resources)
	tls := make([]networkingv1.IngressTLS, 0, len(secrets))
	for secret, hosts := range secrets {
		tls = append(tls, networkingv1.IngressTLS{
			Hosts:      hosts,
			SecretName: secret,
		})
	}
	return tls
}

func newHTTPIngressPath(path, service string, port int) networkingv1.HTTPIngressPath {
	prefix := networkingv1.PathTypePrefix
	return networkingv1.HTTPIngressPath{
		Path:     path,
		PathType: &prefix,
		Backend: networkingv1.IngressBackend{Service: &networkingv1.IngressServiceBackend{
			Name: service,
			Port: networkingv1.ServiceBackendPort{Number: int32(port)},
		}},
	}
}

func mapTLSSecrets(resources *v1alpha1.APIMockList) map[string][]string {
	tlsSecrets := make(map[string][]string, 0)
	for _, r := range resources.Items {
		if r.Spec.Ingress != nil && r.Spec.Ingress.TLS != nil {
			secretName := r.Spec.Ingress.TLS.SecretName
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.Spec.Ingress.Hostname)
		}
	}
	return tlsSecrets
}

func mapHosts(resources *v1alpha1.APIMockList) map[string][]networkingv1.HTTPIngressPath {
	hosts := make(map[string][]networkingv1.HTTPIngressPath, 0)
	for _, r := range resources.Items {
		if r.Spec.Ingress != nil {
			hosts[r.Spec.Ingress.Hostname] = append(
				hosts[r.Spec.Ingress.Hostname],
				newHTTPIngressPath(r.Spec.Ingress.Path, r.GetName(), r.Spec.Service.Port),
			)
		}
	}
	return hosts
}

func merge(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func dedupe(a []string, b ...string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}
	for letter := range check {
		res = append(res, letter)
	}
	return res
}
