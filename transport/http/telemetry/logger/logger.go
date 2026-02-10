package logger

import (
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	snoop "github.com/felixge/httpsnoop"
)

// Logger is an alias for logger.Logger.
type Logger = logger.Logger

// NewHandler constructs an HTTP server logging handler.
func NewHandler(logger *Logger) *Handler {
	return &Handler{logger: logger}
}

// Handler logs HTTP server requests and responses.
type Handler struct {
	logger *Logger
}

// ServeHTTP logs the request after next completes.
//
// It logs attributes including system ("http"), service/method (derived from the request),
// duration, and response code. Log level is derived from the status code:
// 4xx -> warn, 5xx -> error, otherwise -> info. Ignorable paths bypass logging.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsIgnorable(req.URL.Path) {
		next(res, req)
		return
	}

	service, method := http.ParseServiceMethod(req)
	ctx := req.Context()

	attrs := make([]logger.Attr, 0, 5)
	attrs = append(attrs, logger.String(meta.SystemKey, "http"))
	attrs = append(attrs, logger.String(meta.ServiceKey, service))
	attrs = append(attrs, logger.String(meta.MethodKey, method))

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })
	attrs = append(attrs, logger.String(meta.DurationKey, m.Duration.String()), logger.Int(meta.CodeKey, m.Code))
	message := logger.NewText(message(strings.Join(strings.Space, method, service)))

	h.logger.LogAttrs(ctx, codeToLevel(m.Code), message, attrs...)
}

// NewRoundTripper constructs an HTTP client logging RoundTripper.
func NewRoundTripper(logger *Logger, r http.RoundTripper) *RoundTripper {
	return &RoundTripper{logger: logger, RoundTripper: r}
}

// RoundTripper logs HTTP client requests and responses.
type RoundTripper struct {
	logger *Logger
	http.RoundTripper
}

// RoundTrip logs the request/response and delegates to the underlying RoundTripper.
//
// It logs attributes including system ("http"), service/method (derived from the request),
// duration, and (when available) response code. Log level is derived from the status code:
// 4xx -> warn, 5xx -> error, otherwise -> info. If resp is nil, it is treated as a 500 for level selection.
// Ignorable paths bypass logging.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.IsIgnorable(req.URL.Path) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := http.ParseServiceMethod(req)
	start := time.Now()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)

	attrs := make([]logger.Attr, 0, 5)
	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
	attrs = append(attrs, logger.String(meta.SystemKey, "http"))
	attrs = append(attrs, logger.String(meta.ServiceKey, service))
	attrs = append(attrs, logger.String(meta.MethodKey, method))
	if resp != nil {
		attrs = append(attrs, logger.Int(meta.CodeKey, resp.StatusCode))
	}

	message := logger.NewText(message(strings.Join(strings.Space, method, service)))

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
