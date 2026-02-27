package logger

import (
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	snoop "github.com/felixge/httpsnoop"
)

// Logger is an alias for `telemetry/logger.Logger`.
//
// It is re-exported here so transport-layer code can depend on a single logger type when composing
// middleware.
type Logger = logger.Logger

// NewHandler constructs HTTP server logging middleware.
//
// The returned handler logs the outcome of each request after next has completed, including duration
// and response status code. Ignorable paths (health/metrics/etc.) are skipped.
func NewHandler(logger *Logger) *Handler {
	return &Handler{logger: logger}
}

// Handler logs HTTP server requests and responses.
type Handler struct {
	logger *Logger
}

// ServeHTTP logs the request outcome after next completes.
//
// Ignorable paths (health/metrics/etc.) bypass logging (see `transport/strings.IsIgnorable`).
//
// Logged attributes include:
//   - system: "http"
//   - service/method: derived from the request (see `http.ParseServiceMethod`)
//   - duration: wall-clock elapsed time
//   - code: HTTP response status code
//
// Log level is derived from the status code:
//   - 4xx → warn
//   - 5xx → error
//   - otherwise → info
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

// NewRoundTripper constructs HTTP client logging middleware.
//
// The returned RoundTripper logs request outcomes (duration and status) and then delegates to the
// underlying transport. Ignorable paths (health/metrics/etc.) are skipped.
func NewRoundTripper(logger *Logger, r http.RoundTripper) *RoundTripper {
	return &RoundTripper{logger: logger, RoundTripper: r}
}

// RoundTripper logs HTTP client requests and responses.
type RoundTripper struct {
	logger *Logger
	http.RoundTripper
}

// RoundTrip logs the request outcome and delegates to the underlying RoundTripper.
//
// Ignorable paths (health/metrics/etc.) bypass logging (see `transport/strings.IsIgnorable`).
//
// Logged attributes include:
//   - system: "http"
//   - service/method: derived from the request (see `http.ParseServiceMethod`)
//   - duration: wall-clock elapsed time
//   - code: HTTP response status code (when available)
//
// Log level is derived from the status code:
//   - 4xx → warn
//   - 5xx → error
//   - otherwise → info
//
// If resp is nil (for example, due to a transport error), it is treated as HTTP 500 for level selection.
// The log message includes the derived method and service.
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
