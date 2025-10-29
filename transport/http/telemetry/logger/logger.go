package logger

import (
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	snoop "github.com/felixge/httpsnoop"
)

// Logger is an alias of logger.Logger.
type Logger = logger.Logger

// NewHandler for logger.
func NewHandler(logger *Logger) *Handler {
	return &Handler{logger: logger}
}

// Handler for logger.
type Handler struct {
	logger *Logger
}

// ServeHTTP or logger.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	p, method := http.Path(req), strings.ToLower(req.Method)
	if strings.IsObservable(p) {
		next(res, req)
		return
	}

	ctx := req.Context()
	attrs := []logger.Attr{
		logger.String(meta.ServiceKey, "http"),
		logger.String(meta.PathKey, p),
		logger.String(meta.MethodKey, method),
	}
	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })
	attrs = append(attrs, logger.String(meta.DurationKey, m.Duration.String()), logger.Int(meta.CodeKey, m.Code))
	message := logger.NewText(message(strings.Join(strings.Space, method, p)))

	h.logger.LogAttrs(ctx, codeToLevel(m.Code), message, attrs...)
}

// NewRoundTripper for logger.
func NewRoundTripper(logger *Logger, r http.RoundTripper) *RoundTripper {
	return &RoundTripper{logger: logger, RoundTripper: r}
}

// RoundTripper for logger.
type RoundTripper struct {
	logger *Logger
	http.RoundTripper
}

// RoundTrip for logger.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	p, method := http.Path(req), strings.ToLower(req.Method)
	if strings.IsObservable(p) {
		return r.RoundTripper.RoundTrip(req)
	}

	start := time.Now()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)
	attrs := []logger.Attr{
		logger.String(meta.DurationKey, time.Since(start).String()),
		logger.String(meta.ServiceKey, "http"),
		logger.String(meta.PathKey, p),
		logger.String(meta.MethodKey, method),
	}
	if resp != nil {
		attrs = append(attrs, logger.Int(meta.CodeKey, resp.StatusCode))
	}
	message := logger.NewText(message(strings.Join(strings.Space, method, p)))

	r.logger.LogAttrs(ctx, respToLevel(resp), message, attrs...)
	return resp, err
}

func respToLevel(resp *http.Response) logger.Level {
	var code int
	if resp != nil {
		code = resp.StatusCode
	} else {
		code = 500
	}

	return codeToLevel(code)
}

func codeToLevel(code int) logger.Level {
	if code >= 400 && code <= 499 {
		return logger.LevelWarn
	}

	if code >= 500 && code <= 599 {
		return logger.LevelError
	}

	return logger.LevelInfo
}

func message(msg string) string {
	return "http: " + msg
}
