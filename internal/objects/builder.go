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

package objects

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Builder struct {
	scheme *runtime.Scheme
}

func NewBuilder(scheme *runtime.Scheme) *Builder {
	return &Builder{scheme: scheme}
}

func (b *Builder) SetOwnerReference(owner, object metav1.Object) error {
	if err := controllerutil.SetOwnerReference(owner, object, b.scheme); err != nil {
		return fmt.Errorf("failed to set %T %q owner reference: %v", object, object.GetName(), err)
	}
	return nil
}

func (b *Builder) SetControllerReference(controller, object metav1.Object) error {
	if err := controllerutil.SetControllerReference(controller, object, b.scheme); err != nil {
		return fmt.Errorf("failed to set %T %q controller reference: %v", object, object.GetName(), err)
	}
	return nil
}
