package condition

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// These are valid conditions of a deployment APIMock.
const (
	Provisioned string = "Provisioned"
	Error              = "Error"
	Validated          = "Validated"
	Waiting            = "Waiting"
)

func NewValidOpenAPIDefinition(status metav1.ConditionStatus, message string) metav1.Condition {
	return metav1.Condition{
		Type:    Validated,
		Status:  status,
		Reason:  "ValidOpenAPIDefinition",
		Message: message,
	}
}

func NewError(status metav1.ConditionStatus, message string) metav1.Condition {
	return metav1.Condition{
		Type:    Error,
		Status:  status,
		Reason:  "UnexpectedError",
		Message: message,
	}
}
