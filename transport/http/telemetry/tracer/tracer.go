package tracer

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	snoop "github.com/felixge/httpsnoop"
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
	if strings.IsIgnorable(req.URL.Path) {
		next(res, req)
		return
	}

	service, method := http.ParseServiceMethod(req)
	ctx := extract(req.Context(), req)
	op := operationName(strings.Join(strings.Space, method, service))

	ctx, span := h.tracer.StartServer(ctx, op,
		attributes.HTTPRoute(service),
		attributes.HTTPRequestMethod(method))
	defer span.End()

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })
	span.SetAttributes(attributes.HTTPResponseStatusCode(m.Code))
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
	if strings.IsIgnorable(req.URL.Path) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := http.ParseServiceMethod(req)
	ctx := req.Context()
	op := operationName(strings.Join(strings.Space, method, req.URL.Redacted()))

	ctx, span := r.tracer.StartClient(ctx, op,
		attributes.HTTPRoute(service),
		attributes.HTTPRequestMethod(method))
	defer span.End()

	inject(ctx, req)

	resp, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))

	tracer.Meta(ctx, span)
	tracer.Error(err, span)

	if resp != nil {
		span.SetAttributes(attributes.HTTPResponseStatusCode(resp.StatusCode))
	}

	return resp, err
}

func operationName(name string) string {
	return tracer.OperationName("http", name)
}
