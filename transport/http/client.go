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

// NewRoundTripper for HTTP.
func NewRoundTripper(cfg *Config, opts ...ClientOption) (http.RoundTripper, error) {
	os := &clientOptions{tracer: tracer.NewNoopTracer("http")}
	for _, o := range opts {
		o.apply(os)
	}

	hrt := os.roundTripper
	if hrt == nil {
		hrt = transport()
	}

	if os.logger != nil {
		hrt = lzap.NewRoundTripper(os.logger, hrt)
	}

	if os.meter != nil {
		rt, err := metrics.NewRoundTripper(os.meter, hrt)
		if err != nil {
			return nil, err
		}

		hrt = rt
	}

	hrt = htracer.NewRoundTripper(os.tracer, hrt)

	if os.retry {
		hrt = retry.NewRoundTripper(&cfg.Retry, hrt)
	}

	if os.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	hrt = meta.NewRoundTripper(cfg.UserAgent, hrt)

	return hrt, nil
}

// NewClient for HTTP.
func NewClient(cfg *Config, opts ...ClientOption) (*http.Client, error) {
	defaultOptions := &clientOptions{tracer: tracer.NewNoopTracer("http")}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	rt, err := NewRoundTripper(cfg, opts...)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: rt}

	return client, nil
}

func transport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return t
}
