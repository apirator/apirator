package inventory

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Object struct {
	Create []client.Object
	Update []client.Object
	Delete []client.Object
}
