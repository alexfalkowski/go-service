package zap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	ss "github.com/alexfalkowski/go-service/transport/strings"
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
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	path, method := req.URL.Path, strings.ToLower(req.Method)
	if ss.IsHealth(path) {
		next(resp, req)

		return
	}

	ctx := req.Context()
	fields := []zapcore.Field{
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, path),
		zap.String(tm.MethodKey, method),
	}

	m := snoop.CaptureMetricsFn(resp, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })

	fields = append(fields, zap.Stringer(tm.DurationKey, m.Duration), zap.Int(tm.CodeKey, m.Code))
	fields = append(fields, tz.Meta(ctx)...)

	tz.LogWithFunc(message(fmt.Sprintf("%s %s", method, path)), nil, codeToLevel(m.Code, h.logger), fields...)
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
	if ss.IsHealth(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	path, method := req.URL.Path, strings.ToLower(req.Method)
	start := time.Now()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)
	fields := []zapcore.Field{
		zap.Stringer(tm.DurationKey, time.Since(start)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, path),
		zap.String(tm.MethodKey, method),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if resp != nil {
		fields = append(fields, zap.Int(tm.CodeKey, resp.StatusCode))
	}

	tz.LogWithFunc(message(fmt.Sprintf("%s %s", method, req.URL.Redacted())), err, respToLevel(resp, r.logger), fields...)

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
