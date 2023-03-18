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
