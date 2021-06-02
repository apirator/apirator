package resources

import "k8s.io/apimachinery/pkg/runtime"

type Builder struct {
	scheme *runtime.Scheme
}

func NewBuilder(scheme *runtime.Scheme) *Builder {
	return &Builder{scheme: scheme}
}
