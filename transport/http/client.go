package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/http/breaker"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	lzap "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface{ apply(*clientOptions) }

type clientOptions struct {
	logger       *zap.Logger
	tracer       htracer.Tracer
	meter        metric.Meter
	retry        bool
	breaker      bool
	roundTripper http.RoundTripper
}

type clientOptionFunc func(*clientOptions)

func (f clientOptionFunc) apply(o *clientOptions) { f(o) }

// WithClientRoundTripper for HTTP.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.roundTripper = rt
	})
}

// WithClientRetry for HTTP.
func WithClientRetry() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.retry = true
	})
}

// WithClientBreaker for HTTP.
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.breaker = true
	})
}

// WithClientLogger for HTTP.
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.logger = logger
	})
}

// WithClientTracer for HTTP.
func WithClientTracer(tracer htracer.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// WithClientMetrics for HTTP.
func WithClientMetrics(meter metric.Meter) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.meter = meter
	})
}

// NewClient for HTTP.
func NewClient(cfg *Config, opts ...ClientOption) (*http.Client, error) {
	defaultOptions := &clientOptions{tracer: tracer.NewNoopTracer("http")}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	rt, err := newRoundTripper(cfg, defaultOptions)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: rt}

	return client, nil
}

func newRoundTripper(cfg *Config, opts *clientOptions) (http.RoundTripper, error) {
	hrt := opts.roundTripper
	if hrt == nil {
		hrt = http.DefaultTransport
	}

	if opts.logger != nil {
		hrt = lzap.NewRoundTripper(opts.logger, hrt)
	}

	if opts.meter != nil {
		rt, err := metrics.NewRoundTripper(opts.meter, hrt)
		if err != nil {
			return nil, err
		}

		hrt = rt
	}

	hrt = htracer.NewRoundTripper(opts.tracer, hrt)

	if opts.retry {
		hrt = retry.NewRoundTripper(&cfg.Retry, hrt)
	}

	if opts.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	hrt = meta.NewRoundTripper(cfg.UserAgent, hrt)

	return hrt, nil
}
