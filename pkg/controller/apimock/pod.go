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
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/apirator/apirator/pkg/controller/k8s/util/labels"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
)

const (
	mockVolumeMountName = "oas"
	mockVolumeMountPath = "/etc/oas/"
	mockPortName        = "mock-port"
	mockImageName       = "danielgtaylor/apisprout"
)

func BuildPodTemplate(mock *v1alpha1.APIMock) v1.PodTemplateSpec {
	volumes := []v1.Volume{{
		Name: mockVolumeMountName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: mock.GetName(),
				},
			},
		},
	}}
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mock.GetName(),
			Namespace: mock.GetNamespace(),
			Labels:    labels.LabelForAPIMock(mock),
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{mockContainer(mock)},
			Volumes:    volumes,
		},
	}
}

func mockContainer(mock *v1alpha1.APIMock) v1.Container {
	vm := []v1.VolumeMount{{
		Name:      mockVolumeMountName,
		MountPath: filepath.Dir(mockVolumeMountPath),
	}}
	var ports []v1.ContainerPort
	if mock.Spec.ContainerPort != 0 {
		ports = append(ports, v1.ContainerPort{
			ContainerPort: int32(mock.Spec.ContainerPort),
			Name:          mockPortName,
		})
	}

	// Configure Requests
	requests := v1.ResourceList{}
	requests[v1.ResourceCPU] = resource.MustParse("10m")
	requests[v1.ResourceMemory] = resource.MustParse("5Mi")

	// Configure Limits
	limits := v1.ResourceList{}
	limits[v1.ResourceCPU] = resource.MustParse("20m")
	limits[v1.ResourceMemory] = resource.MustParse("10Mi")

	requirements := v1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
	return v1.Container{
		Name:    mock.GetName(),
		Image:   mockImageName,
		Command: []string{"apisprout"},
		Args: []string{
			mockVolumeMountPath + "/" + "oas.yaml",
		},
		VolumeMounts: vm,
		Ports:        ports,
		Resources:    requirements,
	}

}
