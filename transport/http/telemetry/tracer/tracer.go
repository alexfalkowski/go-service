package tracer

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	snoop "github.com/felixge/httpsnoop"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

// Tracer is an alias for tracer.Tracer.
type Tracer = tracer.Tracer

// NewHandler for tracer.
func NewHandler(tracer *Tracer) *Handler {
	return &Handler{tracer: tracer}
}

// Handler for tracer.
type Handler struct {
	tracer *Tracer
}

// ServeHTTP for tracer.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	p, method := http.Path(req), strings.ToLower(req.Method)
	if strings.IsObservable(p) {
		next(res, req)

		return
	}

	ctx := extract(req.Context(), req)
	attrs := []attribute.KeyValue{
		semconv.HTTPRoute(p),
		semconv.HTTPRequestMethodKey.String(method),
	}

	op := operationName(strings.Join(" ", method, p))

	ctx, span := h.tracer.StartServer(ctx, op, attrs...)
	defer span.End()

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })

	span.SetAttributes(semconv.HTTPResponseStatusCode(m.Code))
	tracer.Meta(ctx, span)
}

// NewRoundTripper for tracer.
func NewRoundTripper(tracer *Tracer, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{tracer: tracer, RoundTripper: hrt}
}

// RoundTripper for tracer.
type RoundTripper struct {
	tracer *Tracer
	http.RoundTripper
}

// RoundTrip for tracer.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	p, method := http.Path(req), strings.ToLower(req.Method)
	if strings.IsObservable(p) {
		return r.RoundTripper.RoundTrip(req)
	}

	ctx := req.Context()
	attrs := []attribute.KeyValue{
		semconv.HTTPRoute(p),
		semconv.HTTPRequestMethodKey.String(method),
	}

	op := operationName(strings.Join(" ", method, req.URL.Redacted()))

	ctx, span := r.tracer.StartClient(ctx, op, attrs...)
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
