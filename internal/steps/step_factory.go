package steps

import (
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// create a deployment created step
func NewIngressEntryCreated() v1alpha1.Step {
	return newStep(v1alpha1.IngressEntryCreated, "Ingress entry created successfully")
}

// create a deployment created step
func NewDeploymentCreated() v1alpha1.Step {
	return newStep(v1alpha1.DeploymentCreated, "Deployment created successfully")
}

// create a service created step
func NewServiceCreated() v1alpha1.Step {
	return newStep(v1alpha1.ServiceCreated, "Service created successfully")
}

// create a config map created step
func NewConfigMapCreated() v1alpha1.Step {
	return newStep(v1alpha1.CfgMapCreated, "Config map created successfully")
}

func newStep(action string, description string) v1alpha1.Step {
	return v1alpha1.Step{
		Action:      action,
		Description: description,
		LastUpdate:  metav1.Now(),
	}
}
