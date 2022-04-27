package opentracing

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	sstrings "github.com/alexfalkowski/go-service/strings"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	httpURL             = "http.url"
	httpMethod          = "http.method"
	httpDuration        = "http.duration_ms"
	httpStartTime       = "http.start_time"
	httpRequestDeadline = "http.request.deadline"
	httpStatusCode      = "http.status_code"
	component           = "component"
	httpComponent       = "http"
)

// Tracer for opentracing.
type Tracer opentracing.Tracer

// Handler for opentracing.
type Handler struct {
	tracer Tracer
	http.Handler
}

// NewHandler for opentracing.
func NewHandler(tracer Tracer, handler http.Handler) *Handler {
	return &Handler{tracer: tracer, Handler: handler}
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	if sstrings.IsHealth(url) {
		h.Handler.ServeHTTP(resp, req)

		return
	}

	start := time.Now().UTC()
	ctx := req.Context()
	method := strings.ToLower(req.Method)
	operationName := fmt.Sprintf("%s %s", method, url)
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	traceCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, carrier)
	opts := []opentracing.StartSpanOption{
		ext.RPCServerOption(traceCtx),
		opentracing.Tag{Key: httpStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: httpURL, Value: url},
		opentracing.Tag{Key: httpMethod, Value: method},
		opentracing.Tag{Key: component, Value: httpComponent},
		ext.SpanKindRPCServer,
	}

	span := h.tracer.StartSpan(operationName, opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(httpRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	ctx = opentracing.ContextWithSpan(ctx, span)

	h.Handler.ServeHTTP(resp, req.WithContext(ctx))

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	span.SetTag(httpDuration, stime.ToMilliseconds(time.Since(start)))
}

// NewRoundTripper for opentracing.
func NewRoundTripper(tracer Tracer, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{tracer: tracer, RoundTripper: hrt}
}

// RoundTripper for opentracing.
type RoundTripper struct {
	tracer Tracer
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	if sstrings.IsHealth(url) {
		return r.RoundTripper.RoundTrip(req)
	}

	start := time.Now().UTC()
	ctx := req.Context()
	method := strings.ToLower(req.Method)
	operationName := fmt.Sprintf("%s %s", method, url)
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: httpStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: httpURL, Value: url},
		opentracing.Tag{Key: httpMethod, Value: method},
		opentracing.Tag{Key: component, Value: httpComponent},
		ext.SpanKindRPCClient,
	}

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, operationName, opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(httpRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	if err := r.tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier); err != nil {
		return nil, err
	}

	resp, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	span.SetTag(httpDuration, stime.ToMilliseconds(time.Since(start)))

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))

		return nil, err
	}

	span.SetTag(httpStatusCode, resp.StatusCode)

	return resp, nil
}
