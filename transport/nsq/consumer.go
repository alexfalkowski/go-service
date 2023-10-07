package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/telemetry/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/telemetry/metrics"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	"github.com/nsqio/go-nsq"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerOption for NSQ.
type ConsumerOption interface{ apply(*consumerOptions) }

type consumerOptions struct {
	logger *zap.Logger
	tracer ntracer.Tracer
	meter  metric.Meter
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
func WithConsumerTracer(tracer ntracer.Tracer) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.tracer = tracer
	})
}

// WithConsumerMetrics for NSQ.
func WithConsumerMetrics(meter metric.Meter) ConsumerOption {
	return consumerOptionFunc(func(o *consumerOptions) {
		o.meter = meter
	})
}

// RegisterConsumer for NSQ.
func RegisterConsumer(lc fx.Lifecycle, topic, channel string, cfg *Config, h handler.Handler, m marshaller.Marshaller, opts ...ConsumerOption) error {
	defaultOptions := &consumerOptions{tracer: tracer.NewNoopTracer("nsq")}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	c, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return err
	}

	c.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	if defaultOptions.logger != nil {
		h = lzap.NewHandler(topic, channel, defaultOptions.logger, h)
	}

	if defaultOptions.meter != nil {
		handler, err := metrics.NewHandler(topic, channel, defaultOptions.meter, h)
		if err != nil {
			return err
		}

		h = handler
	}

	h = ntracer.NewHandler(topic, channel, defaultOptions.tracer, h)
	h = meta.NewHandler(h)

	c.AddHandler(handler.New(h, m))

	err = c.ConnectToNSQLookupd(cfg.LookupHost)
	if err != nil {
		return err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			c.Stop()

			return nil
		},
	})

	return nil
}
