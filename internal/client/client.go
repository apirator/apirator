package client

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Kubernetes struct {
	client.Client
}
