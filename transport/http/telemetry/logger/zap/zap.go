package zap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	shttp "github.com/alexfalkowski/go-service/http"
	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	httpURL        = "http.url"
	httpMethod     = "http.method"
	httpDuration   = "http.duration_ms"
	httpStartTime  = "http.start_time"
	httpDeadline   = "http.deadline"
	httpStatusCode = "http.status_code"
	kind           = "kind"
	httpKind       = "http"
	client         = "client"
	server         = "server"
)

// NewHandler for zap.
func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{logger: logger}
}

// Handler for zap.
type Handler struct {
	logger *zap.Logger
}

// ServeHTTP or zap.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	service, method := req.URL.Path, strings.ToLower(req.Method)
	if tstrings.IsHealth(service) {
		next(resp, req)

		return
	}

	start := time.Now().UTC()
	ctx := req.Context()

	res := &shttp.ResponseWriter{ResponseWriter: resp, StatusCode: http.StatusOK}
	next(res, req)

	fields := []zapcore.Field{
		zap.Int64(httpDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(httpStartTime, start.Format(time.RFC3339)),
		zap.String(httpURL, service),
		zap.String(httpMethod, method),
		zap.String("http.kind", server),
		zap.String(kind, httpKind),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(httpDeadline, d.UTC().Format(time.RFC3339)))
	}

	fields = append(fields, zap.Int(httpStatusCode, res.StatusCode))

	loggerLevel := codeToLevel(res.StatusCode, h.logger)
	loggerLevel(fmt.Sprintf("finished call with code %s", res.Status()), fields...)
}

// NewRoundTripper for zap.
func NewRoundTripper(logger *zap.Logger, r http.RoundTripper) *RoundTripper {
	return &RoundTripper{logger: logger, RoundTripper: r}
}

// RoundTripper for zap.
type RoundTripper struct {
	logger *zap.Logger

	http.RoundTripper
}

// RoundTrip for zap.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if tstrings.IsHealth(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := req.URL.Hostname(), strings.ToLower(req.Method)
	start := time.Now().UTC()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)
	fields := []zapcore.Field{
		zap.Int64(httpDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(httpStartTime, start.Format(time.RFC3339)),
		zap.String(httpURL, service),
		zap.String(httpMethod, method),
		zap.String("http.kind", client),
		zap.String(kind, httpKind),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(httpDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		r.logger.Error("finished call with error", fields...)

		return nil, err
	}

	fields = append(fields, zap.Int(httpStatusCode, resp.StatusCode))

	loggerLevel := codeToLevel(resp.StatusCode, r.logger)
	loggerLevel(fmt.Sprintf("finished call with code %s", resp.Status), fields...)

	return resp, nil
}

func codeToLevel(code int, logger *zap.Logger) func(msg string, fields ...zapcore.Field) {
	if code >= 400 && code <= 499 {
		return logger.Warn
	}

	if code >= 500 && code <= 599 {
		return logger.Error
	}

	return logger.Info
}
