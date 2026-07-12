package http

import (
	"github.com/alexfalkowski/go-service/v2/config/client"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/retry"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
)

// ClientOption configures HTTP client construction.
//
// Client options are applied when assembling the client's [http.RoundTripper] stack (see [NewRoundTripper])
// and when determining the [http.Client] timeout (see [NewClient]).
//
// Most options are orthogonal and can be combined. Some options enable behavior by providing a non-nil
// dependency (for example, retries are enabled when [WithClientRetry] provides a non-nil config).
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	gen           token.Generator
	generator     id.Generator
	roundTripper  http.RoundTripper
	tls           *tls.Config
	retry         *retry.Config
	limiter       *limiter.Client
	logger        *logger.Logger
	breaker       *breaker.Config
	userAgent     env.UserAgent
	id            env.UserID
	retryPolicies []retry.Policy
	timeout       time.Duration
	compression   bool
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithClientCompression enables gzip compression for HTTP client requests.
//
// When enabled, the composed RoundTripper will advertise support for gzip and transparently decompress
// gzipped responses when the server sends them.
//
// This option uses [github.com/klauspost/compress/gzhttp] and wraps the underlying transport.
func WithClientCompression() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.compression = true
	})
}

// WithClientTokenGenerator enables token injection using gen and id.
//
// When configured, the client will generate an Authorization token per request and add it to the outbound
// request using the `Bearer` scheme. Token generation is scoped to the request
// method plus URL path (for example `DELETE /users/123`) and the configured user
// id. Query parameters are not included in the audience.
//
// Token-enabled clients use [http.SameOriginRedirect] so credentials are not forwarded to cross-origin
// redirect targets. Cross-origin redirects are returned to the caller as [http.ErrUseLastResponse].
func WithClientTokenGenerator(id env.UserID, gen token.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.id = id
		o.gen = gen
	})
}

// WithClientTimeout sets the [http.Client] timeout used by [NewClient].
//
// This is the overall time limit for requests made by the constructed client (including connection time,
// redirects, and reading the response body), as enforced by Go's [http.Client.Timeout] semantics.
//
// If unset or negative, a default timeout is applied (see `options()` defaults).
func WithClientTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = timeout
	})
}

// WithClientRoundTripper sets the base HTTP RoundTripper to wrap.
//
// This is an escape hatch for supplying a fully configured transport (for example, one with custom proxy,
// connection pooling, or test doubles).
//
// If set, this round tripper is used as-is and the package will not perform default transport selection
// based on TLS configuration (i.e. TLS config construction and `http.Transport(...)` selection are skipped).
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientRetry enables retry behavior for outbound HTTP requests.
//
// When configured, the composed RoundTripper will retry failed requests according to cfg (attempt count,
// backoff, and which status codes are considered retryable). Optional policies decide whether a logical
// request is eligible for retry.
//
// When no policy is provided, only side-effect-safe requests are eligible for retry: safe HTTP methods, or
// requests carrying a request-id idempotency contract. Callers that need different behavior can pass an
// explicit policy.
//
// If cfg is nil, retries are not enabled.
func WithClientRetry(cfg *retry.Config, policies ...retry.Policy) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.retry = cfg
		o.retryPolicies = policies
	})
}

// WithClientBreaker enables circuit breaking for outbound HTTP requests from cfg.
//
// Circuit breakers are applied in a RoundTripper wrapper. Breaker instances are keyed per request
// (by method + host) so that failures are isolated by downstream destination.
// Failure classification is controlled by cfg.
func WithClientBreaker(cfg *breaker.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.breaker = cfg
	})
}

// WithClientLogger enables HTTP client logging middleware.
//
// When configured, the composed RoundTripper logs request outcomes (duration and status classification).
func WithClientLogger(logger *logger.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// WithClientUserAgent sets the user agent value used for metadata injection.
//
// The value is injected into outbound requests by the metadata RoundTripper ([github.com/alexfalkowski/go-service/v2/net/http/meta]).
func WithClientUserAgent(userAgent env.UserAgent) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// WithClientTLS enables TLS for the default HTTP transport selection.
//
// This option only affects transport selection when [WithClientRoundTripper] is not used. When a base
// round tripper is provided explicitly, it is used as-is.
//
// If TLS is enabled and no base round tripper is provided, TLS config is constructed using the package-
// registered filesystem dependency (see [Register] in this package) to resolve TLS source strings.
func WithClientTLS(sec *tls.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tls = sec
	})
}

