package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/transport/http/breaker"
	lzap "github.com/alexfalkowski/go-service/transport/http/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface{ apply(*clientOptions) }

type clientOptions struct {
	logger       *zap.Logger
	tracer       opentracing.Tracer
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
func WithClientTracer(tracer opentracing.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// ClientParams for HTTP.
type ClientParams struct {
	Version version.Version
	Config  *Config
}

// NewClient for HTTP.
func NewClient(params ClientParams, opts ...ClientOption) *http.Client {
	defaultOptions := &clientOptions{}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	return &http.Client{Transport: newRoundTripper(params, defaultOptions)}
}

func newRoundTripper(params ClientParams, opts *clientOptions) http.RoundTripper {
	hrt := opts.roundTripper
	if hrt == nil {
		hrt = http.DefaultTransport
	}

	if opts.logger != nil {
		hrt = lzap.NewRoundTripper(opts.logger, hrt)
	}

	if opts.tracer != nil {
		hrt = opentracing.NewRoundTripper(opts.tracer, hrt)
	}

	if opts.retry {
		hrt = retry.NewRoundTripper(&params.Config.Retry, hrt)
	}

	if opts.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	hrt = meta.NewRoundTripper(params.Config.UserAgent, params.Version, hrt)

	return hrt
}
