package apimock

import (
	"context"
	apirator "github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileAPIMock) processAnnotations(mock *apirator.APIMock) (updated bool, err error) {
	svcPresent := mock.Spec.ServiceDefinition.Port != 0
	if svcPresent {
		svcK8s := &v1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      mock.GetName(),
			Namespace: mock.Namespace,
		}, svcK8s)
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		updateByIP := mock.AnnotateClusterIP(svcK8s.Spec.ClusterIP)
		updateByPort := mock.AnnotatePorts(svcK8s.Spec.Ports)
		if updateByIP || updateByPort {
			log.Info("APIMock annotations update successfully", "Service.Namespace", svcK8s.Namespace, "Service.Name", svcK8s.Name, "Service.ClusterIP", mock.Annotations["apirator.io/cluster-ip"], "Service.Ports", mock.Annotations["apirator.io/ports"])
			return true, nil
		}
	}
	return false, nil
}
