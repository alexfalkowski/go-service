package telemetry

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/os"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

// ClientMetrics for prometheus.
type ClientMetrics struct {
	clientStartedCounter   *prometheus.CounterVec
	clientHandledCounter   *prometheus.CounterVec
	clientMsgReceived      *prometheus.CounterVec
	clientMsgSent          *prometheus.CounterVec
	clientHandledHistogram *prometheus.HistogramVec
}

// NewClientMetrics for prometheus.
//
//nolint:dupl
func NewClientMetrics(lc fx.Lifecycle, version version.Version) *ClientMetrics {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ClientMetrics{
		clientStartedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_started_total",
				Help:        "Total number of RPCs started on the client.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		clientHandledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_handled_total",
				Help:        "Total number of RPCs completed by the client, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method", "http_code"}),
		clientMsgReceived: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_msg_received_total",
				Help:        "Total number of RPC stream messages received by the client.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		clientMsgSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_msg_sent_total",
				Help:        "Total number of HTTP stream messages sent by the client.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		clientHandledHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "http_client_handling_seconds",
			Help:        "Histogram of response latency (seconds) of the HTTP until it is finished by the application.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: labels,
		}, []string{"http_service", "http_method"}),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(metrics)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(metrics)

			return nil
		},
	})

	return metrics
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ClientMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.clientStartedCounter.Describe(ch)
	m.clientHandledCounter.Describe(ch)
	m.clientMsgReceived.Describe(ch)
	m.clientMsgSent.Describe(ch)
	m.clientHandledHistogram.Describe(ch)
}

// Collect is called by the prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ClientMetrics) Collect(ch chan<- prometheus.Metric) {
	m.clientStartedCounter.Collect(ch)
	m.clientHandledCounter.Collect(ch)
	m.clientMsgReceived.Collect(ch)
	m.clientMsgSent.Collect(ch)
	m.clientHandledHistogram.Collect(ch)
}

// RoundTripper for prometheus.
func (m *ClientMetrics) RoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &roundTripper{Metrics: m, RoundTripper: rt}
}

type roundTripper struct {
	Metrics *ClientMetrics
	http.RoundTripper
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if tstrings.IsHealth(req.URL.String()) {
		return r.RoundTripper.RoundTrip(req)
	}

	service, method := req.URL.Hostname(), strings.ToLower(req.Method)
	monitor := newClientReporter(r.Metrics, service, method)
	monitor.SentMessage()

	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	monitor.ReceivedMessage()
	monitor.Handled(resp.StatusCode)

	return resp, nil
}

type clientReporter struct {
	metrics     *ClientMetrics
	serviceName string
	methodName  string
	startTime   time.Time
}

func newClientReporter(m *ClientMetrics, service, method string) *clientReporter {
	r := &clientReporter{metrics: m, startTime: time.Now(), serviceName: service, methodName: method}
	r.metrics.clientStartedCounter.WithLabelValues(r.serviceName, r.methodName).Inc()

	return r
}

func (r *clientReporter) ReceivedMessage() {
	r.metrics.clientMsgReceived.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *clientReporter) SentMessage() {
	r.metrics.clientMsgSent.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *clientReporter) Handled(code int) {
	r.metrics.clientHandledCounter.WithLabelValues(r.serviceName, r.methodName, strconv.Itoa(code)).Inc()
	r.metrics.clientHandledHistogram.WithLabelValues(r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
}

// ServerMetrics for prometheus.
type ServerMetrics struct {
	serverStartedCounter   *prometheus.CounterVec
	serverHandledCounter   *prometheus.CounterVec
	serverMsgReceived      *prometheus.CounterVec
	serverMsgSent          *prometheus.CounterVec
	serverHandledHistogram *prometheus.HistogramVec
}

// NewServerMetrics for prometheus.
//
//nolint:dupl
func NewServerMetrics(lc fx.Lifecycle, version version.Version) *ServerMetrics {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ServerMetrics{
		serverStartedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_started_total",
				Help:        "Total number of RPCs started on the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		serverHandledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_handled_total",
				Help:        "Total number of RPCs completed on the server, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method", "http_code"}),
		serverMsgReceived: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_msg_received_total",
				Help:        "Total number of RPC messages received on the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		serverMsgSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_msg_sent_total",
				Help:        "Total number of RPC messages sent by the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		serverHandledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_server_handling_seconds",
				Help:        "Histogram of response latency (seconds) of HTTP that had been application-level handled by the server.",
				Buckets:     prometheus.DefBuckets,
				ConstLabels: labels,
			}, []string{"http_service", "http_method"},
		),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(metrics)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(metrics)

			return nil
		},
	})

	return metrics
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ServerMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.serverStartedCounter.Describe(ch)
	m.serverHandledCounter.Describe(ch)
	m.serverMsgReceived.Describe(ch)
	m.serverMsgSent.Describe(ch)
	m.serverHandledHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ServerMetrics) Collect(ch chan<- prometheus.Metric) {
	m.serverStartedCounter.Collect(ch)
	m.serverHandledCounter.Collect(ch)
	m.serverMsgReceived.Collect(ch)
	m.serverMsgSent.Collect(ch)
	m.serverHandledHistogram.Collect(ch)
}

// Handler for prometheus.
func (m *ServerMetrics) Handler(h http.Handler) http.Handler {
	return &handler{metrics: m, Handler: h}
}

type handler struct {
	metrics *ServerMetrics
	http.Handler
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	service, method := req.URL.Path, strings.ToLower(req.Method)
	if tstrings.IsHealth(service) {
		h.Handler.ServeHTTP(resp, req)

		return
	}

	monitor := newServerReporter(h.metrics, service, method)
	monitor.ReceivedMessage()

	res := &responseWriter{ResponseWriter: resp, Status: http.StatusOK}
	h.Handler.ServeHTTP(res, req)

	monitor.Handled(res.Status)

	if res.Status >= 200 && res.Status <= 299 {
		monitor.SentMessage()
	}
}

type responseWriter struct {
	http.ResponseWriter
	Status int
}

func (r *responseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type serverReporter struct {
	metrics     *ServerMetrics
	serviceName string
	methodName  string
	startTime   time.Time
}

func newServerReporter(m *ServerMetrics, service, method string) *serverReporter {
	r := &serverReporter{metrics: m, startTime: time.Now(), serviceName: service, methodName: method}
	r.metrics.serverStartedCounter.WithLabelValues(r.serviceName, r.methodName).Inc()

	return r
}

func (r *serverReporter) ReceivedMessage() {
	r.metrics.serverMsgReceived.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *serverReporter) SentMessage() {
	r.metrics.serverMsgSent.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *serverReporter) Handled(code int) {
	r.metrics.serverHandledCounter.WithLabelValues(r.serviceName, r.methodName, strconv.Itoa(code)).Inc()
	r.metrics.serverHandledHistogram.WithLabelValues(r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
}
