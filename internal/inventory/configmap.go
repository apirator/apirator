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
	core "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ForConfigMaps(existing, desired []core.ConfigMap) Object {
	var update []client.Object
	mcreate := configsMap(desired)
	mdelete := configsMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
			diff := cmp.Diff(v, t, ignore(configMapFields...))
			if diff != "" {
				tp := t.DeepCopy()

				tp.Data = v.Data
				tp.BinaryData = v.BinaryData
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
		Create: configsList(mcreate),
		Update: update,
		Delete: configsList(mdelete),
	}
}

func configsMap(deps []core.ConfigMap) map[string]core.ConfigMap {
	m := map[string]core.ConfigMap{}
	for _, d := range deps {
		m[fmt.Sprintf("%s.%s", d.Namespace, d.Name)] = d
	}
	return m
}

func configsList(m map[string]core.ConfigMap) []client.Object {
	var l []client.Object
	for _, v := range m {
		l = append(l, &v)
	}
	return l
}

var configMapFields = []string{
	"TypeMeta.Kind",
	"TypeMeta.APIVersion",
	"ObjectMeta.SelfLink",
	"ObjectMeta.UID",
	"ObjectMeta.ResourceVersion",
	"ObjectMeta.CreationTimestamp",
	"ObjectMeta.ManagedFields",
}
