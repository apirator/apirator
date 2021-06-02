package tracing

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-logr/logr"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	ctrl "sigs.k8s.io/controller-runtime"
)

type (
	Span struct {
		opentracing.Span
	}
	Info struct {
		TraceID   string
		SpanID    string
		ParentID  string
		IsSampled bool
	}
)

func (s *Span) String() string {
	return fmt.Sprintf("%+v", s.Span)
}

func (s *Span) Info() *Info {
	if sc, ok := s.Context().(jaeger.SpanContext); ok {
		return &Info{
			TraceID:   sc.TraceID().String(),
			SpanID:    sc.SpanID().String(),
			ParentID:  sc.ParentID().String(),
			IsSampled: sc.IsSampled(),
		}
	}
	return nil
}

func (s *Span) SetError(err error) {
	ext.Error.Set(s, true)
	fields := make([]log.Field, 0)
	fields = append(fields, log.String("event", "error"), log.String("message", err.Error()))

	var serr *kerrors.StatusError
	if errors.As(err, &serr) {
		status := serr.Status()
		fields = append(fields, log.Int32("code", status.Code), log.String("reason", string(status.Reason)))
	}
	s.LogFields(fields...)
}

func (s *Span) HandleError(err error) error {
	s.SetError(err)
	return err
}

func (s *Span) SetHTTPResponseStatus(status int) {
	ext.HTTPStatusCode.Set(s, uint16(status))

	if status >= 500 && status < 600 {
		ext.Error.Set(s, true)
		s.SetTag("error.type", fmt.Sprintf("%d: %s", status, http.StatusText(status)))
		s.LogKV(
			"event", "error",
			"message", fmt.Sprintf("%d: %s", status, http.StatusText(status)),
		)
	}
}

func (s *Span) Panic(err interface{}) {
	ext.HTTPStatusCode.Set(s, uint16(500))
	ext.Error.Set(s, true)
	s.SetTag("error.type", "panic")
	s.LogKV("event", "error",
		"error.kind", "panic",
		"message", err,
		"stack", string(debug.Stack()))
	panic(err)
}

func (s *Span) LoggerWithName(name string) logr.Logger {
	return ctrl.Log.WithName(name).WithValues("trace", s.String())
}

func (s *Span) Logger() logr.Logger {
	return ctrl.Log.WithValues("trace", s.String())
}
