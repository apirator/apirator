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

	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForIngresses(existing, desired []networkingv1.Ingress) Object {
	var update []client.Object
	mcreate := ingressMap(desired)
	mdelete := ingressMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(v, t, ignore(ingressFields...))
			if diff != "" {
				tp := t.DeepCopy()

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
		Create: ingressList(mcreate),
		Update: update,
		Delete: ingressList(mdelete),
	}
}

func ingressMap(deps []networkingv1.Ingress) map[string]networkingv1.Ingress {
	m := map[string]networkingv1.Ingress{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

func ingressList(m map[string]networkingv1.Ingress) []client.Object {
	var l []client.Object
	for _, v := range m {
		l = append(l, &v)
	}
	return l
}

var ingressFields = []string{
	"ObjectMeta.CreationTimestamp",
	"ObjectMeta.Generation",
	"ObjectMeta.ManagedFields",
	"ObjectMeta.ResourceVersion",
	"ObjectMeta.SelfLink",
	"ObjectMeta.UID",
	"TypeMeta.APIVersion",
	"TypeMeta.Kind",
	"Status",
}
