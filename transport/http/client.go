package http

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/meta"
	"github.com/alexfalkowski/go-service/v2/transport/http/retry"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
)

// ClientOption configures HTTP client construction.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	gen            token.Generator
	generator      id.Generator
	roundTripper   http.RoundTripper
	limiter        *limiter.Client
	retry          *retry.Config
	tls            *tls.Config
	logger         *logger.Logger
	userAgent      env.UserAgent
	id             env.UserID
	breakerOptions []breaker.Option
	timeout        time.Duration
	breaker        bool
	compression    bool
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithClientCompression enables gzip compression for HTTP client requests.
func WithClientCompression() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.compression = true
	})
}

// WithClientTokenGenerator enables token injection using gen and id.
func WithClientTokenGenerator(id env.UserID, gen token.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.id = id
		o.gen = gen
	})
}

// WithClientTimeout sets the http.Client timeout used by NewClient.
//
// If unset, a default timeout is applied (see options()).
func WithClientTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = timeout
	})
}

// WithClientRoundTripper sets the base HTTP RoundTripper to wrap.
//
// If set, this round tripper is used as-is (TLS config and default transport selection are skipped).
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientRetry enables retry behavior for outbound HTTP requests.
func WithClientRetry(cfg *retry.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.retry = cfg
	})
}

// WithClientBreaker enables circuit breaking for outbound HTTP requests.
func WithClientBreaker(options ...breaker.Option) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.breaker = true
		o.breakerOptions = options
	})
}

// WithClientLogger enables HTTP client logging middleware.
func WithClientLogger(logger *logger.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// WithClientUserAgent sets the user agent value used for metadata injection.
func WithClientUserAgent(userAgent env.UserAgent) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// WithClientTLS enables TLS for the default HTTP transport selection.
//
// If TLS is enabled and no base round tripper is provided, TLS config is constructed using the registered
// filesystem dependency (see Register in this package).
func WithClientTLS(sec *tls.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tls = sec
	})
}

// WithClientID sets the request id generator used by metadata injection middleware.
func WithClientID(generator id.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.generator = generator
	})
}

// WithClientLimiter enables client-side rate limiting middleware.
func WithClientLimiter(limiter *limiter.Client) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.limiter = limiter
	})
}

// NewRoundTripper constructs an HTTP RoundTripper by composing middleware derived from opts.
//
// Defaults (see options()):
//   - timeout: 30s
//   - request-id generator: uuid
//
// If no base round tripper is configured (via WithClientRoundTripper), a transport is selected based on TLS:
//   - if TLS is disabled: the default HTTP transport is used
//   - if TLS is enabled: TLS config is constructed using the registered filesystem (see Register)
//
// Middleware is applied in the following order (outermost first):
//   - meta injection (User-Agent, Request-Id)
//   - logger (optional)
//   - breaker (optional)
//   - retry (optional)
//   - limiter (optional)
//   - compression (optional)
//   - token injection (optional)
//   - base transport
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
		hrt = breaker.NewRoundTripper(hrt, os.breakerOptions...)
	}

	if os.logger != nil {
		hrt = logger.NewRoundTripper(os.logger, hrt)
	}

	hrt = meta.NewRoundTripper(os.userAgent, os.generator, hrt)

	return hrt, nil
}

// NewClient constructs a new instrumented http.Client.
//
// The returned client uses the RoundTripper built by NewRoundTripper and applies the configured timeout.
// Note: http.NewClient wraps the transport with OpenTelemetry instrumentation.
func NewClient(opts ...ClientOption) (*http.Client, error) {
	os := options(opts...)

	transport, err := NewRoundTripper(opts...)
	if err != nil {
		return nil, err
	}

	return http.NewClient(transport, os.timeout), nil
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
		os.generator = uuid.NewGenerator()
	}

	return os
}
