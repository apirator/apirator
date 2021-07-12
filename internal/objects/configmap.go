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
	"path/filepath"

	"github.com/apirator/apirator/api/v1alpha1"
	yu "github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	yamlConfigPath = "/etc/oas/oas.yaml"
	jsonConfigPath = "/etc/oas/oas.json"
)

func (b *Builder) ConfigMapFor(apimock *v1alpha1.APIMock) (*corev1.ConfigMap, error) {
	bJson, err := yu.YAMLToJSON([]byte(apimock.Spec.Definition))
	if err != nil {
		return nil, fmt.Errorf("failed to convert openapi definition to JSON: %w", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apimock.GetName(),
			Namespace: apimock.GetNamespace(),
			Labels:    apimock.MatchLabels(),
		},
		Data: map[string]string{
			filepath.Base(yamlConfigPath): apimock.Spec.Definition,
			filepath.Base(jsonConfigPath): string(bJson),
		},
	}

	if err = b.SetControllerReference(apimock, cm); err != nil {
		return nil, err
	}

	return cm, nil
}
