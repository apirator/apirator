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

package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

// Initialize create an instance of Jaeger Tracer and sets it as GlobalTracer.
func Initialize(service string, module string) (io.Closer, error) {
	mod = module
	cfg, err := (&config.Configuration{ServiceName: service}).FromEnv()
	if err != nil {
		return nil, fmt.Errorf("cannot init Jaeger: %w", err)
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, fmt.Errorf("cannot init Jaeger: %w", err)
	}

	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}
