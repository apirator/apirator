package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type DeploymentStatus struct{ *appsv1.DeploymentStatus }

func (d *DeploymentStatus) HasAvailableCondition() bool {
	for _, condition := range d.Conditions {
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
