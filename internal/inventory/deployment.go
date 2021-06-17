// Copyright 2020 apirator.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inventory

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForDeployments(existing, desired []appsv1.Deployment) Object {
	var update []client.Object
	mcreate := deploymentMap(desired)
	mdelete := deploymentMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(v, t, ignore(deploymentFields...))
			if diff != "" {
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
			}
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

var deploymentFields = []string{
	"ObjectMeta.Annotations",
	"ObjectMeta.CreationTimestamp",
	"ObjectMeta.Generation",
	"ObjectMeta.ManagedFields",
	"ObjectMeta.ResourceVersion",
	"ObjectMeta.SelfLink",
	"ObjectMeta.UID",
	"Spec.ProgressDeadlineSeconds",
	"Spec.RevisionHistoryLimit",
	"Spec.Template.Spec.Containers.ImagePullPolicy",
	"Spec.Template.Spec.Containers.LivenessProbe.FailureThreshold",
	"Spec.Template.Spec.Containers.LivenessProbe.SuccessThreshold",
	"Spec.Template.Spec.Containers.ReadinessProbe.FailureThreshold",
	"Spec.Template.Spec.Containers.ReadinessProbe.SuccessThreshold",
	"Spec.Template.Spec.Containers.TerminationMessagePath",
	"Spec.Template.Spec.Containers.TerminationMessagePolicy",
	"Spec.Template.Spec.DNSPolicy",
	"Spec.Template.Spec.RestartPolicy",
	"Spec.Template.Spec.SchedulerName",
	"Spec.Template.Spec.SecurityContext",
	"Spec.Template.Spec.TerminationGracePeriodSeconds",
	"Spec.Template.Spec.Volumes.VolumeSource.ConfigMap.DefaultMode",
	"Status",
	"TypeMeta.APIVersion",
	"TypeMeta.Kind",
}
