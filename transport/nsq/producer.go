package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/transport/nsq/breaker"
	"github.com/alexfalkowski/go-service/transport/nsq/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/metrics/prometheus"
	notel "github.com/alexfalkowski/go-service/transport/nsq/otel"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/transport/nsq/retry"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProducerOption for NSQ.
type ProducerOption interface{ apply(*producerOptions) }

type producerOptions struct {
	logger  *zap.Logger
	tracer  notel.Tracer
	metrics *prometheus.ProducerMetrics
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
func WithProducerTracer(tracer notel.Tracer) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.tracer = tracer
	})
}

// WithProducerMetrics for NSQ.
func WithProducerMetrics(metrics *prometheus.ProducerMetrics) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.metrics = metrics
	})
}

// ProducerParams for NSQ.
type ProducerParams struct {
	Lifecycle fx.Lifecycle

	Config     *Config
	Marshaller marshaller.Marshaller
}

// NewProducer for NSQ.
func NewProducer(params ProducerParams, opts ...ProducerOption) producer.Producer {
	defaultOptions := &producerOptions{tracer: otel.NewNoopTracer("nsq")}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	cfg := nsq.NewConfig()
	p, _ := nsq.NewProducer(params.Config.Host, cfg)

	p.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			p.Stop()

			return nil
		},
	})

	var pr producer.Producer = &nsqProducer{marshaller: params.Marshaller, Producer: p}

	if defaultOptions.logger != nil {
		pr = lzap.NewProducer(lzap.ProducerParams{Logger: defaultOptions.logger, Producer: pr})
	}

	if defaultOptions.metrics != nil {
		pr = defaultOptions.metrics.Producer(pr)
	}

	pr = notel.NewProducer(defaultOptions.tracer, pr)

	if defaultOptions.retry {
		pr = retry.NewProducer(&params.Config.Retry, pr)
	}

	if defaultOptions.breaker {
		pr = breaker.NewProducer(pr)
	}

	pr = meta.NewProducer(params.Config.UserAgent, pr)

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
