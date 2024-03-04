package nsq

import (
	"context"

	gn "github.com/alexfalkowski/go-service/nsq"
	r "github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/nsq/breaker"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/retry"
	"github.com/alexfalkowski/go-service/transport/nsq/telemetry/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/telemetry/metrics"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	"github.com/nsqio/go-nsq"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProducerOption for NSQ.
type ProducerOption interface{ apply(opts *producerOptions) }

type producerOptions struct {
	logger    *zap.Logger
	tracer    ntracer.Tracer
	meter     metric.Meter
	retry     *r.Config
	userAgent string
	breaker   bool
}

type producerOptionFunc func(*producerOptions)

func (f producerOptionFunc) apply(o *producerOptions) { f(o) }

// WithProducerRetry for NSQ.
func WithProducerRetry(cfg *r.Config) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.retry = cfg
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
func WithProducerMetrics(meter metric.Meter) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.meter = meter
	})
}

// WithUserAgent for NSQ.
func WithProducerUserAgent(userAgent string) ProducerOption {
	return producerOptionFunc(func(o *producerOptions) {
		o.userAgent = userAgent
	})
}

// NewProducer for NSQ.
func NewProducer(lc fx.Lifecycle, host string, m gn.Marshaller, opts ...ProducerOption) (gn.Producer, error) {
	os := &producerOptions{tracer: tracer.NewNoopTracer("nsq")}
	for _, o := range opts {
		o.apply(os)
	}

	p, _ := nsq.NewProducer(host, nsq.NewConfig())
	p.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			p.Stop()

			return nil
		},
	})

	var pr gn.Producer = &nsqProducer{marshaller: m, Producer: p}

	if os.retry != nil {
		pr = retry.NewProducer(os.retry, pr)
	}

	if os.breaker {
		pr = breaker.NewProducer(pr)
	}

	if os.logger != nil {
		pr = lzap.NewProducer(os.logger, pr)
	}

	if os.meter != nil {
		producer, err := metrics.NewProducer(os.meter, pr)
		if err != nil {
			return nil, err
		}

		pr = producer
	}

	pr = ntracer.NewProducer(os.tracer, pr)
	pr = meta.NewProducer(os.userAgent, pr)

	return pr, nil
}

type nsqProducer struct {
	marshaller gn.Marshaller
	*nsq.Producer
}

// Produce a message to a topic.
func (p *nsqProducer) Produce(_ context.Context, topic string, msg *gn.Message) error {
	bytes, err := p.marshaller.Marshal(msg)
	if err != nil {
		return err
	}

	return p.Producer.Publish(topic, bytes)
}
