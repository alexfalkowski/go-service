package prometheus

import (
	"time"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

// ProducerCollector for prometheus.
type ProducerCollector struct {
	started          *prometheus.CounterVec
	handled          *prometheus.CounterVec
	sent             *prometheus.CounterVec
	handledHistogram *prometheus.HistogramVec
}

// NewProducerCollector for prometheus.
func NewProducerCollector(lc fx.Lifecycle, version version.Version) *ProducerCollector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ProducerCollector{
		started: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_producer_started_total",
				Help:        "Total number of messages started by the producer.",
				ConstLabels: labels,
			}, []string{"nsq_topic"}),
		handled: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_producer_handled_total",
				Help:        "Total number of messages published, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"nsq_topic"}),
		sent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_producer_msg_sent_total",
				Help:        "Total number of stream messages sent by the producer.",
				ConstLabels: labels,
			}, []string{"nsq_topic"}),
		handledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "nsq_producer_handling_seconds",
				Help:        "Histogram of response latency (seconds) of messages that had been application-level handled by the producer.",
				Buckets:     prometheus.DefBuckets,
				ConstLabels: labels,
			}, []string{"nsq_topic"}),
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
func (m *ProducerCollector) Describe(ch chan<- *prometheus.Desc) {
	m.started.Describe(ch)
	m.handled.Describe(ch)
	m.sent.Describe(ch)
	m.handledHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ProducerCollector) Collect(ch chan<- prometheus.Metric) {
	m.started.Collect(ch)
	m.handled.Collect(ch)
	m.sent.Collect(ch)
	m.handledHistogram.Collect(ch)
}

// Producer for prometheus.
func (m *ProducerCollector) Producer(p producer.Producer) producer.Producer {
	return &promProducer{metrics: m, Producer: p}
}

type promProducer struct {
	metrics *ProducerCollector

	producer.Producer
}

func (p *promProducer) Publish(ctx context.Context, topic string, msg *message.Message) error {
	monitor := newProducerReporter(p.metrics, topic)
	monitor.SentMessage()

	err := p.Producer.Publish(ctx, topic, msg)
	if err != nil {
		return err
	}

	monitor.Handled()

	return nil
}

type producerReporter struct {
	metrics   *ProducerCollector
	topicName string
	startTime time.Time
}

func newProducerReporter(m *ProducerCollector, topic string) *producerReporter {
	r := &producerReporter{metrics: m, startTime: time.Now(), topicName: topic}
	r.metrics.started.WithLabelValues(r.topicName).Inc()

	return r
}

func (r *producerReporter) SentMessage() {
	r.metrics.sent.WithLabelValues(r.topicName).Inc()
}

func (r *producerReporter) Handled() {
	r.metrics.handled.WithLabelValues(r.topicName).Inc()
	r.metrics.handledHistogram.WithLabelValues(r.topicName).Observe(time.Since(r.startTime).Seconds())
}
