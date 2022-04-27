package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/breaker"
	"github.com/alexfalkowski/go-service/transport/nsq/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/transport/nsq/retry"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProducerOption for NSQ.
type ProducerOption interface{ apply(*producerOptions) }

type producerOptions struct {
	lifecycle  fx.Lifecycle
	config     *Config
	logger     *zap.Logger
	tracer     opentracing.Tracer
	retry      bool
	breaker    bool
	marshaller marshaller.Marshaller
	version    version.Version
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

// WithProducerConfig for NSQ.
func WithProducerConfig(config *Config) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.config = config
	})
}

// WithProducerLogger for NSQ.
func WithProducerLogger(logger *zap.Logger) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.logger = logger
	})
}

// WithProducerConfig for NSQ.
func WithProducerTracer(tracer opentracing.Tracer) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.tracer = tracer
	})
}

// WithProducerLifecycle for NSQ.
func WithProducerLifecycle(lifecycle fx.Lifecycle) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.lifecycle = lifecycle
	})
}

// WithProducerMarshaller for NSQ.
func WithProducerMarshaller(marshaller marshaller.Marshaller) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.marshaller = marshaller
	})
}

// WithCProducerVersion for NSQ.
func WithProducerVersion(version version.Version) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.version = version
	})
}

// NewProducer for NSQ.
func NewProducer(opts ...ProducerOption) producer.Producer {
	defaultOptions := &producerOptions{}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	cfg := nsq.NewConfig()
	p, _ := nsq.NewProducer(defaultOptions.config.Host, cfg)

	p.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	defaultOptions.lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			p.Stop()

			return nil
		},
	})

	var pr producer.Producer = &nsqProducer{marshaller: defaultOptions.marshaller, Producer: p}
	pr = lzap.NewProducer(defaultOptions.logger, pr)
	pr = opentracing.NewProducer(defaultOptions.tracer, pr)

	if defaultOptions.retry {
		pr = retry.NewProducer(&defaultOptions.config.Retry, pr)
	}

	if defaultOptions.breaker {
		pr = breaker.NewProducer(pr)
	}

	pr = meta.NewProducer(defaultOptions.config.UserAgent, defaultOptions.version, pr)

	return pr
}

type nsqProducer struct {
	marshaller marshaller.Marshaller
	*nsq.Producer
}

// Publish a message to a topic.
func (p *nsqProducer) Publish(ctx context.Context, topic string, msg *message.Message) error {
	bytes, err := p.marshaller.Marshal(msg)
	if err != nil {
		return err
	}

	return p.Producer.Publish(topic, bytes)
}
