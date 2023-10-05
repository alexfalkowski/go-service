package tracer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/semconv/v1.19.0/httpconv"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Params for tracer.
type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *tracer.Config
	Version   version.Version
}

// NewTracer for tracer.
func NewTracer(params Params) (Tracer, error) {
	return tracer.NewTracer(params.Lifecycle, "http", params.Version, params.Config)
}

// Tracer for tracer.
type Tracer trace.Tracer

// Handler for tracer.
type Handler struct {
	tracer Tracer
	http.Handler
}

// NewHandler for tracer.
func NewHandler(tracer Tracer, handler http.Handler) *Handler {
	return &Handler{tracer: tracer, Handler: handler}
}

// ServeHTTP for tracer.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	service, method := req.URL.Path, strings.ToLower(req.Method)
	if tstrings.IsHealth(service) {
		h.Handler.ServeHTTP(resp, req)

		return
	}

	ctx := extract(req.Context(), req)
	attrs := []attribute.KeyValue{
		semconv.HTTPRoute(service),
		semconv.HTTPMethod(method),
	}
	attrs = append(attrs, httpconv.ServerRequest("", req)...)
	operationName := fmt.Sprintf("%s %s", method, service)

	ctx, span := h.tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		operationName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	h.Handler.ServeHTTP(resp, req.WithContext(ctx))

	for k, v := range meta.Attributes(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}
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
		semconv.HTTPMethod(method),
	}
	attrs = append(attrs, httpconv.ClientRequest(req)...)

	ctx, span := r.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	inject(ctx, req)

	resp, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))

	for k, v := range meta.Attributes(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return nil, err
	}

	span.SetAttributes(semconv.HTTPStatusCode(resp.StatusCode))
	span.SetAttributes(httpconv.ClientResponse(resp)...)

	return resp, nil
}
