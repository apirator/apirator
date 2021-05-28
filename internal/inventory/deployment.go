package inventory

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	Create []appsv1.Deployment
	Update []appsv1.Deployment
	Delete []appsv1.Deployment
}

func ForDeployments(existing, desired []appsv1.Deployment) Object {
	var update []client.Object
	mcreate := deploymentMap(desired)
	mdelete := deploymentMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			tp := t.DeepCopy()

			if tp.Spec.Replicas != nil && v.Spec.Replicas == nil {
				v.Spec.Replicas = tp.Spec.Replicas
			}

			tp.Spec = v.Spec
			tp.ObjectMeta.OwnerReferences = v.ObjectMeta.OwnerReferences

			for k, v := range v.ObjectMeta.Annotations {
				tp.ObjectMeta.Annotations[k] = v
			}

			for k, v := range v.ObjectMeta.Labels {
				tp.ObjectMeta.Labels[k] = v
			}

			update = append(update, tp)
			delete(mcreate, k)
			delete(mdelete, k)
		}
	}

	return Object{
		Create: deploymentList(mcreate),
		Update: update,
		Delete: deploymentList(mdelete),
	}
}

func deploymentMap(deps []appsv1.Deployment) map[string]appsv1.Deployment {
	m := map[string]appsv1.Deployment{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

func deploymentList(m map[string]appsv1.Deployment) []client.Object {
	var l []client.Object
	for _, v := range m {
		l = append(l, &v)
	}
	return l
}
