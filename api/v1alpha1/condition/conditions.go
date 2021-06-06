package condition

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// These are valid conditions of a deployment APIMock.
const (
	Provisioned string = "Provisioned"
	Error              = "Error"
	Validated          = "Validated"
	Available          = "Available"
	Waiting            = "Waiting"
)

func NewOpenAPIValidation(status bool, message string) metav1.Condition {
	return metav1.Condition{
		Type:    Validated,
		Status:  statusFor(&status),
		Reason:  "ValidOpenAPIDefinition",
		Message: message,
	}
}

func NewError(status bool, message string) metav1.Condition {
	return metav1.Condition{
		Type:    Error,
		Status:  statusFor(&status),
		Reason:  "UnexpectedError",
		Message: message,
	}
}

func NewAvailable(status bool, message string) metav1.Condition {
	return metav1.Condition{
		Type:    Available,
		Status:  statusFor(&status),
		Reason:  "MinimumReplicasAvailable",
		Message: message,
	}
}

func statusFor(b *bool) metav1.ConditionStatus {
	if b == nil {
		return metav1.ConditionUnknown
	}
	if *b {
		return metav1.ConditionTrue
	}
	return metav1.ConditionFalse
}
