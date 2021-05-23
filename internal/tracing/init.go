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
