package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/nsq/breaker"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/transport/nsq/retry"
	"github.com/alexfalkowski/go-service/transport/nsq/telemetry/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/telemetry/metrics/prometheus"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProducerOption for NSQ.
type ProducerOption interface{ apply(*producerOptions) }

type producerOptions struct {
	logger  *zap.Logger
	tracer  ntracer.Tracer
	metrics *prometheus.ProducerCollector
	retry   bool
	breaker bool
}

type producerOptionFunc func(*producerOptions)

func (f producerOptionFunc) apply(o *producerOptions) { f(o) }

// WithProducerRetry for NSQ.
func WithProducerRetry() ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.retry = true
	})
}

// WithProducerBreaker for NSQ.
func WithProducerBreaker() ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.breaker = true
	})
}

// WithProducerLogger for NSQ.
func WithProducerLogger(logger *zap.Logger) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.logger = logger
	})
}

// WithProducerConfig for NSQ.
func WithProducerTracer(tracer ntracer.Tracer) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.tracer = tracer
	})
}

// WithProducerMetrics for NSQ.
func WithProducerMetrics(metrics *prometheus.ProducerCollector) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.metrics = metrics
	})
}

// NewProducer for NSQ.
func NewProducer(lc fx.Lifecycle, cfg *Config, m marshaller.Marshaller, opts ...ProducerOption) producer.Producer {
	defaultOptions := &producerOptions{tracer: tracer.NewNoopTracer("nsq")}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	p, _ := nsq.NewProducer(cfg.Host, nsq.NewConfig())
	p.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			p.Stop()

			return nil
		},
	})

	var pr producer.Producer = &nsqProducer{marshaller: m, Producer: p}

	if defaultOptions.logger != nil {
		pr = lzap.NewProducer(defaultOptions.logger, pr)
	}

	if defaultOptions.metrics != nil {
		pr = defaultOptions.metrics.Producer(pr)
	}

	pr = ntracer.NewProducer(defaultOptions.tracer, pr)

	if defaultOptions.retry {
		pr = retry.NewProducer(&cfg.Retry, pr)
	}

	if defaultOptions.breaker {
		pr = breaker.NewProducer(pr)
	}

	pr = meta.NewProducer(cfg.UserAgent, pr)

	return pr
}

type nsqProducer struct {
	marshaller marshaller.Marshaller
	*nsq.Producer
}

// Publish a message to a topic.
func (p *nsqProducer) Publish(_ context.Context, topic string, msg *message.Message) error {
	bytes, err := p.marshaller.Marshal(msg)
	if err != nil {
		return err
	}

	return p.Producer.Publish(topic, bytes)
}
