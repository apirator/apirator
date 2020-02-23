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

package apimock

import (
	"context"
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
	"github.com/apirator/apirator/pkg/controller/k8s/util/labels"
	"github.com/apirator/apirator/pkg/controller/k8s/util/owner"
	yu "github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"path/filepath"
)

const (
	yamlConfigPath = "/etc/oas/oas.yaml"
	jsonConfigPath = "/etc/oas/oas.json"
)

func (r *ReconcileAPIMock) EnsureConfigMap(mock *v1alpha1.APIMock) error {
	// find the configmap element
	cMap := &v1.ConfigMap{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      mock.GetName(),
		Namespace: mock.Namespace,
	}, cMap)

	// the OAS Definition was re-defined. We should update the item definition
	// https://itnext.io/how-to-automatically-update-your-kubernetes-app-configuration-d750e0ca79ab
	if err != nil && cMap != nil && v1alpha1.PROVISIONED == mock.Status.Phase && mock.Spec.Watch {
		log.Info("Starting updating configmap", "ConfigMap.Namespace", mock.Namespace, "ConfigMap.Name", mock.Name)
		bJson, jsonErr := yu.YAMLToJSON([]byte(mock.Spec.Definition))
		if jsonErr != nil {
			log.Error(err, "[UPDATE] - Failed to read oas spec in json model", "APIMock.Namespace", mock.Namespace, "APIMock.Name", mock.Name)
			return err
		}
		cMap.Data = map[string]string{filepath.Base(yamlConfigPath): mock.Spec.Definition,
			filepath.Base(jsonConfigPath): string(bJson)}
		updErr := r.client.Update(context.TODO(), cMap)
		if updErr != nil {
			log.Error(err, "Failed to update ConfigMap", "ConfigMap.Namespace", mock.Namespace, "ConfigMap.Name", mock.Name)
			return err
		}
		log.Info("ConfigMap update successfully", "ConfigMap.Namespace", mock.Namespace, "ConfigMap.Name", mock.Name)
		return nil
	}

	// The first interaction. Things should be created.
	if err != nil && errors.IsNotFound(err) {
		log.Info("ConfigMap not found. Starting creation...", "ConfigMap.Namespace", mock.Namespace, "ConfigMap.Name", mock.Name)

		bJson, jsonErr := yu.YAMLToJSON([]byte(mock.Spec.Definition))
		if jsonErr != nil {
			log.Error(err, "[CREATE] - Failed to read oas spec in json model", "APIMock.Namespace", mock.Namespace, "APIMock.Name", mock.Name)
			return err
		}
		cm := &v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: mock.Namespace,
				Name:      mock.GetName(),
				Labels:    labels.LabelForAPIMock(mock),
			},
			Data: map[string]string{filepath.Base(yamlConfigPath): mock.Spec.Definition,
				filepath.Base(jsonConfigPath): string(bJson)},
		}
		owner.AddOwnerRefToObject(cm, owner.AsOwner(&mock.ObjectMeta))
		err = r.client.Create(context.TODO(), cm)
		if err != nil {
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return err
		}
		log.Info("ConfigMap created successfully", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		return nil
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return err
	}
	return nil
}
