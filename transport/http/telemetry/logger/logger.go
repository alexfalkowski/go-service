package logger

import (
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/meta"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/status"
	"github.com/alexfalkowski/go-service/v3/runtime"
	"github.com/alexfalkowski/go-service/v3/telemetry/logger"
	"github.com/alexfalkowski/go-service/v3/time"
	snoop "github.com/felixge/httpsnoop"
)

// Logger is an alias for [github.com/alexfalkowski/go-service/v3/telemetry/logger.Logger].
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
//   - http.request.method and url.path: lowercase HTTP verb and request path
//   - http.route: matched route template, when available
//   - service/method: logical names for registered RPC routes with both names in the path
//   - duration: wall-clock elapsed time
//   - code: HTTP response status code
//   - error: the request diagnostic error, when present
//
// Recovered panics are logged at error level with their original diagnostic error. An aborted committed
// stream (see [http.ErrAbortHandler]) is the one exception: it still produces the ordinary outcome log line
// above, not a panic log, since the response was already committed.
//
// Log level is derived from the status code:
//   - 4xx → warn
//   - 5xx → error
//   - otherwise → info
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	message, attrs := serverRequestLog(req, h.routePolicy)
	ctx := status.WithRequestError(req.Context())
	start := time.Now()

	isOperation := h.routePolicy.IsOperation(req)
	m := snoop.Metrics{Code: http.StatusOK}
	defer func() {
		if value := recover(); value != nil {
			// A streaming handler that already committed its response aborts via
			// http.ErrAbortHandler (see net/http/content/stream's error contract) so net/http can
			// sever the in-flight connection. That is an intentional, expected outcome of a truncated
			// stream, not a bug, so it must not be logged as a recovered panic — but the request still
			// gets its ordinary outcome log line, carrying the diagnostic error finalizeStream recorded.
			if err, ok := value.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				if !isOperation {
					h.logOutcome(ctx, attrs, message, m)
				}
				panic(value)
			}

			h.logPanic(ctx, attrs, runtime.ConvertRecover(value), time.Since(start).String())
			panic(value)
		}
	}()
	if isOperation {
		next(res, req.WithContext(ctx))
		if err := status.RequestError(ctx); errors.Is(err, runtime.ErrRecovered) {
			h.logPanic(ctx, attrs, err, time.Since(start).String())
		}

		return
	}

	m.CaptureMetrics(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })
	if err := status.RequestError(ctx); errors.Is(err, runtime.ErrRecovered) {
		attrs = append(attrs, logger.String(meta.DurationKey, m.Duration.String()), logger.Int(meta.CodeKey, m.Code))
		h.logPanic(ctx, attrs, err, "")
		return
	}

	h.logOutcome(ctx, attrs, message, m)
}

func (h *Handler) logPanic(ctx context.Context, attrs []logger.Attr, err error, duration string) {
	if duration != "" {
		attrs = append(attrs, logger.String(meta.DurationKey, duration))
	}

	h.logger.LogAttrs(ctx, logger.LevelError, logger.NewMessage("http: panic", err), attrs...)
}

func (h *Handler) logOutcome(ctx context.Context, attrs []logger.Attr, message string, m snoop.Metrics) {
	attrs = append(attrs, logger.String(meta.DurationKey, m.Duration.String()), logger.Int(meta.CodeKey, m.Code))
	h.logger.LogAttrs(ctx, codeToLevel(m.Code), logger.NewMessage(message, status.RequestError(ctx)), attrs...)
}

// NewRoundTripper constructs HTTP client logging middleware.
//
// The returned RoundTripper logs request outcomes (duration and status) after the underlying transport returns.
func NewRoundTripper(logger *Logger, r http.RoundTripper) *RoundTripper {
	return &RoundTripper{logger: logger, RoundTripper: r}
}

// RoundTripper logs HTTP client requests and responses.
type RoundTripper struct {
	logger *Logger
	http.RoundTripper
}

// RoundTrip delegates to the underlying RoundTripper and logs the request outcome.
//
// Logged attributes include:
//   - system: "http"
//   - http.request.method, server.address and url.path: lowercase HTTP verb, destination hostname and request path
//   - url.template: path template supplied by the caller, when available
//   - service/method: logical names for explicitly identified RPC calls
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
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	message, attrs := clientRequestLog(req)
	start := time.Now()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)

	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
	code := responseCode(resp, err)
	attrs = append(attrs, logger.Int(meta.CodeKey, code))

	r.logger.LogAttrs(ctx, codeToLevel(code), logger.NewMessage(message, err), attrs...)
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