// WithClientID sets the request id generator used by metadata injection middleware.
//
// The generator is used to create a request id when one is not already present on the outgoing context
// (as propagated by the meta package) and/or request headers.
func WithClientID(generator id.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.generator = generator
	})
}

// WithClientLimiter enables client-side rate limiting middleware.
//
// When configured, outbound requests are rate-limited before being sent. If limiter is nil, rate limiting
// is not enabled.
func WithClientLimiter(limiter *limiter.Client) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.limiter = limiter
	})
}

// NewRoundTripper constructs an [http.RoundTripper] by composing middleware derived from opts.
//
// Defaults (see `options()`):
//   - timeout: 30s (used by [NewClient], not by the transport itself)
//   - request-id generator: uuid
//
// Base transport selection:
// If no base round tripper is configured (via [WithClientRoundTripper]), a transport is selected based on TLS:
//   - if TLS is disabled: the default HTTP transport is used (`http.Transport(nil)`)
//   - if TLS is enabled: TLS config is constructed using the registered filesystem (see [Register]) and used
//     to construct a TLS-enabled transport (`http.Transport(conf)`).
//
// Composition:
// Middleware is applied by wrapping the base RoundTripper. The order in which wrappers are applied matters.
// Given the implementation below, the resulting call stack (outermost → innermost) is:
//
//   - meta injection (User-Agent, Request-Id)
//   - logger (optional)
//   - retry (optional)
//   - limiter (optional)
//   - breaker (optional)
//   - token injection (optional)
//   - compression (optional)
//   - base transport
//
// Notes:
//   - [WithClientRoundTripper] short-circuits base transport selection and TLS config construction.
//   - Token injection remains inside retry so a fresh token is generated for each retry attempt.
//   - Retry wraps limiter and breaker so each retry attempt consumes quota and breaker capacity.
//   - The limiter stays outside the breaker so local quota denials are not counted as upstream failures.
func NewRoundTripper(opts ...ClientOption) (http.RoundTripper, error) {
	resolved := options(opts...)

	hrt, err := roundTripper(resolved)
	if err != nil {
		return nil, err
	}

	if resolved.compression {
		hrt = gzhttp.Transport(hrt, gzhttp.TransportEnableGzip(true))
	}

	if resolved.gen != nil {
		hrt = token.NewRoundTripper(resolved.id, resolved.gen, hrt)
	}

	if resolved.breaker != nil {
		hrt = breaker.NewRoundTripper(hrt, resolved.breaker.Options()...)
	}

	if resolved.limiter != nil {
		hrt = limiter.NewRoundTripper(resolved.limiter, hrt)
	}

	if resolved.retry != nil {
		hrt = retry.NewRoundTripper(resolved.retry, hrt, resolved.retryPolicies...)
	}

	if resolved.logger != nil {
		hrt = logger.NewRoundTripper(resolved.logger, hrt)
	}

	hrt = meta.NewRoundTripper(resolved.userAgent, resolved.generator, hrt)

	return hrt, nil
}

// NewClient constructs a new [http.Client].
//
// The returned client:
//   - uses the RoundTripper stack built by [NewRoundTripper], and
//   - applies the configured client timeout (see [WithClientTimeout]).
//
// If token generation is enabled, the returned client uses [http.SameOriginRedirect] so cross-origin
// redirects are returned to the caller instead of forwarding credentials.
//
// Note: [http.NewClient] wraps the provided transport with OpenTelemetry instrumentation when tracing
// or metrics are enabled.
func NewClient(opts ...ClientOption) (*http.Client, error) {
	resolved := options(opts...)

	transport, err := NewRoundTripper(opts...)
	if err != nil {
		return nil, err
	}

	client := http.NewClient(transport, resolved.timeout)
	if resolved.gen != nil {
		client.CheckRedirect = http.SameOriginRedirect
	}

	return client, nil
}

func roundTripper(resolved *clientOpts) (http.RoundTripper, error) {
	hrt := resolved.roundTripper
	if hrt != nil {
		return hrt, nil
	}

	if !resolved.tls.IsEnabled() {
		return http.Transport(nil), nil
	}

	conf, err := client.NewConfig(fs, resolved.tls)
	if err != nil {
		return nil, err
	}

	return http.Transport(conf), nil
}

func options(opts ...ClientOption) *clientOpts {
	resolved := &clientOpts{}
	for _, o := range opts {
		o.apply(resolved)
	}

	if resolved.timeout <= 0 {
		resolved.timeout = time.DefaultTimeout
	}

	if resolved.generator == nil {
		resolved.generator = uuid.NewGenerator()
	}

	return resolved
}
