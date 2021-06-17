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

package adapter

import (
	"context"

	"github.com/apirator/apirator/controllers"
	"github.com/apirator/apirator/internal/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Factory struct {
	userCases *UserCases
	svc       *k8s.Service
}

func NewFactory(userCases *UserCases, svc *k8s.Service) *Factory {
	return &Factory{userCases: userCases, svc: svc}
}

func (a *Factory) CreateAPIMockAdapter(ctx context.Context, key client.ObjectKey) (controllers.APIMockAdapter, error) {
	resource, err := a.svc.GetAPIMock(ctx, key)
	if err != nil {
		return nil, err
	}
	return newAdapter(a.userCases, resource), err
}
