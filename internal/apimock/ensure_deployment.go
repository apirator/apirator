package apimock

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/apirator/apirator/internal/inventory"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/apirator/apirator/internal/operation"
	"github.com/apirator/apirator/internal/tracing"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (a *Adapter) EnsureDeployment(ctx context.Context) (*operation.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	desired, err := a.newDesiredDeployment()
	if err != nil {
		return nil, span.HandleError(err)
	}

	list, err := a.listDeployments()
	if err != nil {
		return nil, span.HandleError(err)
	}

	inv := inventory.ForDeployments(list.Items, []appsv1.Deployment{*desired})
	err = a.svc.Apply(ctx, inv)
	if err != nil {
		return nil, span.HandleError(err)
	}

	return operation.ContinueProcessing()
}

func (a *Adapter) listDeployments() (*appsv1.DeploymentList, error) {
	opts := []client.ListOption{
		client.InNamespace(a.resource.Namespace),
		client.MatchingLabels(Labels),
	}
	list := new(appsv1.DeploymentList)
	if err := a.svc.client.List(context.TODO(), list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Deployments: %w", err)
	}
	return list, nil
}

func (a *Adapter) newDesiredDeployment() (*appsv1.Deployment, error) {
	reps := int32(1)
	labels := Labels

	volumes := a.newVolumes()
	containers := []v1.Container{a.newMockContainer(), a.newDocContainer()}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.resource.GetName(),
			Namespace: a.resource.GetNamespace(),
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &reps,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      a.resource.GetName(),
					Namespace: a.resource.GetNamespace(),
					Labels:    labels,
				},
				Spec: v1.PodSpec{
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

	if err := controllerutil.SetControllerReference(a.resource, dep, a.scheme); err != nil {
		return nil, fmt.Errorf("failed to set Deployment %q owner reference: %v", dep.GetName(), err)
	}

	return dep, nil
}

func (a *Adapter) newMockContainer() v1.Container {
	var ports []v1.ContainerPort
	ports = append(ports, v1.ContainerPort{
		ContainerPort: mockPort,
		Name:          mockPortName,
	})
	cnPort := v1.EnvVar{
		Name:  "PORT",
		Value: strconv.Itoa(mockPort),
	}
	cnWatch := v1.EnvVar{
		Name:  "WATCH",
		Value: strconv.FormatBool(a.resource.Spec.Watch),
	}

	// Handler for probes
	rh := v1.Handler{
		HTTPGet: &v1.HTTPGetAction{
			Path:   "/__health",
			Port:   intstr.FromInt(8000),
			Scheme: "HTTP",
		},
	}

	// LivenessProbe and Readiness Probe
	rp := &v1.Probe{
		Handler:             rh,
		InitialDelaySeconds: 2,
		TimeoutSeconds:      1,
		PeriodSeconds:       3,
	}

	return v1.Container{
		Name:    mockContainerName,
		Image:   mockImageName,
		Command: []string{"apisprout"},
		Args: []string{
			mockVolumeMountPath + "oas.yaml",
		},
		VolumeMounts:   a.newContainerMounts(),
		Ports:          ports,
		Resources:      a.newContainerRequirements(),
		Env:            []v1.EnvVar{cnPort, cnWatch},
		ReadinessProbe: rp,
		LivenessProbe:  rp,
	}
}

func (a *Adapter) newDocContainer() v1.Container {
	var ports []v1.ContainerPort
	ports = append(ports, v1.ContainerPort{
		ContainerPort: docPort,
		Name:          docPortName,
	})
	cnPort := v1.EnvVar{
		Name:  "PORT",
		Value: strconv.Itoa(docPort),
	}
	oasPath := v1.EnvVar{
		Name:  "SWAGGER_JSON",
		Value: "/etc/oas/oas.json",
	}
	baseUrl := v1.EnvVar{
		Name:  "BASE_URL",
		Value: "/" + a.resource.GetName() + "/docs",
	}
	return v1.Container{
		Name:         docContainerName,
		Image:        docImageName,
		VolumeMounts: a.newContainerMounts(),
		Ports:        ports,
		Resources:    a.newContainerRequirements(),
		Env:          []v1.EnvVar{cnPort, oasPath, baseUrl},
	}
}

func (a *Adapter) newContainerMounts() []v1.VolumeMount {
	return []v1.VolumeMount{{
		Name:      mockVolumeMountName,
		MountPath: filepath.Dir(mockVolumeMountPath),
	}}
}

func (a *Adapter) newContainerRequirements() v1.ResourceRequirements {
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

func (a *Adapter) newVolumes() []v1.Volume {
	return []v1.Volume{{
		Name: mockVolumeMountName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: a.resource.GetName(),
				},
			},
		},
	}}
}
