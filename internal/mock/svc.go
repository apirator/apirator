package mock

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct {
	client.Client
}

func NewService(client client.Client) *Service {
	return &Service{Client: client}
}
