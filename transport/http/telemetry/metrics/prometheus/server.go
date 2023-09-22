package prometheus

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/os"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

// ServerCollector for prometheus.
type ServerCollector struct {
	started          *prometheus.CounterVec
	handled          *prometheus.CounterVec
	received         *prometheus.CounterVec
	sent             *prometheus.CounterVec
	handledHistogram *prometheus.HistogramVec
}

// NewServerCollector for prometheus.
//
//nolint:dupl
func NewServerCollector(lc fx.Lifecycle, version version.Version) *ServerCollector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ServerCollector{
		started: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_started_total",
				Help:        "Total number of RPCs started on the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		handled: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_handled_total",
				Help:        "Total number of RPCs completed on the server, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method", "http_code"}),
		received: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_msg_received_total",
				Help:        "Total number of RPC messages received on the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		sent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_msg_sent_total",
				Help:        "Total number of RPC messages sent by the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		handledHistogram: prometheus.NewHistogramVec(
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
func (m *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	m.started.Describe(ch)
	m.handled.Describe(ch)
	m.received.Describe(ch)
	m.sent.Describe(ch)
	m.handledHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	m.started.Collect(ch)
	m.handled.Collect(ch)
	m.received.Collect(ch)
	m.sent.Collect(ch)
	m.handledHistogram.Collect(ch)
}

// Handler for prometheus.
func (m *ServerCollector) Handler(h http.Handler) http.Handler {
	return &handler{metrics: m, Handler: h}
}

type handler struct {
	metrics *ServerCollector
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
	metrics     *ServerCollector
	serviceName string
	methodName  string
	startTime   time.Time
}

func newServerReporter(m *ServerCollector, service, method string) *serverReporter {
	r := &serverReporter{metrics: m, startTime: time.Now(), serviceName: service, methodName: method}
	r.metrics.started.WithLabelValues(r.serviceName, r.methodName).Inc()

	return r
}

func (r *serverReporter) ReceivedMessage() {
	r.metrics.received.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *serverReporter) SentMessage() {
	r.metrics.sent.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *serverReporter) Handled(code int) {
	r.metrics.handled.WithLabelValues(r.serviceName, r.methodName, strconv.Itoa(code)).Inc()
	r.metrics.handledHistogram.WithLabelValues(r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
}
