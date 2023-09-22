package prometheus

import (
	"time"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

// ConsumerCollector for prometheus.
type ConsumerCollector struct {
	started          *prometheus.CounterVec
	handled          *prometheus.CounterVec
	received         *prometheus.CounterVec
	handledHistogram *prometheus.HistogramVec
}

// NewConsumerCollector for prometheus.
func NewConsumerCollector(lc fx.Lifecycle, version version.Version) *ConsumerCollector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ConsumerCollector{
		started: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_consumer_started_total",
				Help:        "Total number of messages started to be consumed.",
				ConstLabels: labels,
			}, []string{"nsq_topic", "nsq_channel"}),
		handled: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_consumer_handled_total",
				Help:        "Total number of messages consumed, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"nsq_topic", "nsq_channel"}),
		received: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_consumer_msg_received_total",
				Help:        "Total number of messages consumned.",
				ConstLabels: labels,
			}, []string{"nsq_topic", "nsq_channel"}),
		handledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "nsq_consumer_handling_seconds",
				Help:        "Histogram of response latency (seconds) of messages that had been consumed.",
				Buckets:     prometheus.DefBuckets,
				ConstLabels: labels,
			}, []string{"nsq_topic", "nsq_channel"},
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
func (m *ConsumerCollector) Describe(ch chan<- *prometheus.Desc) {
	m.started.Describe(ch)
	m.handled.Describe(ch)
	m.received.Describe(ch)
	m.handledHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ConsumerCollector) Collect(ch chan<- prometheus.Metric) {
	m.started.Collect(ch)
	m.handled.Collect(ch)
	m.received.Collect(ch)
	m.handledHistogram.Collect(ch)
}

// ServerHandler for prometheus.
func (m *ConsumerCollector) Handler(topic, channel string, h handler.Handler) handler.Handler {
	return &promConsumer{topic: topic, channel: channel, metrics: m, Handler: h}
}

type promConsumer struct {
	topic, channel string
	metrics        *ConsumerCollector
	handler.Handler
}

func (h *promConsumer) Handle(ctx context.Context, message *message.Message) error {
	monitor := newConsumerReporter(h.metrics, h.topic, h.channel)
	monitor.ReceivedMessage()

	if err := h.Handler.Handle(ctx, message); err != nil {
		return err
	}

	monitor.Handled()

	return nil
}

type consumerReporter struct {
	metrics     *ConsumerCollector
	topicName   string
	channelName string
	startTime   time.Time
}

func newConsumerReporter(m *ConsumerCollector, topic, channel string) *consumerReporter {
	r := &consumerReporter{metrics: m, startTime: time.Now(), topicName: topic, channelName: channel}
	r.metrics.started.WithLabelValues(r.topicName, r.channelName).Inc()

	return r
}

func (r *consumerReporter) ReceivedMessage() {
	r.metrics.received.WithLabelValues(r.topicName, r.channelName).Inc()
}

func (r *consumerReporter) Handled() {
	r.metrics.handled.WithLabelValues(r.topicName, r.channelName).Inc()
	r.metrics.handledHistogram.WithLabelValues(r.topicName, r.channelName).Observe(time.Since(r.startTime).Seconds())
}
