package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	ntel "github.com/alexfalkowski/go-service/transport/nsq/telemetry"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerOption for NSQ.
type ConsumerOption interface{ apply(*consumerOptions) }

type consumerOptions struct {
	logger  *zap.Logger
	tracer  ntel.Tracer
	metrics *ntel.ConsumerMetrics
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
func WithConsumerTracer(tracer ntel.Tracer) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.tracer = tracer
	})
}

// WithConsumerMetrics for NSQ.
func WithConsumerMetrics(metrics *ntel.ConsumerMetrics) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.metrics = metrics
	})
}

// ConsumerParams for NSQ.
type ConsumerParams struct {
	Lifecycle fx.Lifecycle

	Topic, Channel string
	Config         *Config
	Handler        handler.Handler
	Marshaller     marshaller.Marshaller
}

// RegisterConsumer for NSQ.
func RegisterConsumer(params ConsumerParams, opts ...ConsumerOption) error {
	defaultOptions := &consumerOptions{tracer: telemetry.NewNoopTracer("nsq")}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	cfg := nsq.NewConfig()

	c, err := nsq.NewConsumer(params.Topic, params.Channel, cfg)
	if err != nil {
		return err
	}

	c.SetLogger(&logger{}, nsq.LogLevelInfo)

	h := params.Handler

	if defaultOptions.logger != nil {
		h = ntel.NewLoggerHandler(ntel.LoggerHandlerParams{Topic: params.Topic, Channel: params.Channel, Logger: defaultOptions.logger, Handler: h})
	}

	if defaultOptions.metrics != nil {
		h = defaultOptions.metrics.Handler(params.Topic, params.Channel, h)
	}

	h = ntel.NewTracerHandler(params.Topic, params.Channel, defaultOptions.tracer, h)
	h = meta.NewHandler(h)

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
