package opentracing

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/http/encoder"
	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	httpRequest         = "http.request"
	httpResponse        = "http.response"
	httpURL             = "http.url"
	httpMethod          = "http.method"
	httpDuration        = "http.duration_ms"
	httpStartTime       = "http.start_time"
	httpRequestDeadline = "http.request.deadline"
	httpStatusCode      = "http.status_code"
	component           = "component"
	httpComponent       = "http"
)

// NewRoundTripper for opentracing.
func NewRoundTripper(hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt}
}

// RoundTripper for opentracing.
type RoundTripper struct {
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now().UTC()
	ctx := req.Context()
	tracer := opentracing.GlobalTracer()
	operationName := fmt.Sprintf("%s %s", req.Method, req.URL.Hostname())
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: httpStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: httpURL, Value: req.URL.String()},
		opentracing.Tag{Key: httpMethod, Value: req.Method},
		opentracing.Tag{Key: httpRequest, Value: encoder.Request(req)},
		opentracing.Tag{Key: component, Value: httpComponent},
		ext.SpanKindRPCClient,
	}

	for k, v := range meta.Attributes(ctx) {
		opts = append(opts, opentracing.Tag{Key: k, Value: v})
	}

	clientSpan, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName, opts...)

	defer clientSpan.Finish()

	if d, ok := ctx.Deadline(); ok {
		clientSpan.SetTag(httpRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, carrier); err != nil {
		return nil, err
	}

	resp, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))

	clientSpan.SetTag(httpDuration, time.ToMilliseconds(time.Since(start)))

	if err != nil {
		ext.Error.Set(clientSpan, true)
		clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))

		return nil, err
	}

	clientSpan.SetTag(httpStatusCode, resp.StatusCode)
	clientSpan.SetTag(httpResponse, encoder.Response(resp))

	return resp, nil
}
