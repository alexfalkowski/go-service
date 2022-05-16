package opentracing

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	sstrings "github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/fx"
)

const (
	httpURL        = "http.url"
	httpMethod     = "http.method"
	httpDeadline   = "http.deadline"
	httpStatusCode = "http.status_code"
	component      = "component"
	httpComponent  = "http"
)

// TracerParams for otr.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *opentracing.Config
	Version   version.Version
}

// NewTracer for otr.
func NewTracer(params TracerParams) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: params.Lifecycle, Name: "http", Config: params.Config, Version: params.Version})
}

// Tracer for otr.
type Tracer otr.Tracer

// Handler for otr.
type Handler struct {
	tracer Tracer
	http.Handler
}

// NewHandler for otr.
func NewHandler(tracer Tracer, handler http.Handler) *Handler {
	return &Handler{tracer: tracer, Handler: handler}
}

// ServeHTTP for otr.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	service, method := req.URL.Path, strings.ToLower(req.Method)
	if sstrings.IsHealth(service) {
		h.Handler.ServeHTTP(resp, req)

		return
	}

	ctx := req.Context()
	operationName := fmt.Sprintf("%s %s", method, service)
	carrier := otr.HTTPHeadersCarrier(req.Header)
	traceCtx, _ := h.tracer.Extract(otr.HTTPHeaders, carrier)
	opts := []otr.StartSpanOption{
		ext.RPCServerOption(traceCtx),
		otr.Tag{Key: httpURL, Value: service},
		otr.Tag{Key: httpMethod, Value: method},
		otr.Tag{Key: component, Value: httpComponent},
		ext.SpanKindRPCServer,
	}

	span := h.tracer.StartSpan(operationName, opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(httpDeadline, d.UTC().Format(time.RFC3339))
	}

	ctx = otr.ContextWithSpan(ctx, span)

	h.Handler.ServeHTTP(resp, req.WithContext(ctx))

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}
}

// NewRoundTripper for otr.
func NewRoundTripper(tracer Tracer, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{tracer: tracer, RoundTripper: hrt}
}

// RoundTripper for otr.
type RoundTripper struct {
	tracer Tracer
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if sstrings.IsHealth(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := req.URL.Hostname(), strings.ToLower(req.Method)
	ctx := req.Context()
	operationName := fmt.Sprintf("%s %s", method, service)
	opts := []otr.StartSpanOption{
		otr.Tag{Key: httpURL, Value: service},
		otr.Tag{Key: httpMethod, Value: method},
		otr.Tag{Key: component, Value: httpComponent},
		ext.SpanKindRPCClient,
	}

	span, ctx := otr.StartSpanFromContextWithTracer(ctx, r.tracer, operationName, opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(httpDeadline, d.UTC().Format(time.RFC3339))
	}

	carrier := otr.HTTPHeadersCarrier(req.Header)
	if err := r.tracer.Inject(span.Context(), otr.HTTPHeaders, carrier); err != nil {
		return nil, err
	}

	resp, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	if err != nil {
		opentracing.SetError(span, err)

		return nil, err
	}

	span.SetTag(httpStatusCode, resp.StatusCode)

	return resp, nil
}
