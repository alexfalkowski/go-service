package logger

import (
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	snoop "github.com/felixge/httpsnoop"
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
	fields := []logger.Field{
		logger.String(meta.ServiceKey, service),
		logger.String(meta.PathKey, path),
		logger.String(meta.MethodKey, method),
	}

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })

	fields = append(fields, logger.Stringer(meta.DurationKey, m.Duration), logger.Int(meta.CodeKey, m.Code))

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
	fields := []logger.Field{
		logger.Stringer(meta.DurationKey, time.Since(start)),
		logger.String(meta.ServiceKey, service),
		logger.String(meta.PathKey, path),
		logger.String(meta.MethodKey, method),
	}

	if resp != nil {
		fields = append(fields, logger.Int(meta.CodeKey, resp.StatusCode))
	}

	r.logger.LogFunc(ctx, respToLevel(resp, r.logger), message(method+" "+req.URL.Redacted()), err, fields...)

	return resp, err
}

func respToLevel(resp *http.Response, logger *logger.Logger) func(msg string, fields ...logger.Field) {
	var code int

	if resp != nil {
		code = resp.StatusCode
	} else {
		code = 500
	}

	return codeToLevel(code, logger)
}

func codeToLevel(code int, logger *logger.Logger) func(msg string, fields ...logger.Field) {
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
