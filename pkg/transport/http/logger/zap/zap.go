package zap

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"github.com/alexfalkowski/go-service/pkg/transport/http/encoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	client              = "client"
)

var (
	sensitiveURLs = regexp.MustCompile(`oauth|token|jwks|health|liveness|readiness`)
)

// NewRoundTripper for zap.
func NewRoundTripper(logger *zap.Logger, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{logger: logger, RoundTripper: hrt}
}

// RoundTripper for zap.
type RoundTripper struct {
	logger *zap.Logger

	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now().UTC()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)
	isSensitive := sensitiveURLs.Match([]byte(req.URL.String()))
	fields := []zapcore.Field{
		zap.Int64(httpDuration, time.ToMilliseconds(time.Since(start))),
		zap.String(httpStartTime, start.Format(time.RFC3339)),
		zap.String(httpURL, req.URL.String()),
		zap.String(httpMethod, req.Method),
		zap.String("span.kind", client),
		zap.String(component, httpComponent),
	}

	if !isSensitive {
		fields = append(fields, zap.String(httpRequest, encoder.Request(req)))
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(httpRequestDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		r.logger.Error("finished call with error", fields...)

		return nil, err
	}

	fields = append(fields, zap.Int(httpStatusCode, resp.StatusCode))

	if !isSensitive {
		fields = append(fields, zap.String(httpResponse, encoder.Response(resp)))
	}

	r.logger.Info(fmt.Sprintf("finished call with code %s", resp.Status), fields...)

	return resp, nil
}
