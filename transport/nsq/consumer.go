package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerOption for NSQ.
type ConsumerOption interface{ apply(*consumerOptions) }

type consumerOptions struct {
	lifecycle  fx.Lifecycle
	config     *Config
	logger     *zap.Logger
	tracer     opentracing.Tracer
	marshaller marshaller.Marshaller
	handler    handler.Handler
}

type consumerOptionFunc func(*consumerOptions)

func (f consumerOptionFunc) apply(o *consumerOptions) { f(o) }

// WithConsumerConfig for NSQ.
func WithConsumerConfig(config *Config) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.config = config
	})
}

// WithConsumerLogger for NSQ.
func WithConsumerLogger(logger *zap.Logger) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.logger = logger
	})
}

// WithConsumerConfig for NSQ.
func WithConsumerTracer(tracer opentracing.Tracer) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.tracer = tracer
	})
}

// WithConsumerLifecycle for NSQ.
func WithConsumerLifecycle(lifecycle fx.Lifecycle) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.lifecycle = lifecycle
	})
}

// WithConsumerMarshaller for NSQ.
func WithConsumerMarshaller(marshaller marshaller.Marshaller) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.marshaller = marshaller
	})
}

// WithConsumerMarshaller for NSQ.
func WithConsumerHandler(handler handler.Handler) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.handler = handler
	})
}

// RegisterConsumer for NSQ.
func RegisterConsumer(topic string, channel string, opts ...ConsumerOption) error {
	defaultOptions := &consumerOptions{}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	cfg := nsq.NewConfig()

	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return err
	}

	c.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	var h handler.Handler = lzap.NewHandler(topic, channel, defaultOptions.logger, defaultOptions.handler)
	h = opentracing.NewHandler(topic, channel, defaultOptions.tracer, h)
	h = meta.NewHandler(h)

	c.AddHandler(handler.New(handler.Params{Handler: h, Marshaller: defaultOptions.marshaller}))

	err = c.ConnectToNSQLookupd(defaultOptions.config.LookupHost)
	if err != nil {
		return err
	}

	defaultOptions.lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			c.Stop()

			return nil
		},
	})

	return nil
}
