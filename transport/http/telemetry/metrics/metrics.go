package metrics

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	snoop "github.com/felixge/httpsnoop"
	prometheus "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	nameAttribute   = attributes.Key("service_name")
	pathAttribute   = attributes.Key("service_path")
	methodAttribute = attributes.Key("service_method")
	codeAttribute   = attributes.Key("service_code")
)

// Meter is an alias for metrics.Meter.
type Meter = metrics.Meter

// Register prometheus.
func Register(name env.Name, cfg *metrics.Config, mux *http.ServeMux) {
	if cfg.IsEnabled() && cfg.IsPrometheus() {
		mux.Handle("GET "+http.Pattern(name, "/metrics"), prometheus.Handler())
	}
}

// NewHandler for metrics.
func NewHandler(name env.Name, meter *Meter) *Handler {
	started := meter.MustInt64Counter("http_server_started_total", "Total number of RPCs started on the server.")
	received := meter.MustInt64Counter("http_server_msg_received_total", "Total number of RPC messages received on the server.")
	sent := meter.MustInt64Counter("http_server_msg_sent_total", "Total number of RPC messages sent by the server.")
	handled := meter.MustInt64Counter("http_server_handled_total", "Total number of RPCs completed on the server, regardless of success or failure.")
	handledHist := meter.MustFloat64Histogram("http_server_handling_seconds",
		"Histogram of response latency (seconds) of HTTP that had been application-level handled by the server.")

	return &Handler{
		name:        name,
		started:     started,
		received:    received,
		sent:        sent,
		handled:     handled,
		handledHist: handledHist,
	}
}

// Handler for metrics.
type Handler struct {
	started     metrics.Int64Counter
	received    metrics.Int64Counter
	sent        metrics.Int64Counter
	handled     metrics.Int64Counter
	handledHist metrics.Float64Histogram
	name        env.Name
}

// ServeHTTP for metrics.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsObservable(req.URL.Path) {
		next(res, req)
		return
	}

	service, method := http.ParseServiceMethod(req)
	opts := metrics.WithAttributes(
		nameAttribute.String(h.name.String()),
		pathAttribute.String(service),
		methodAttribute.String(method),
	)
	start := time.Now()
	ctx := req.Context()

	h.started.Add(ctx, 1, opts)
	h.received.Add(ctx, 1, opts)

	captured := snoop.CaptureMetricsFn(res, func(res http.ResponseWriter) { next(res, req.WithContext(ctx)) })
	if captured.Code >= 200 && captured.Code <= 299 {
		h.sent.Add(ctx, 1, opts)
	}

	h.handled.Add(ctx, 1, opts, metrics.WithAttributes(codeAttribute.String(strconv.Itoa(captured.Code))))
	h.handledHist.Record(ctx, time.Since(start).Seconds(), opts)
}

// NewRoundTripper for metrics.
func NewRoundTripper(name env.Name, meter *Meter, r http.RoundTripper) *RoundTripper {
	started := meter.MustInt64Counter("http_client_started_total", "Total number of RPCs started on the client.")
	received := meter.MustInt64Counter("http_client_msg_received_total", "Total number of RPC messages received on the client.")
	sent := meter.MustInt64Counter("http_client_msg_sent_total", "Total number of RPC messages sent by the client.")
	handled := meter.MustInt64Counter("http_client_handled_total", "Total number of RPCs completed on the client, regardless of success or failure.")
	handledHist := meter.MustFloat64Histogram("http_client_handling_seconds",
		"Histogram of response latency (seconds) of HTTP that had been application-level handled by the client.")

	return &RoundTripper{
		name:         name,
		started:      started,
		received:     received,
		sent:         sent,
		handled:      handled,
		handledHist:  handledHist,
		RoundTripper: r,
	}
}

// RoundTripper for metrics.
type RoundTripper struct {
	started     metrics.Int64Counter
	received    metrics.Int64Counter
	sent        metrics.Int64Counter
	handled     metrics.Int64Counter
	handledHist metrics.Float64Histogram
	http.RoundTripper
	name env.Name
}

// RoundTrip for metrics.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.IsObservable(req.URL.Path) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := http.ParseServiceMethod(req)
	start := time.Now()
	ctx := req.Context()
	opts := metrics.WithAttributes(
		nameAttribute.String(r.name.String()),
		pathAttribute.String(service),
		methodAttribute.String(method),
	)

	r.started.Add(ctx, 1, opts)
	r.sent.Add(ctx, 1, opts)

	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	r.received.Add(ctx, 1, opts)
	r.handled.Add(ctx, 1, opts, metrics.WithAttributes(codeAttribute.String(strconv.Itoa(resp.StatusCode))))
	r.handledHist.Record(ctx, time.Since(start).Seconds(), opts)

	return resp, nil
}
