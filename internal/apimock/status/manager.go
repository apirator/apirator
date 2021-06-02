package status

import (
	"github.com/apirator/apirator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Manager struct{}

func (m *Manager) SetCondition(apimock *v1alpha1.APIMock, conditionType v1alpha1.APIMockConditionType, status corev1.ConditionStatus, reason, message string) {
	now := metav1.Now()
	condition, _ := m.findConditionOrInitialize(apimock, conditionType)
	if message != condition.Message || status != condition.Status || reason != condition.Reason || conditionType != condition.Type {
		condition.LastTransitionTime = now
	}
	if message != "" {
		condition.Message = message
	}
	condition.LastProbeTime = now
	condition.Reason = reason
	condition.Status = status
}

func (m *Manager) HasCondition(apimock *v1alpha1.APIMock, conditionType v1alpha1.APIMockConditionType) bool {
	conditions := apimock.Status.Conditions
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return true
		}
	}
	return false
}

func (m *Manager) findConditionOrInitialize(apimock *v1alpha1.APIMock, conditionType v1alpha1.APIMockConditionType) (*v1alpha1.APIMockCondition, bool) {
	if apimock.Status.Conditions == nil {
		apimock.Status.Conditions = make([]*v1alpha1.APIMockCondition, 0)
	}
	conditions := apimock.Status.Conditions
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return conditions[i], true
		}
	}
	condition := &v1alpha1.APIMockCondition{Type: conditionType}
	conditions = append(conditions, condition)
	return condition, false
}
