package zap

import (
	"net/http"
	"strings"
	"time"

	sh "github.com/alexfalkowski/go-service/net/http"
	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	st "github.com/alexfalkowski/go-service/time"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	ss "github.com/alexfalkowski/go-service/transport/strings"
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

	start := time.Now()
	ctx := req.Context()

	res := &sh.ResponseWriter{ResponseWriter: resp, StatusCode: http.StatusOK}
	next(res, req)

	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, path),
		zap.String(tm.MethodKey, method),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	fields = append(fields, zap.Int(tm.CodeKey, res.StatusCode))

	loggerLevel := codeToLevel(res.StatusCode, h.logger)
	loggerLevel("finished call with code "+res.Status(), fields...)
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
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, path),
		zap.String(tm.MethodKey, method),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		r.logger.Error("finished call with error", fields...)

		return nil, err
	}

	fields = append(fields, zap.Int(tm.CodeKey, resp.StatusCode))

	loggerLevel := codeToLevel(resp.StatusCode, r.logger)
	loggerLevel("finished call with code "+resp.Status, fields...)

	return resp, nil
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
