package http

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerRoundTripper struct {
	logger *zap.Logger

	http.RoundTripper
}

func (r *loggerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now().UTC()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)
	fields := []zapcore.Field{
		zap.Int64(httpDuration, time.ToMilliseconds(time.Since(start))),
		zap.String(httpStartTime, start.Format(time.RFC3339)),
		zap.String(httpURL, req.URL.String()),
		zap.String(httpMethod, req.Method),
		zap.String(httpRequest, encodeRequest(req)),
		zap.String("span.kind", client),
		zap.String(component, httpComponent),
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

	fields = append(fields, zap.Int(httpStatusCode, resp.StatusCode), zap.String(httpResponse, encodeResponse(resp)))

	r.logger.Info(fmt.Sprintf("finished call with code %s", resp.Status), fields...)

	return resp, nil
}
