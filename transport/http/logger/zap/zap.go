package zap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	httpURL        = "http.url"
	httpMethod     = "http.method"
	httpDuration   = "http.duration_ms"
	httpStartTime  = "http.start_time"
	httpDeadline   = "http.deadline"
	httpStatusCode = "http.status_code"
	component      = "component"
	httpComponent  = "http"
	client         = "client"
	server         = "server"
)

// HandlerParams for zap.
type HandlerParams struct {
	Logger  *zap.Logger
	Handler http.Handler
}

// NewHandler for zap.
func NewHandler(params HandlerParams) *Handler {
	return &Handler{logger: params.Logger, Handler: params.Handler}
}

// Handler for zap.
type Handler struct {
	logger *zap.Logger
	http.Handler
}

// ServeHTTP  or zap.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	service, method := req.URL.Path, strings.ToLower(req.Method)
	if tstrings.IsHealth(service) {
		h.Handler.ServeHTTP(resp, req)

		return
	}

	start := time.Now().UTC()
	ctx := req.Context()

	h.Handler.ServeHTTP(resp, req)

	fields := []zapcore.Field{
		zap.Int64(httpDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(httpStartTime, start.Format(time.RFC3339)),
		zap.String(httpURL, service),
		zap.String(httpMethod, method),
		zap.String("span.kind", server),
		zap.String(component, httpComponent),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(httpDeadline, d.UTC().Format(time.RFC3339)))
	}

	h.logger.Info("finished call", fields...)
}

// RoundTripperParams for zap.
type RoundTripperParams struct {
	Logger       *zap.Logger
	RoundTripper http.RoundTripper
}

// NewRoundTripper for zap.
func NewRoundTripper(params RoundTripperParams) *RoundTripper {
	return &RoundTripper{logger: params.Logger, RoundTripper: params.RoundTripper}
}

// RoundTripper for zap.
type RoundTripper struct {
	logger *zap.Logger

	http.RoundTripper
}

// RoundTrip for zap.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if tstrings.IsHealth(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := req.URL.Hostname(), strings.ToLower(req.Method)
	start := time.Now().UTC()
	ctx := req.Context()
	resp, err := r.RoundTripper.RoundTrip(req)
	fields := []zapcore.Field{
		zap.Int64(httpDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(httpStartTime, start.Format(time.RFC3339)),
		zap.String(httpURL, service),
		zap.String(httpMethod, method),
		zap.String("span.kind", client),
		zap.String(component, httpComponent),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(httpDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		r.logger.Error("finished call with error", fields...)

		return nil, err
	}

	fields = append(fields, zap.Int(httpStatusCode, resp.StatusCode))

	r.logger.Info(fmt.Sprintf("finished call with code %s", resp.Status), fields...)

	return resp, nil
}
