package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"runtime"
	"strings"
)

var mod string

type SpanOptions struct {
	operationName  string
	customResource *types.NamespacedName
}

func (o *SpanOptions) OperationName() string {
	if o.operationName == "" {
		pc, _, _, _ := runtime.Caller(2)
		details := runtime.FuncForPC(pc)
		name := details.Name()
		return strings.Replace(name, mod, "", 1)
	}
	return o.operationName
}

type SpanOptionFunc func(*SpanOptions)

func WithOperationName(operation string) SpanOptionFunc {
	return func(o *SpanOptions) {
		o.operationName = operation
	}
}

func WithCustomResource(cr types.NamespacedName) SpanOptionFunc {
	return func(o *SpanOptions) {
		o.customResource = &cr
	}
}

func StartSpanFromContext(ctx context.Context, options ...SpanOptionFunc) (*Span, context.Context) {
	opt := new(SpanOptions)

	for _, fn := range options {
		if fn == nil {
			continue
		}
		fn(opt)
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, opt.OperationName())
	if opt.customResource != nil {
		span.SetTag("kubernetes.resource", opt.customResource.String())
	}
	return &Span{Span: span}, ctx
}

func StartSpanFromRequest(r *http.Request) (*Span, context.Context) {
	ctx := ExtractSpanContextFromRequest(r)
	span, ctxWithSpan := opentracing.StartSpanFromContext(r.Context(), r.Method+" "+r.URL.Path, ext.RPCServerOption(ctx))

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	ext.HTTPUrl.Set(span, fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))
	ext.HTTPMethod.Set(span, r.Method)
	span.SetTag("http.protocol", r.Proto)

	return &Span{Span: span}, ctxWithSpan
}

func SpanFromContext(ctx context.Context) *Span {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}

	return &Span{Span: span}
}

func ExtractSpanContextFromRequest(r *http.Request) opentracing.SpanContext {
	tracer := opentracing.GlobalTracer()
	ctx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	return ctx
}
