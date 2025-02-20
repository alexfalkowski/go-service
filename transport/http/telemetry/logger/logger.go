package logger

import (
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	snoop "github.com/felixge/httpsnoop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const service = "http"

// NewHandler for logger.
func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{logger: logger}
}

// Handler for logger.
type Handler struct {
	logger *logger.Logger
}

// ServeHTTP or logger.
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

	h.logger.LogFunc(ctx, codeToLevel(m.Code, h.logger), message(method+" "+path), nil, fields...)
}

// NewRoundTripper for logger.
func NewRoundTripper(logger *logger.Logger, r http.RoundTripper) *RoundTripper {
	return &RoundTripper{logger: logger, RoundTripper: r}
}

// RoundTripper for logger.
type RoundTripper struct {
	logger *logger.Logger

	http.RoundTripper
}

// RoundTrip for logger.
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

	if resp != nil {
		fields = append(fields, zap.Int(meta.CodeKey, resp.StatusCode))
	}

	r.logger.LogFunc(ctx, respToLevel(resp, r.logger), message(method+" "+req.URL.Redacted()), err, fields...)

	return resp, err
}

func respToLevel(resp *http.Response, logger *logger.Logger) func(msg string, fields ...zapcore.Field) {
	var code int

	if resp != nil {
		code = resp.StatusCode
	} else {
		code = 500
	}

	return codeToLevel(code, logger)
}

func codeToLevel(code int, logger *logger.Logger) func(msg string, fields ...zapcore.Field) {
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
