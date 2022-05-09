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

// ProducerMetrics for prometheus.
type ProducerMetrics struct {
	producerStartedCounter   *prometheus.CounterVec
	producerHandledCounter   *prometheus.CounterVec
	producerMsgSent          *prometheus.CounterVec
	producerHandledHistogram *prometheus.HistogramVec
}

// NewProducerMetrics for prometheus.
func NewProducerMetrics(lc fx.Lifecycle, version version.Version) *ProducerMetrics {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ProducerMetrics{
		producerStartedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_producer_started_total",
				Help:        "Total number of messages started by the producer.",
				ConstLabels: labels,
			}, []string{"nsq_topic"}),
		producerHandledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_producer_handled_total",
				Help:        "Total number of messages published, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"nsq_topic"}),
		producerMsgSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "nsq_producer_msg_sent_total",
				Help:        "Total number of stream messages sent by the producer.",
				ConstLabels: labels,
			}, []string{"nsq_topic"}),
		producerHandledHistogram: prometheus.NewHistogramVec(
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
func (m *ProducerMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.producerStartedCounter.Describe(ch)
	m.producerHandledCounter.Describe(ch)
	m.producerMsgSent.Describe(ch)
	m.producerHandledHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ProducerMetrics) Collect(ch chan<- prometheus.Metric) {
	m.producerStartedCounter.Collect(ch)
	m.producerHandledCounter.Collect(ch)
	m.producerMsgSent.Collect(ch)
	m.producerHandledHistogram.Collect(ch)
}

// Producer for prometheus.
func (m *ProducerMetrics) Producer(p producer.Producer) producer.Producer {
	return &promProducer{metrics: m, Producer: p}
}

type promProducer struct {
	metrics *ProducerMetrics

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
	metrics   *ProducerMetrics
	topicName string
	startTime time.Time
}

func newProducerReporter(m *ProducerMetrics, topic string) *producerReporter {
	r := &producerReporter{metrics: m, startTime: time.Now(), topicName: topic}
	r.metrics.producerStartedCounter.WithLabelValues(r.topicName).Inc()

	return r
}

func (r *producerReporter) SentMessage() {
	r.metrics.producerMsgSent.WithLabelValues(r.topicName).Inc()
}

func (r *producerReporter) Handled() {
	r.metrics.producerHandledCounter.WithLabelValues(r.topicName).Inc()
	r.metrics.producerHandledHistogram.WithLabelValues(r.topicName).Observe(time.Since(r.startTime).Seconds())
}
