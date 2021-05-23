package k8s

import (
	"fmt"

	"github.com/google/wire"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

var Providers = wire.NewSet(NewCoreClient)

func NewCoreClient(cfg *rest.Config) (*core.CoreV1Client, error) {
	client, err := core.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get core/v1 client: %w", err)
	}
	return client, nil
}
