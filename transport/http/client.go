package http

import (
	"net/http"

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

// NewClient for HTTP.
func NewClient(config *Config, logger *zap.Logger, opts ...ClientOption) *http.Client {
	defaultOptions := &clientOptions{config: config, logger: logger}

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
	hrt = opentracing.NewRoundTripper(hrt)

	if opts.retry {
		hrt = retry.NewRoundTripper(&opts.config.Retry, hrt)
	}

	if opts.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	hrt = meta.NewRoundTripper(opts.config.UserAgent, hrt)

	return hrt
}
