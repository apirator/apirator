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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	svcPortName  = "http"
	docsPortName = "docs"
)

func (b *Builder) ServiceFor(apimock *v1alpha1.APIMock) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        apimock.GetName(),
			Namespace:   apimock.GetNamespace(),
			Labels:      apimock.MatchLabels(),
			Annotations: apimock.Spec.Service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector:              apimock.MatchLabels(),
			Type:                  apimock.Spec.Service.Type,
			ExternalTrafficPolicy: apimock.Spec.Service.ExternalTrafficPolicy,
			Ports: []corev1.ServicePort{
				{
					Name:       svcPortName,
					Protocol:   "TCP",
					Port:       int32(apimock.Spec.Service.Port),
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

	if err := b.SetControllerReference(apimock, svc); err != nil {
		return nil, err
	}

	return svc, nil
}
