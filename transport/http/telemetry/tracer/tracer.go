package tracer

import (
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	snoop "github.com/felixge/httpsnoop"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// Handler for tracer.
type Handler struct {
	tracer *tracer.Tracer
}

// NewHandler for tracer.
func NewHandler(tracer *tracer.Tracer) *Handler {
	return &Handler{tracer: tracer}
}

// ServeHTTP for tracer.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	path, method := req.URL.Path, strings.ToLower(req.Method)
	if ts.IsObservable(path) {
		next(res, req)

		return
	}

	ctx := extract(req.Context(), req)
	attrs := []attribute.KeyValue{
		semconv.HTTPRoute(path),
		semconv.HTTPRequestMethodKey.String(method),
	}

	ctx, span := h.tracer.StartServer(ctx, operationName(method+" "+path), attrs...)
	defer span.End()

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })

	span.SetAttributes(semconv.HTTPResponseStatusCode(m.Code))
	tracer.Meta(ctx, span)
}

// NewRoundTripper for tracer.
func NewRoundTripper(tracer *tracer.Tracer, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{tracer: tracer, RoundTripper: hrt}
}

// RoundTripper for tracer.
type RoundTripper struct {
	tracer *tracer.Tracer
	http.RoundTripper
}

// RoundTrip for tracer.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if ts.IsObservable(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	method := strings.ToLower(req.Method)
	ctx := req.Context()
	attrs := []attribute.KeyValue{
		semconv.HTTPRoute(req.URL.Path),
		semconv.HTTPRequestMethodKey.String(method),
	}

	ctx, span := r.tracer.StartClient(ctx, operationName(method+" "+req.URL.Redacted()), attrs...)
	defer span.End()

	inject(ctx, req)

	resp, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))

	tracer.Meta(ctx, span)
	tracer.Error(err, span)

	if resp != nil {
		span.SetAttributes(semconv.HTTPResponseStatusCode(resp.StatusCode))
	}

	return resp, err
}

func operationName(name string) string {
	return tracer.OperationName("http", name)
}
