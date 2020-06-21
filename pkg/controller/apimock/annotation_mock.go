package apimock

import (
	"context"
	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileAPIMock) processAnnotations(mock *apirator.APIMock) (updated bool, err error) {
	var updateByIP = false
	var updateByPort = false
	var updateByIng = false
	svcPresent := mock.Spec.ServiceDefinition.Port != 0
	if svcPresent {
		svcK8s := &v1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      mock.GetName(),
			Namespace: mock.Namespace,
		}, svcK8s)
		if err != nil && errors.IsNotFound(err) {
			log.Error(err, "Service not found.", "Service.Name", mock.GetName())
			return false, err
		}
		updateByIP = mock.AnnotateClusterIP(svcK8s.Spec.ClusterIP)
		updateByPort = mock.AnnotatePorts(svcK8s.Spec.Ports)
		if updateByIP || updateByPort {
			log.Info("APIMock annotations update successfully", "Service.Namespace", svcK8s.Namespace, "Service.Name", svcK8s.Name, "Service.ClusterIP", mock.Annotations["apirator.io/cluster-ip"], "Service.Ports", mock.Annotations["apirator.io/ports"])
		}
	}

	if mock.ExposeInIngress() {
		ingK8s := &v1beta1.Ingress{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      mock.Spec.Selector[apirator.IngressTag],
			Namespace: mock.Spec.Selector[apirator.NamespaceTag],
		}, ingK8s)
		if err != nil && errors.IsNotFound(err) {
			log.Error(err, "Ingress not found.", "Ingress.Name", mock.Spec.Selector[apirator.IngressTag])
			return false, err
		}
		log.Info("Configuring mock address...", "Mock.Name", mock.GetName())
		if len(ingK8s.Status.LoadBalancer.Ingress) > 0 {
			updateByIng = mock.AnnotateAddress(ingK8s.Status.LoadBalancer.Ingress[0].IP)
		}
	}
	return updateByIng || updateByPort || updateByIP, nil
}
