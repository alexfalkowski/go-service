package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerOption for NSQ.
type ConsumerOption interface{ apply(*consumerOptions) }

type consumerOptions struct {
	logger *zap.Logger
	tracer opentracing.Tracer
}

type consumerOptionFunc func(*consumerOptions)

func (f consumerOptionFunc) apply(o *consumerOptions) { f(o) }

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

// ConsumerParams for NSQ.
type ConsumerParams struct {
	Lifecycle fx.Lifecycle

	Topic, Channel string
	Config         *Config
	Version        version.Version
	Handler        handler.Handler
	Marshaller     marshaller.Marshaller
}

// RegisterConsumer for NSQ.
func RegisterConsumer(params ConsumerParams, opts ...ConsumerOption) error {
	defaultOptions := &consumerOptions{}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	cfg := nsq.NewConfig()

	c, err := nsq.NewConsumer(params.Topic, params.Channel, cfg)
	if err != nil {
		return err
	}

	c.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	var h handler.Handler = params.Handler

	if defaultOptions.logger != nil {
		h = lzap.NewHandler(params.Topic, params.Channel, defaultOptions.logger, h)
	}

	if defaultOptions.tracer != nil {
		h = opentracing.NewHandler(params.Topic, params.Channel, defaultOptions.tracer, h)
	}

	h = meta.NewHandler(params.Version, h)

	c.AddHandler(handler.New(handler.Params{Handler: h, Marshaller: params.Marshaller}))

	err = c.ConnectToNSQLookupd(params.Config.LookupHost)
	if err != nil {
		return err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			c.Stop()

			return nil
		},
	})

	return nil
}
