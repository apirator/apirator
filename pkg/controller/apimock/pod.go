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
	docPortName         = "doc-port"
	mockImageName       = "danielgtaylor/apisprout"
	docImageName        = "apirator/mockdoc"
)

// it will create the pod template, with doc-container and mock-container
func BuildPodTemplate(mock *v1alpha1.APIMock) v1.PodTemplateSpec {
	volumes := podVolume(mock)
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mock.GetName(),
			Namespace: mock.GetNamespace(),
			Labels:    labels.LabelForAPIMock(mock),
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{mockContainer(mock), docContainer(mock)},
			Volumes:    volumes,
		},
	}
}

// create mock container, it will deploy the mock api
func mockContainer(mock *v1alpha1.APIMock) v1.Container {
	var ports []v1.ContainerPort
	if mock.Spec.MockContainerPort != 0 {
		ports = append(ports, v1.ContainerPort{
			ContainerPort: int32(mock.Spec.MockContainerPort),
			Name:          mockPortName,
		})
	}
	return v1.Container{
		Name:    mock.GetName(),
		Image:   mockImageName,
		Command: []string{"apisprout"},
		Args: []string{
			mockVolumeMountPath + "/" + "oas.yaml",
		},
		VolumeMounts: volumeMount(),
		Ports:        ports,
		Resources:    requirements(),
	}

}

// create documentation container, it will used for display the swagger-ui
func docContainer(mock *v1alpha1.APIMock) v1.Container {
	var ports []v1.ContainerPort
	if mock.Spec.DocContainerPort != 0 {
		ports = append(ports, v1.ContainerPort{
			ContainerPort: int32(mock.Spec.DocContainerPort),
			Name:          docPortName,
		})
	}
	return v1.Container{
		Name:         mock.GetName(),
		Image:        docImageName,
		VolumeMounts: volumeMount(),
		Ports:        ports,
		Resources:    requirements(),
	}

}

// configure the volume mount for containers
func volumeMount() []v1.VolumeMount {
	return []v1.VolumeMount{{
		Name:      mockVolumeMountName,
		MountPath: filepath.Dir(mockVolumeMountPath),
	}}
}

// configure minimum requirements
func requirements() (resourcesReq v1.ResourceRequirements) {
	// Configure Requests
	requests := v1.ResourceList{}
	requests[v1.ResourceCPU] = resource.MustParse("10m")
	requests[v1.ResourceMemory] = resource.MustParse("5Mi")
	// Configure Limits
	limits := v1.ResourceList{}
	limits[v1.ResourceCPU] = resource.MustParse("20m")
	limits[v1.ResourceMemory] = resource.MustParse("10Mi")
	return v1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

// configure the pod volumes
func podVolume(mock *v1alpha1.APIMock) []v1.Volume {
	return []v1.Volume{{
		Name: mockVolumeMountName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: mock.GetName(),
				},
			},
		},
	}}
}
