package http

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type traceRoundTripper struct {
	http.RoundTripper
}

func (r *traceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now().UTC()
	ctx := req.Context()
	tracer := opentracing.GlobalTracer()
	operationName := fmt.Sprintf("%s %s", req.Method, req.URL.Hostname())
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: httpStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: httpURL, Value: req.URL.String()},
		opentracing.Tag{Key: httpMethod, Value: req.Method},
		opentracing.Tag{Key: httpRequest, Value: encodeRequest(req)},
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
	clientSpan.SetTag(httpResponse, encodeResponse(resp))

	return resp, nil
}
