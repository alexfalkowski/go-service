package logger

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	snoop "github.com/felixge/httpsnoop"
)

// Logger is an alias for [github.com/alexfalkowski/go-service/v2/telemetry/logger.Logger].
//
// It is re-exported here so transport-layer code can depend on a single logger type when composing
// middleware.
type Logger = logger.Logger

// NewHandler constructs HTTP server logging middleware.
//
// The returned handler logs the outcome of each request after next has completed, including duration and
// response status code. Registered operation paths (health/metrics/etc.) skip ordinary access logs but still
// log recovered panics.
func NewHandler(routePolicy *http.RoutePolicy, logger *Logger) *Handler {
	return &Handler{routePolicy: routePolicy, logger: logger}
}

// Handler logs HTTP server requests and responses.
type Handler struct {
	logger      *Logger
	routePolicy *http.RoutePolicy
}

// ServeHTTP logs the request outcome after next completes.
//
// Registered operation paths (health/metrics/etc.) bypass ordinary access logging but retain recovered-panic
// logging.
//
// Logged attributes include:
//   - system: "http"
//   - service/method: derived from the request (see [http.ParseServiceMethod])
//   - duration: wall-clock elapsed time
//   - code: HTTP response status code
//   - error: the request diagnostic error, when present
//
// Recovered panics are logged at error level with their original diagnostic error.
//
// Log level is derived from the status code:
//   - 4xx → warn
//   - 5xx → error
//   - otherwise → info
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	service, method := http.ParseServiceMethod(req)
	ctx := status.WithRequestError(req.Context())
	start := time.Now()

	attrs := make([]logger.Attr, 0, 5)
	attrs = append(attrs, logger.String(meta.SystemKey, "http"))
	attrs = append(attrs, logger.String(meta.ServiceKey, service))
	attrs = append(attrs, logger.String(meta.MethodKey, method))
	defer func() {
		if value := recover(); value != nil {
			h.logPanic(ctx, attrs, runtime.ConvertRecover(value), time.Since(start).String())
			panic(value)
		}
	}()
	if h.routePolicy.IsOperation(req) {
		next(res, req.WithContext(ctx))
		if err := status.RequestError(ctx); errors.Is(err, runtime.ErrRecovered) {
			h.logPanic(ctx, attrs, err, time.Since(start).String())
		}

		return
	}

	m := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })
	attrs = append(attrs, logger.String(meta.DurationKey, m.Duration.String()), logger.Int(meta.CodeKey, m.Code))
	if err := status.RequestError(ctx); errors.Is(err, runtime.ErrRecovered) {
		h.logPanic(ctx, attrs, err, "")
		return
	}

	message := logger.NewMessage(httpMessage(strings.Join(strings.Space, method, service)), status.RequestError(ctx))

	h.logger.LogAttrs(ctx, codeToLevel(m.Code), message, attrs...)
}

func (h *Handler) logPanic(ctx context.Context, attrs []logger.Attr, err error, duration string) {
	if duration != "" {
		attrs = append(attrs, logger.String(meta.DurationKey, duration))
	}

	h.logger.LogAttrs(ctx, logger.LevelError, logger.NewMessage("http: panic", err), attrs...)
}

// NewRoundTripper constructs HTTP client logging middleware.
//
// The returned RoundTripper logs request outcomes (duration and status) and then delegates to the
// underlying transport.
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
// Logged attributes include:
//   - system: "http"
//   - service/method: derived from the request (see [http.ParseServiceMethod])
//   - duration: wall-clock elapsed time
//   - code: HTTP response status code, or status-bearing transport error code when no response is available
//   - error: transport error (when present)
//
// Log level is derived from the status code:
//   - 4xx → warn
//   - 5xx → error
//   - otherwise → info
//
// If resp is nil and err does not carry a status code, it is treated as HTTP 500 for level selection.
// The log message includes the derived method and service.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	service, method := http.ParseServiceMethod(req)
	start := time.Now()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)

	attrs := make([]logger.Attr, 0, 5)
	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
	attrs = append(attrs, logger.String(meta.SystemKey, "http"))
	attrs = append(attrs, logger.String(meta.ServiceKey, service))
	attrs = append(attrs, logger.String(meta.MethodKey, method))
	code := responseCode(resp, err)
	attrs = append(attrs, logger.Int(meta.CodeKey, code))

	message := logger.NewMessage(httpMessage(strings.Join(strings.Space, method, service)), err)

	r.logger.LogAttrs(ctx, codeToLevel(code), message, attrs...)
	return resp, err
}

func responseCode(resp *http.Response, err error) int {
	if resp != nil {
		return resp.StatusCode
	}

	return status.Code(err)
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

func httpMessage(msg string) string {
	return "http: " + msg
}
