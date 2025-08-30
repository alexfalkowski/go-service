package http

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/meta"
	"github.com/alexfalkowski/go-service/v2/transport/http/retry"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
)

// ClientOption for HTTP.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	gen          token.Generator
	generator    id.Generator
	roundTripper http.RoundTripper
	tracer       *tracer.Tracer
	logger       *logger.Logger
	retry        *retry.Config
	tls          *tls.Config
	meter        *metrics.Meter
	limiter      *limiter.Client
	userAgent    env.UserAgent
	id           env.UserID
	timeout      time.Duration
	breaker      bool
	compression  bool
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithClientCompression for HTTP.
func WithClientCompression() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.compression = true
	})
}

// WithClientTokenGenerator for HTTP.
func WithClientTokenGenerator(id env.UserID, gen token.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.id = id
		o.gen = gen
	})
}

// WithClientTimeout for HTTP.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = time.MustParseDuration(timeout)
	})
}

// WithClientRoundTripper for HTTP.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientRetry for HTTP.
func WithClientRetry(cfg *retry.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.retry = cfg
	})
}

// WithClientBreaker for HTTP.
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.breaker = true
	})
}

// WithClientLogger for HTTP.
func WithClientLogger(logger *logger.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// WithClientTracer for HTTP.
func WithClientTracer(tracer *tracer.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tracer = tracer
	})
}

// WithClientMetrics for HTTP.
func WithClientMetrics(meter *metrics.Meter) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.meter = meter
	})
}

// WithClientUserAgent for HTTP.
func WithClientUserAgent(userAgent env.UserAgent) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// WithClientTLS for HTTP.
func WithClientTLS(sec *tls.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tls = sec
	})
}

// WithClientID for HTTP.
func WithClientID(generator id.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.generator = generator
	})
}

// WithClientLimiter for HTTP.
func WithClientLimiter(limiter *limiter.Client) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.limiter = limiter
	})
}

// NewRoundTripper for HTTP.
func NewRoundTripper(opts ...ClientOption) (http.RoundTripper, error) {
	os := options(opts...)

	hrt, err := roundTripper(os)
	if err != nil {
		return nil, err
	}

	if os.gen != nil {
		hrt = token.NewRoundTripper(os.id, os.gen, hrt)
	}

	if os.compression {
		hrt = gzhttp.Transport(hrt, gzhttp.TransportEnableGzip(true))
	}

	if os.limiter != nil {
		hrt = limiter.NewRoundTripper(os.limiter, hrt)
	}

	if os.retry != nil {
		hrt = retry.NewRoundTripper(os.retry, hrt)
	}

	if os.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	if os.logger != nil {
		hrt = logger.NewRoundTripper(os.logger, hrt)
	}

	if os.meter != nil {
		hrt = metrics.NewRoundTripper(os.meter, hrt)
	}

	if os.tracer != nil {
		hrt = tracer.NewRoundTripper(os.tracer, hrt)
	}

	hrt = meta.NewRoundTripper(os.userAgent, os.generator, hrt)

	return hrt, nil
}

// NewClient for HTTP.
func NewClient(opts ...ClientOption) (*http.Client, error) {
	os := options(opts...)

	transport, err := NewRoundTripper(opts...)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   os.timeout,
	}

	return client, nil
}

func roundTripper(os *clientOpts) (http.RoundTripper, error) {
	hrt := os.roundTripper
	if hrt != nil {
		return hrt, nil
	}

	if !os.tls.IsEnabled() {
		return http.Transport(nil), nil
	}

	conf, err := tls.NewConfig(fs, os.tls)
	if err != nil {
		return nil, err
	}

	return http.Transport(conf), nil
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	if os.timeout == 0 {
		os.timeout = 30 * time.Second
	}

	if os.generator == nil {
		os.generator = &id.UUID{}
	}

	return os
}
