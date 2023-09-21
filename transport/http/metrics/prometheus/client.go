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

// ClientCollector for prometheus.
type ClientCollector struct {
	started          *prometheus.CounterVec
	handled          *prometheus.CounterVec
	received         *prometheus.CounterVec
	sent             *prometheus.CounterVec
	handledHistogram *prometheus.HistogramVec
}

// NewClientCollector for prometheus.
//
//nolint:dupl
func NewClientCollector(lc fx.Lifecycle, version version.Version) *ClientCollector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ClientCollector{
		started: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_started_total",
				Help:        "Total number of RPCs started on the client.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		handled: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_handled_total",
				Help:        "Total number of RPCs completed by the client, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method", "http_code"}),
		received: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_msg_received_total",
				Help:        "Total number of RPC stream messages received by the client.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		sent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_client_msg_sent_total",
				Help:        "Total number of HTTP stream messages sent by the client.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		handledHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
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
func (m *ClientCollector) Describe(ch chan<- *prometheus.Desc) {
	m.started.Describe(ch)
	m.handled.Describe(ch)
	m.received.Describe(ch)
	m.sent.Describe(ch)
	m.handledHistogram.Describe(ch)
}

// Collect is called by the prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ClientCollector) Collect(ch chan<- prometheus.Metric) {
	m.started.Collect(ch)
	m.handled.Collect(ch)
	m.received.Collect(ch)
	m.sent.Collect(ch)
	m.handledHistogram.Collect(ch)
}

// RoundTripper for prometheus.
func (m *ClientCollector) RoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &roundTripper{Metrics: m, RoundTripper: rt}
}

type roundTripper struct {
	Metrics *ClientCollector
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
	metrics     *ClientCollector
	serviceName string
	methodName  string
	startTime   time.Time
}

func newClientReporter(m *ClientCollector, service, method string) *clientReporter {
	r := &clientReporter{metrics: m, startTime: time.Now(), serviceName: service, methodName: method}
	r.metrics.started.WithLabelValues(r.serviceName, r.methodName).Inc()

	return r
}

func (r *clientReporter) ReceivedMessage() {
	r.metrics.received.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *clientReporter) SentMessage() {
	r.metrics.sent.WithLabelValues(r.serviceName, r.methodName).Inc()
}

func (r *clientReporter) Handled(code int) {
	r.metrics.handled.WithLabelValues(r.serviceName, r.methodName, strconv.Itoa(code)).Inc()
	r.metrics.handledHistogram.WithLabelValues(r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
}
