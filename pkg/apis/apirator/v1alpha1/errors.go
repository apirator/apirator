package v1alpha1

import "fmt"

// ingress controller not configured
type IngressControllerNotFound struct {
	Name string
}

// host already configured
type HostAlreadyExist struct {
	Host    string
	Ingress string
}

// when ingress-controller is necessary and it was not found
func (e IngressControllerNotFound) Error() string {
	return fmt.Sprintf("Ingress-Controller not found %s", e.Name)
}

// when a host was already configured in the desired ingress
func (e HostAlreadyExist) Error() string {
	return fmt.Sprintf("Host %s already configured in ingress %s", e.Host, e.Ingress)
}
