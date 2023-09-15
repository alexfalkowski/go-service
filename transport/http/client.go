package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/transport/http/breaker"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	htel "github.com/alexfalkowski/go-service/transport/http/telemetry"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface{ apply(*clientOptions) }

type clientOptions struct {
	logger       *zap.Logger
	tracer       htel.Tracer
	metrics      *htel.ClientMetrics
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
func WithClientTracer(tracer htel.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// WithClientMetrics for HTTP.
func WithClientMetrics(metrics *htel.ClientMetrics) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.metrics = metrics
	})
}

// ClientParams for HTTP.
type ClientParams struct {
	Config *Config
}

// NewClient for HTTP.
func NewClient(params ClientParams, opts ...ClientOption) *http.Client {
	defaultOptions := &clientOptions{tracer: telemetry.NewNoopTracer("http")}
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
		hrt = htel.NewLoggerRoundTripper(htel.LoggerRoundTripperParams{Logger: opts.logger, RoundTripper: hrt})
	}

	if opts.metrics != nil {
		hrt = opts.metrics.RoundTripper(hrt)
	}

	hrt = htel.NewTracerRoundTripper(opts.tracer, hrt)

	if opts.retry {
		hrt = retry.NewRoundTripper(&params.Config.Retry, hrt)
	}

	if opts.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	hrt = meta.NewRoundTripper(params.Config.UserAgent, hrt)

	return hrt
}
