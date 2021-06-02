package status

import (
	"github.com/apirator/apirator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Manager struct {
}

func (m *Manager) SetCondition(conditions []v1alpha1.APIMockCondition, conditionType v1alpha1.APIMockConditionType, status corev1.ConditionStatus, reason string, message string) {
	now := metav1.Now()
	condition, _ := m.FindCondition(conditions, conditionType)
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

func (m *Manager) FindCondition(conditions []v1alpha1.APIMockCondition, conditionType v1alpha1.APIMockConditionType) (*v1alpha1.APIMockCondition, bool) {
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return &conditions[i], true
		}
	}
	conditions = append(conditions, v1alpha1.APIMockCondition{Type: conditionType})
	return &conditions[len(conditions)-1], false
}

func (m *Manager) HasCondition(conditions []v1alpha1.APIMockCondition, conditionType v1alpha1.APIMockConditionType) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return true
		}
	}
	return false
}
