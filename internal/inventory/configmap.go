package inventory

import (
	"fmt"

	core "k8s.io/api/core/v1"
)

type ConfigMap struct {
	Create []core.ConfigMap
	Update []core.ConfigMap
	Delete []core.ConfigMap
}

func ForConfigMaps(existing, desired []core.ConfigMap) ConfigMap {
	var update []core.ConfigMap
	mcreate := configsMap(desired)
	mdelete := configsMap(existing)

	for k, v := range mcreate {
		if t, ok := mdelete[k]; ok {
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

			update = append(update, *tp)
			delete(mcreate, k)
			delete(mdelete, k)
		}
	}

	return ConfigMap{
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

func configsList(m map[string]core.ConfigMap) []core.ConfigMap {
	var l []core.ConfigMap
	for _, v := range m {
		l = append(l, v)
	}
	return l
}
