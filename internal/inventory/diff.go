package inventory

import "github.com/google/go-cmp/cmp"

func ignore(fields ...string) cmp.Option {
	return cmp.FilterPath(func(path cmp.Path) bool {
		for _, p := range fields {
			if p == path.String() {
				return true
			}
		}
		return false
	}, cmp.Ignore())
}
