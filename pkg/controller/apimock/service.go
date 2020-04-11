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

package apimock

import (
	"context"

	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/apirator/apirator/pkg/controller/k8s/util/labels"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	svcPortName  = "http"
	docsPortName = "docs"
)

func (r *ReconcileAPIMock) EnsureService(mock *v1alpha1.APIMock) error {
	svcPresent := mock.Spec.ServiceDefinition.Port != 0
	if svcPresent {
		svcK8s := &v1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      mock.GetName(),
			Namespace: mock.Namespace,
		}, svcK8s)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Service not found. Starting creation...", "Service.Namespace", mock.Namespace, "Service.Name", mock.Name)
			mockPort := v1.ServicePort{
				Name:       svcPortName,
				Protocol:   "TCP",
				Port:       int32(mock.Spec.ServiceDefinition.Port),
				TargetPort: intstr.FromInt(8000),
			}
			docPort := v1.ServicePort{
				Name:       docsPortName,
				Protocol:   "TCP",
				Port:       int32(8080),
				TargetPort: intstr.FromInt(8000),
			}
			svc := &v1.Service{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            mock.GetName(),
					Namespace:       mock.GetNamespace(),
					OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(mock, mock.GroupVersionKind())},
				},
				Spec: v1.ServiceSpec{
					Selector: labels.LabelForAPIMock(mock),
					Type:     mock.Spec.ServiceDefinition.ServiceType,
					Ports:    []v1.ServicePort{mockPort, docPort},
				},
			}
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
		} else {
			if mock.AnnotateClusterIP(svcK8s.Spec.ClusterIP) || mock.AnnotatePorts(svcK8s.Spec.Ports) {
				err := r.client.Update(context.TODO(), mock)
				if err != nil {
					log.Error(err, "Failed to update APIMock annotations", "Service.Namespace", svcK8s.Namespace, "Service.Name", svcK8s.Name, "Service.ClusterIP", mock.Annotations["apirator.io/cluster-ip"], "Service.Ports", mock.Annotations["apirator.io/ports"])
					return err
				}
				log.Info("APIMock annotations update successfully", "Service.Namespace", svcK8s.Namespace, "Service.Name", svcK8s.Name, "Service.ClusterIP", mock.Annotations["apirator.io/cluster-ip"], "Service.Ports", mock.Annotations["apirator.io/ports"])
			}
		}
	} else {
		log.Info("Service is not necessary")
	}
	return nil
}
