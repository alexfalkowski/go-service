package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/http/breaker"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	lzap "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/metrics/prometheus"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface{ apply(*clientOptions) }

type clientOptions struct {
	logger       *zap.Logger
	tracer       htracer.Tracer
	metrics      *prometheus.ClientCollector
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
func WithClientMetrics(metrics *prometheus.ClientCollector) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.metrics = metrics
	})
}

// NewClient for HTTP.
func NewClient(cfg *Config, opts ...ClientOption) *http.Client {
	defaultOptions := &clientOptions{tracer: tracer.NewNoopTracer("http")}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	return &http.Client{Transport: newRoundTripper(cfg, defaultOptions)}
}

func newRoundTripper(cfg *Config, opts *clientOptions) http.RoundTripper {
	hrt := opts.roundTripper
	if hrt == nil {
		hrt = http.DefaultTransport
	}

	if opts.logger != nil {
		hrt = lzap.NewRoundTripper(opts.logger, hrt)
	}

	if opts.metrics != nil {
		hrt = opts.metrics.RoundTripper(hrt)
	}

	hrt = htracer.NewRoundTripper(opts.tracer, hrt)

	if opts.retry {
		hrt = retry.NewRoundTripper(&cfg.Retry, hrt)
	}

	if opts.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	hrt = meta.NewRoundTripper(cfg.UserAgent, hrt)

	return hrt
}
