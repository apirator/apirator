package manager

import (
	"github.com/google/wire"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var Providers = wire.NewSet(
	Client,
	RestConfig,
	Scheme,
)

func Scheme(mgr manager.Manager) *runtime.Scheme {
	return mgr.GetScheme()
}

func Client(mgr manager.Manager) client.Client {
	return mgr.GetClient()
}

func RestConfig(mgr manager.Manager) *rest.Config {
	return mgr.GetConfig()
}
