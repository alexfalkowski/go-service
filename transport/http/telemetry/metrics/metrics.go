package metrics

import (
	"net/http"
	"strings"
	"time"

	shttp "github.com/alexfalkowski/go-service/http"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// NewHandler for metrics.
func NewHandler(meter metric.Meter) (*Handler, error) {
	started, err := meter.Int64Counter("http_server_started_total", metric.WithDescription("Total number of RPCs started on the server."))
	if err != nil {
		return nil, err
	}

	received, err := meter.Int64Counter("http_server_msg_received_total", metric.WithDescription("Total number of RPC messages received on the server."))
	if err != nil {
		return nil, err
	}

	sent, err := meter.Int64Counter("http_server_msg_sent_total", metric.WithDescription("Total number of RPC messages sent by the server."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Int64Counter("http_server_handled_total",
		metric.WithDescription("Total number of RPCs completed on the server, regardless of success or failure."))
	if err != nil {
		return nil, err
	}

	handledHist, err := meter.Float64Histogram("http_server_handling_seconds",
		metric.WithDescription("Histogram of response latency (seconds) of HTTP that had been application-level handled by the server."))
	if err != nil {
		return nil, err
	}

	h := &Handler{
		started: started, received: received,
		sent: sent, handled: handled,
		handledHist: handledHist,
	}

	return h, nil
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
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	service, method := req.URL.Path, strings.ToLower(req.Method)
	if tstrings.IsHealth(service) {
		next(resp, req)

		return
	}

	opts := metric.WithAttributes(
		attribute.Key("http_service").String(service),
		attribute.Key("http_method").String(method),
	)

	st := time.Now()
	ctx := req.Context()

	h.started.Add(ctx, 1, opts)
	h.received.Add(ctx, 1, opts)

	res := &shttp.ResponseWriter{ResponseWriter: resp, StatusCode: http.StatusOK}
	next(res, req)

	h.handled.Add(ctx, 1, opts, metric.WithAttributes(attribute.Key("http_code").Int(res.StatusCode)))
	h.handledHist.Record(ctx, time.Since(st).Seconds(), opts)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		h.sent.Add(ctx, 1, opts)
	}
}

// NewRoundTripper for metrics.
func NewRoundTripper(meter metric.Meter, r http.RoundTripper) (*RoundTripper, error) {
	started, err := meter.Int64Counter("http_client_started_total", metric.WithDescription("Total number of RPCs started on the client."))
	if err != nil {
		return nil, err
	}

	received, err := meter.Int64Counter("http_client_msg_received_total", metric.WithDescription("Total number of RPC messages received on the client."))
	if err != nil {
		return nil, err
	}

	sent, err := meter.Int64Counter("http_client_msg_sent_total", metric.WithDescription("Total number of RPC messages sent by the client."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Int64Counter("http_client_handled_total",
		metric.WithDescription("Total number of RPCs completed on the client, regardless of success or failure."))
	if err != nil {
		return nil, err
	}

	handledHist, err := meter.Float64Histogram("http_client_handling_seconds",
		metric.WithDescription("Histogram of response latency (seconds) of HTTP that had been application-level handled by the client."))
	if err != nil {
		return nil, err
	}

	rt := &RoundTripper{
		started: started, received: received, sent: sent, handled: handled, handledHist: handledHist,
		RoundTripper: r,
	}

	return rt, nil
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
	if tstrings.IsHealth(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := req.URL.Hostname(), strings.ToLower(req.Method)
	st := time.Now()
	ctx := req.Context()

	opts := metric.WithAttributes(
		attribute.Key("http_service").String(service),
		attribute.Key("http_method").String(method),
	)

	r.started.Add(ctx, 1, opts)
	r.sent.Add(ctx, 1, opts)

	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	r.received.Add(ctx, 1, opts)
	r.handled.Add(ctx, 1, opts, metric.WithAttributes(attribute.Key("http_code").Int(resp.StatusCode)))
	r.handledHist.Record(ctx, time.Since(st).Seconds(), opts)

	return resp, nil
}
