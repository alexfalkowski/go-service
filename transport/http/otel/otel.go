package otel

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/otel"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/semconv/v1.18.0/httpconv"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// TracerParams for otel.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *otel.Config
	Version   version.Version
}

// NewTracer for otel.
func NewTracer(params TracerParams) (Tracer, error) {
	return otel.NewTracer(otel.TracerParams{Lifecycle: params.Lifecycle, Name: "http", Config: params.Config, Version: params.Version})
}

// Tracer for otel.
type Tracer trace.Tracer

// Handler for otel.
type Handler struct {
	tracer Tracer
	http.Handler
}

// NewHandler for otel.
func NewHandler(tracer Tracer, handler http.Handler) *Handler {
	return &Handler{tracer: tracer, Handler: handler}
}

// ServeHTTP for otel.
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

// NewRoundTripper for otr.
func NewRoundTripper(tracer Tracer, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{tracer: tracer, RoundTripper: hrt}
}

// RoundTripper for otel.
type RoundTripper struct {
	tracer Tracer
	http.RoundTripper
}

// RoundTrip for otel.
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
