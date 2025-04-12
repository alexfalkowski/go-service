package logger

import (
	"log/slog"
	"net/http"

	"github.com/alexfalkowski/go-service/strings"
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
	attrs := []slog.Attr{
		slog.String(meta.ServiceKey, service),
		slog.String(meta.PathKey, path),
		slog.String(meta.MethodKey, method),
	}

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })

	attrs = append(attrs, slog.String(meta.DurationKey, m.Duration.String()), slog.Int(meta.CodeKey, m.Code))
	h.logger.LogAttrs(ctx, codeToLevel(m.Code), logger.NewText(message(method+" "+path)), attrs...)
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
	attrs := []slog.Attr{
		slog.String(meta.DurationKey, time.Since(start).String()),
		slog.String(meta.ServiceKey, service),
		slog.String(meta.PathKey, path),
		slog.String(meta.MethodKey, method),
	}

	if resp != nil {
		attrs = append(attrs, slog.Int(meta.CodeKey, resp.StatusCode))
	}

	r.logger.LogAttrs(ctx, respToLevel(resp), logger.NewMessage(message(method+" "+path), err), attrs...)

	return resp, err
}

func respToLevel(resp *http.Response) slog.Level {
	var code int

	if resp != nil {
		code = resp.StatusCode
	} else {
		code = 500
	}

	return codeToLevel(code)
}

func codeToLevel(code int) slog.Level {
	if code >= 400 && code <= 499 {
		return slog.LevelWarn
	}

	if code >= 500 && code <= 599 {
		return slog.LevelError
	}

	return slog.LevelInfo
}

func message(msg string) string {
	return "http: " + msg
}
