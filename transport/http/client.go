package http

import (
	"net/http"

	sopentracing "github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/http/breaker"
	lzap "github.com/alexfalkowski/go-service/transport/http/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface{ apply(*clientOptions) }

type clientOptions struct {
	config       *Config
	logger       *zap.Logger
	tracer       sopentracing.TransportTracer
	retry        bool
	breaker      bool
	roundTripper http.RoundTripper
}

type clientOptionFunc func(*clientOptions)

func (f clientOptionFunc) apply(o *clientOptions) { f(o) }

// WithClientRoundTripper for HTTP.
// nolint:ireturn
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.roundTripper = rt
	})
}

// WithClientRetry for HTTP.
// nolint:ireturn
func WithClientRetry() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.retry = true
	})
}

// WithClientBreaker for HTTP.
// nolint:ireturn
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.breaker = true
	})
}

// WithClientConfig for gRPC.
// nolint:ireturn
func WithClientConfig(config *Config) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.config = config
	})
}

// WithClientLogger for gRPC.
// nolint:ireturn
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.logger = logger
	})
}

// WithClientConfig for gRPC.
// nolint:ireturn
func WithClientTracer(tracer sopentracing.TransportTracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// NewClient for HTTP.
func NewClient(opts ...ClientOption) *http.Client {
	defaultOptions := &clientOptions{}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	return &http.Client{Transport: newRoundTripper(defaultOptions)}
}

func newRoundTripper(opts *clientOptions) http.RoundTripper {
	hrt := opts.roundTripper
	if hrt == nil {
		hrt = http.DefaultTransport
	}

	hrt = lzap.NewRoundTripper(opts.logger, hrt)
	hrt = opentracing.NewRoundTripper(opts.tracer, hrt)

	if opts.retry {
		hrt = retry.NewRoundTripper(&opts.config.Retry, hrt)
	}

	if opts.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	hrt = meta.NewRoundTripper(opts.config.UserAgent, hrt)

	return hrt
}
