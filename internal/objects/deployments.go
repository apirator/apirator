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
	"path/filepath"
	"strconv"

	"github.com/apirator/apirator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	docContainerName = "doc"
	docImageName     = "swaggerapi/swagger-ui:v3.25.0"
	docPort          = 8080
	docPortName      = "doc-port"

	mockContainerName   = "mock"
	mockImageName       = "apirator/mock"
	mockPort            = 8000
	mockPortName        = "mock-port"
	mockVolumeMountName = "oas"
	mockVolumeMountPath = "/etc/oas/"
)

func (b *Builder) DeploymentFor(apimock *v1alpha1.APIMock) (*appsv1.Deployment, error) {
	reps := int32(1)
	labels := apimock.MatchLabels()

	volumes := newVolumes(apimock)
	containers := []corev1.Container{newMockContainer(apimock), newDocContainer(apimock)}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apimock.GetName(),
			Namespace: apimock.GetNamespace(),
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &reps,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      apimock.GetName(),
					Namespace: apimock.GetNamespace(),
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: containers,
					Volumes:    volumes,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
					MaxSurge:       func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
				},
			},
		},
	}

	if err := b.SetControllerReference(apimock, dep); err != nil {
		return nil, err
	}

	return dep, nil
}

func newMockContainer(apimock *v1alpha1.APIMock) corev1.Container {
	var ports []corev1.ContainerPort
	ports = append(ports, corev1.ContainerPort{
		ContainerPort: mockPort,
		Name:          mockPortName,
		Protocol:      corev1.ProtocolTCP,
	})
	cnPort := corev1.EnvVar{
		Name:  "PORT",
		Value: strconv.Itoa(mockPort),
	}
	cnWatch := corev1.EnvVar{
		Name:  "WATCH",
		Value: strconv.FormatBool(apimock.Spec.Watch),
	}

	// Handler for probes
	rh := corev1.Handler{
		HTTPGet: &corev1.HTTPGetAction{
			Path:   "/__health",
			Port:   intstr.FromInt(8000),
			Scheme: "HTTP",
		},
	}

	// LivenessProbe and Readiness Probe
	rp := &corev1.Probe{
		Handler:             rh,
		InitialDelaySeconds: 2,
		TimeoutSeconds:      1,
		PeriodSeconds:       3,
	}

	return corev1.Container{
		Name:    mockContainerName,
		Image:   mockImageName,
		Command: []string{"apisprout"},
		Args: []string{
			mockVolumeMountPath + "oas.yaml",
		},
		VolumeMounts:   newContainerMounts(),
		Ports:          ports,
		Resources:      newContainerRequirements(),
		Env:            []corev1.EnvVar{cnPort, cnWatch},
		ReadinessProbe: rp,
		LivenessProbe:  rp,
	}
}

func newDocContainer(apimock *v1alpha1.APIMock) corev1.Container {
	var ports []corev1.ContainerPort
	ports = append(ports, corev1.ContainerPort{
		ContainerPort: docPort,
		Name:          docPortName,
		Protocol:      corev1.ProtocolTCP,
	})
	cnPort := corev1.EnvVar{
		Name:  "PORT",
		Value: strconv.Itoa(docPort),
	}
	oasPath := corev1.EnvVar{
		Name:  "SWAGGER_JSON",
		Value: "/etc/oas/oas.json",
	}
	baseUrl := corev1.EnvVar{
		Name:  "BASE_URL",
		Value: "/" + apimock.GetName() + "/docs",
	}
	return corev1.Container{
		Name:         docContainerName,
		Image:        docImageName,
		VolumeMounts: newContainerMounts(),
		Ports:        ports,
		Resources:    newContainerRequirements(),
		Env:          []corev1.EnvVar{cnPort, oasPath, baseUrl},
	}
}

func newContainerMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{{
		Name:      mockVolumeMountName,
		MountPath: filepath.Dir(mockVolumeMountPath),
	}}
}

func newContainerRequirements() corev1.ResourceRequirements {
	// Configure Requests
	requests := corev1.ResourceList{}
	requests[corev1.ResourceCPU] = resource.MustParse("10m")
	requests[corev1.ResourceMemory] = resource.MustParse("5Mi")
	// Configure Limits
	limits := corev1.ResourceList{}
	limits[corev1.ResourceCPU] = resource.MustParse("20m")
	limits[corev1.ResourceMemory] = resource.MustParse("10Mi")
	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

func newVolumes(apimock *v1alpha1.APIMock) []corev1.Volume {
	return []corev1.Volume{{
		Name: mockVolumeMountName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: apimock.GetName(),
				},
			},
		},
	}}
}
