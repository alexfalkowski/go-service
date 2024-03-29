package tracer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/env"
	shttp "github.com/alexfalkowski/go-service/http"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Params for tracer.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *tracer.Config
	Environment env.Environment
	Version     version.Version
}

// NewTracer for tracer.
func NewTracer(params Params) (Tracer, error) {
	return tracer.NewTracer(context.Background(), params.Lifecycle, "http", params.Environment, params.Version, params.Config)
}

// Tracer for tracer.
type Tracer trace.Tracer

// Handler for tracer.
type Handler struct {
	tracer Tracer
}

// NewHandler for tracer.
func NewHandler(tracer Tracer) *Handler {
	return &Handler{tracer: tracer}
}

// ServeHTTP for tracer.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	path, method := req.URL.Path, strings.ToLower(req.Method)
	if tstrings.IsHealth(path) {
		next(resp, req)

		return
	}

	ctx := extract(req.Context(), req)
	attrs := []attribute.KeyValue{
		semconv.HTTPRoute(path),
		semconv.HTTPRequestMethodKey.String(method),
	}

	operationName := fmt.Sprintf("%s %s", method, path)

	ctx, span := h.tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		operationName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	res := &shttp.ResponseWriter{ResponseWriter: resp, StatusCode: http.StatusOK}
	next(res, req.WithContext(ctx))

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	span.SetAttributes(semconv.HTTPResponseStatusCode(res.StatusCode))
}

// NewRoundTripper for tracer.
func NewRoundTripper(tracer Tracer, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{tracer: tracer, RoundTripper: hrt}
}

// RoundTripper for tracer.
type RoundTripper struct {
	tracer Tracer
	http.RoundTripper
}

// RoundTrip for tracer.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if tstrings.IsHealth(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := req.URL.Hostname(), strings.ToLower(req.Method)
	ctx := req.Context()
	operationName := fmt.Sprintf("%s %s", method, service)
	attrs := []attribute.KeyValue{
		semconv.HTTPRoute(service),
		semconv.HTTPRequestMethodKey.String(method),
	}

	ctx, span := r.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	inject(ctx, req)

	resp, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return nil, err
	}

	span.SetAttributes(semconv.HTTPResponseStatusCode(resp.StatusCode))

	return resp, nil
}
