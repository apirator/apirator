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

package objects

import (
	"github.com/apirator/apirator/api/v1alpha1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (b *Builder) IngressesFor(apimocks *v1alpha1.APIMockList) (*networkingv1beta1.IngressList, error) {
	filtered := make([]v1alpha1.APIMock, 0)
	for _, r := range apimocks.Items {
		exposed := r.Spec.Ingress != nil
		deleted := r.GetDeletionTimestamp() != nil
		if exposed && !deleted {
			filtered = append(filtered, r)
		}
	}
	if len(filtered) == 0 {
		return &networkingv1beta1.IngressList{}, nil
	}
	ing := &networkingv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "apirator",
			Labels:      map[string]string{"app.kubernetes.io/managed-by": "apirator"},
			Annotations: nil,
		},
		Spec: networkingv1beta1.IngressSpec{
			TLS:   newIngressTLS(filtered),
			Rules: newIngressRules(filtered),
		},
	}

	for _, r := range filtered {
		if r.Spec.Ingress != nil {
			if len(r.Spec.Ingress.Annotations) > 0 {
				ing.Annotations = merge(r.Spec.Ingress.Annotations, ing.Annotations)
			}
			ing.Namespace = r.GetNamespace()
			if err := b.SetOwnerReference(&r, ing); err != nil {
				return nil, err
			}
		}
	}

	return &networkingv1beta1.IngressList{Items: []networkingv1beta1.Ingress{*ing}}, nil
}

func newIngressRules(apimocks []v1alpha1.APIMock) []networkingv1beta1.IngressRule {
	hosts := mapHosts(apimocks)
	rules := make([]networkingv1beta1.IngressRule, 0, len(hosts))
	for host, backends := range hosts {
		rules = append(rules, networkingv1beta1.IngressRule{
			Host: host,
			IngressRuleValue: networkingv1beta1.IngressRuleValue{
				HTTP: &networkingv1beta1.HTTPIngressRuleValue{Paths: backends},
			},
		})
	}
	if len(rules) > 0 {
		return rules
	}
	return nil
}

func newIngressTLS(apimocks []v1alpha1.APIMock) []networkingv1beta1.IngressTLS {
	secrets := mapTLSSecrets(apimocks)
	tls := make([]networkingv1beta1.IngressTLS, 0, len(secrets))
	for secret, hosts := range secrets {
		tls = append(tls, networkingv1beta1.IngressTLS{
			Hosts:      hosts,
			SecretName: secret,
		})
	}
	if len(tls) > 0 {
		return tls
	}
	return nil
}

func newHTTPIngressPath(path, service string, port int) networkingv1beta1.HTTPIngressPath {
	prefix := networkingv1beta1.PathTypePrefix
	return networkingv1beta1.HTTPIngressPath{
		Path:     path,
		PathType: &prefix,
		Backend: networkingv1beta1.IngressBackend{
			ServiceName: service,
			ServicePort: intstr.FromInt(port),
		},
	}
}

func mapTLSSecrets(apimocks []v1alpha1.APIMock) map[string][]string {
	tlsSecrets := make(map[string][]string, 0)
	for _, r := range apimocks {
		if r.Spec.Ingress != nil && r.Spec.Ingress.TLS != nil {
			secretName := r.Spec.Ingress.TLS.SecretName
			tlsSecrets[secretName] = dedupe(tlsSecrets[secretName], r.Spec.Ingress.Hostname)
		}
	}
	return tlsSecrets
}

func mapHosts(apimocks []v1alpha1.APIMock) map[string][]networkingv1beta1.HTTPIngressPath {
	hosts := make(map[string][]networkingv1beta1.HTTPIngressPath, 0)
	for _, r := range apimocks {
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
