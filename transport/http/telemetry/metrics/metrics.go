package metrics

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/time"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	snoop "github.com/felixge/httpsnoop"
	prometheus "github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/metric"
)

// Register prometheus.
func Register(cfg *metrics.Config, mux *http.ServeMux) {
	if !metrics.IsEnabled(cfg) || !cfg.IsPrometheus() {
		return
	}

	mux.Handle("GET /metrics", prometheus.Handler())
}

// NewHandler for metrics.
func NewHandler(meter *metrics.Meter) *Handler {
	started := meter.MustInt64Counter("http_server_started_total", "Total number of RPCs started on the server.")
	received := meter.MustInt64Counter("http_server_msg_received_total", "Total number of RPC messages received on the server.")
	sent := meter.MustInt64Counter("http_server_msg_sent_total", "Total number of RPC messages sent by the server.")
	handled := meter.MustInt64Counter("http_server_handled_total", "Total number of RPCs completed on the server, regardless of success or failure.")
	handledHist := meter.MustFloat64Histogram("http_server_handling_seconds",
		"Histogram of response latency (seconds) of HTTP that had been application-level handled by the server.")

	return &Handler{
		started: started, received: received,
		sent: sent, handled: handled,
		handledHist: handledHist,
	}
}

// Handler for metrics.
type Handler struct {
	started     metric.Int64Counter
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram
}

// ServeHTTP for metrics.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	service, method := req.URL.Path, strings.ToLower(req.Method)
	if ts.IsObservable(service) {
		next(res, req)

		return
	}

	opts := metric.WithAttributes(
		serviceAttribute.String(service),
		methodAttribute.String(method),
	)

	start := time.Now()
	ctx := req.Context()

	h.started.Add(ctx, 1, opts)
	h.received.Add(ctx, 1, opts)

	metrics := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })

	if metrics.Code >= 200 && metrics.Code <= 299 {
		h.sent.Add(ctx, 1, opts)
	}

	h.handled.Add(ctx, 1, opts, metric.WithAttributes(statusCodeAttribute.String(strconv.Itoa(metrics.Code))))
	h.handledHist.Record(ctx, time.Since(start).Seconds(), opts)
}

// NewRoundTripper for metrics.
func NewRoundTripper(meter *metrics.Meter, r http.RoundTripper) *RoundTripper {
	started := meter.MustInt64Counter("http_client_started_total", "Total number of RPCs started on the client.")
	received := meter.MustInt64Counter("http_client_msg_received_total", "Total number of RPC messages received on the client.")
	sent := meter.MustInt64Counter("http_client_msg_sent_total", "Total number of RPC messages sent by the client.")
	handled := meter.MustInt64Counter("http_client_handled_total", "Total number of RPCs completed on the client, regardless of success or failure.")
	handledHist := meter.MustFloat64Histogram("http_client_handling_seconds",
		"Histogram of response latency (seconds) of HTTP that had been application-level handled by the client.")

	return &RoundTripper{
		started: started, received: received, sent: sent, handled: handled, handledHist: handledHist,
		RoundTripper: r,
	}
}

// RoundTripper for metrics.
type RoundTripper struct {
	started     metric.Int64Counter
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram

	http.RoundTripper
}

// RoundTrip for metrics.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if ts.IsObservable(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := req.URL.Hostname(), strings.ToLower(req.Method)
	start := time.Now()
	ctx := req.Context()

	opts := metric.WithAttributes(
		serviceAttribute.String(service),
		methodAttribute.String(method),
	)

	r.started.Add(ctx, 1, opts)
	r.sent.Add(ctx, 1, opts)

	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	r.received.Add(ctx, 1, opts)
	r.handled.Add(ctx, 1, opts, metric.WithAttributes(statusCodeAttribute.String(strconv.Itoa(resp.StatusCode))))
	r.handledHist.Record(ctx, time.Since(start).Seconds(), opts)

	return resp, nil
}
