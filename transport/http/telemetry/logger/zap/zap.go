package zap

import (
	"net/http"
	"strings"
	"time"

	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	snoop "github.com/felixge/httpsnoop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	service = "http"
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
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	path, method := req.URL.Path, strings.ToLower(req.Method)
	if ts.IsObservable(path) {
		next(res, req)

		return
	}

	ctx := req.Context()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, service),
		zap.String(meta.PathKey, path),
		zap.String(meta.MethodKey, method),
	}

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })

	fields = append(fields, zap.Stringer(meta.DurationKey, m.Duration), zap.Int(meta.CodeKey, m.Code))
	fields = append(fields, tz.Meta(ctx)...)

	tz.LogWithFunc(message(method+" "+path), nil, codeToLevel(m.Code, h.logger), fields...)
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
	if ts.IsObservable(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	path, method := req.URL.Path, strings.ToLower(req.Method)
	start := time.Now()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)
	fields := []zapcore.Field{
		zap.Stringer(meta.DurationKey, time.Since(start)),
		zap.String(meta.ServiceKey, service),
		zap.String(meta.PathKey, path),
		zap.String(meta.MethodKey, method),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if resp != nil {
		fields = append(fields, zap.Int(meta.CodeKey, resp.StatusCode))
	}

	tz.LogWithFunc(message(method+" "+req.URL.Redacted()), err, respToLevel(resp, r.logger), fields...)

	return resp, err
}

func respToLevel(resp *http.Response, logger *zap.Logger) func(msg string, fields ...zapcore.Field) {
	var code int

	if resp != nil {
		code = resp.StatusCode
	} else {
		code = 500
	}

	return codeToLevel(code, logger)
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

func message(msg string) string {
	return "http: " + msg
}
