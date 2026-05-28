package grpc

import (
	"github.com/alexfalkowski/go-service/v2/config/client"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// ClientOption configures gRPC client construction.
//
// Client options are applied when building dial options (see `NewDialOptions`) and when assembling
// interceptor chains (see `UnaryClientInterceptors`).
//
// Most options are orthogonal and can be combined. Some options enable behavior by providing a non-nil
// dependency (for example, retries are enabled when `WithClientRetry` provides a non-nil config).
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	gen               token.Generator
	generator         id.Generator
	security          *tls.Config
	logger            *logger.Logger
	retry             *retry.Config
	retryPolicies     []retry.Policy
	limiter           *limiter.Client
	userAgent         env.UserAgent
	id                env.UserID
	opts              []grpc.DialOption
	unary             []grpc.UnaryClientInterceptor
	stream            []grpc.StreamClientInterceptor
	breakerOptions    []breaker.Option
	keepalive_ping    time.Duration
	keepalive_timeout time.Duration
	timeout           time.Duration
	breaker           bool
	compression       bool
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithClientCompression enables gzip compression for gRPC client calls.
//
// This option appends a default call option that requests the "gzip" compressor. The server must also be
// configured to accept gzip-compressed requests for this to have any effect.
func WithClientCompression() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.compression = true
	})
}

// WithClientTokenGenerator enables client-side token injection interceptors.
//
// When configured, the client will generate an Authorization token per RPC and attach it to outgoing
// metadata (unary and streaming). The token is generated via gen and is typically scoped to the RPC's
// full method name and the provided user id.
func WithClientTokenGenerator(id env.UserID, gen token.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.id = id
		o.gen = gen
	})
}

// WithClientTimeout sets the default per-RPC timeout applied by the unary timeout interceptor.
//
// If unset or negative, a default timeout is applied (see `NewDialOptions` defaults).
//
// Note: this timeout is enforced via an interceptor and is independent from any deadlines already set
// on the incoming context; the interceptor will typically only apply a timeout when a deadline is not
// already present. Streaming callers should use explicit context deadlines or a custom stream interceptor.
func WithClientTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = timeout
	})
}

// WithClientKeepalive sets gRPC client keepalive ping and timeout parameters.
//
// Keepalive affects connection liveness detection and can help discover broken connections.
//
// If either value is unset (0), it defaults to the resolved client timeout (see `NewDialOptions`).
func WithClientKeepalive(ping, timeout time.Duration) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.keepalive_ping = ping
		o.keepalive_timeout = timeout
	})
}

// WithClientRetry enables retry behavior for unary client calls.
//
// Retries are applied via a unary client interceptor. The retry policy is derived from cfg and
// typically includes a maximum attempt count, per-retry timeout, and a backoff strategy.
// Optional policies decide whether a logical unary RPC is eligible for retry.
//
// When no policy is provided, only side-effect-safe unary RPCs are eligible for retry: AIP-style read methods,
// or calls carrying a request-id idempotency contract. Callers that need different behavior can pass an
// explicit policy.
//
// If cfg is nil, retries are not enabled.
func WithClientRetry(cfg *retry.Config, policies ...retry.Policy) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.retry = cfg
		o.retryPolicies = policies
	})
}

// WithClientBreaker enables circuit breaking for unary client calls.
//
// Circuit breakers are keyed per RPC full method name. Failure accounting is controlled by the
// breaker options (for example, which gRPC status codes count as failures). Streaming callers should
// use a custom stream interceptor for stream-specific breaker behavior.
func WithClientBreaker(opts ...breaker.Option) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.breaker = true
		o.breakerOptions = opts
	})
}

// WithClientTLS enables TLS for the client connection using sec.
//
// TLS configuration is materialized when building dial options. A non-nil sec enables TLS, including an
// empty config that relies on system roots and the target host name. Source strings may be resolved via the
// package-registered filesystem (see the package `Register` function).
func WithClientTLS(sec *tls.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.security = sec
	})
}

// WithClientDialOption configures raw gRPC dial options.
//
// This is an escape hatch for passing options not modeled by this package. NewDialOptions appends the
// provided options after the package's baseline options, so callers can override behavior when supported
// by gRPC.
//
// Pass all custom dial options for a client construction in one call. Repeating this option follows the
// package's last-wins functional option convention and replaces earlier raw dial options.
func WithClientDialOption(opts ...grpc.DialOption) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.opts = opts
	})
}

// WithClientUnaryInterceptors adds custom unary client interceptors after metadata propagation.
//
// Metadata propagation runs first so custom interceptors see the resolved user-agent and request-id.
// Interceptors provided here then run before the remaining standard interceptors added by this package
// (timeout, retry, breaker, logging, token injection, etc.).
//
// Pass all custom unary interceptors for a client construction in one call. Repeating this option follows
// the package's last-wins functional option convention and replaces earlier custom unary interceptors.
func WithClientUnaryInterceptors(unary ...grpc.UnaryClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.unary = unary
	})
}

// WithClientStreamInterceptors adds custom stream client interceptors after metadata propagation.
//
// Metadata propagation runs first so custom interceptors see the resolved user-agent and request-id.
// Interceptors provided here then run before the remaining standard interceptors added by this package
// (logging, token injection, limiter, etc.).
//
// Pass all custom stream interceptors for a client construction in one call. Repeating this option follows
// the package's last-wins functional option convention and replaces earlier custom stream interceptors.
func WithClientStreamInterceptors(stream ...grpc.StreamClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.stream = stream
	})
}

// WithClientLogger enables gRPC client logging interceptors.
//
// When configured, both unary and stream client interceptors may emit logs about RPC outcomes.
func WithClientLogger(logger *logger.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// WithClientUserAgent sets the user agent string used for the gRPC connection and metadata propagation.
//
// The value is used in two places:
//   - as the gRPC dial user agent (`grpc.WithUserAgent`)
//   - for metadata propagation via the `net/grpc/meta` interceptors
func WithClientUserAgent(userAgent env.UserAgent) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// WithClientID sets the request id generator used by metadata propagation interceptors.
//
// The generator is used to create a request id when one is not already present on the outgoing context
// or outgoing metadata.
func WithClientID(generator id.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.generator = generator
	})
}

// WithClientLimiter enables client-side rate limiting interceptors.
//
// When configured, unary client calls and streams are rate-limited before being sent. If limiter is nil,
// rate limiting is not enabled.
func WithClientLimiter(limiter *limiter.Client) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.limiter = limiter
	})
}

// NewDialOptions builds `grpc.DialOption` values from `ClientOption`.
//
// Defaults (see `options()`):
//   - timeout: 30s
//   - request-id generator: uuid
//
// Keepalive defaults:
//   - if keepalive ping or timeout are not set (0), they default to the resolved timeout.
//
// Transport security:
//   - if TLS is requested, TLS config is constructed using the package-registered filesystem (see `Register`)
//     to resolve TLS source strings.
//   - otherwise, insecure transport credentials are used.
//     This default is intended for local and in-cluster/container traffic where transport security is provided
//     at the platform boundary; use WithClientTLS for calls outside that trusted boundary.
func NewDialOptions(opts ...ClientOption) ([]grpc.DialOption, error) {
	os := options(opts...)

	if os.keepalive_ping == 0 {
		os.keepalive_ping = os.timeout
	}

	if os.keepalive_timeout == 0 {
		os.keepalive_timeout = os.timeout
	}

	var security grpc.DialOption
	if os.security != nil {
		conf, err := client.NewConfig(fs, os.security)
		if err != nil {
			return nil, err
		}

		security = grpc.WithTransportCredentials(grpc.NewTLS(conf))
	} else {
		security = grpc.WithTransportCredentials(grpc.NewInsecureCredentials())
	}

	cis := UnaryClientInterceptors(opts...)
	sto := streamDialOption(os)
	ops := []grpc.DialOption{
		grpc.WithUserAgent(os.userAgent.String()),
		grpc.WithKeepaliveParams(os.keepalive_ping, os.keepalive_timeout),
		grpc.WithChainUnaryInterceptor(cis...), sto, security,
	}

	if os.compression {
		ops = append(ops, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	ops = append(ops, os.opts...)

	return ops, nil
}

// ClientConn is an alias for grpc.ClientConn.
//
// It is exposed so callers can refer to the concrete connection type without importing the underlying
// gRPC wrapper package directly.
type ClientConn = grpc.ClientConn

// NewClient constructs and dials a gRPC client connection to target.
//
// It uses dial options derived from opts (see `NewDialOptions`) and adds OpenTelemetry stats handling
// when tracing or metrics are enabled.
//
// The returned connection should be closed by the caller when no longer needed.
func NewClient(target string, opts ...ClientOption) (*ClientConn, error) {
	os, err := NewDialOptions(opts...)
	if err != nil {
		return nil, err
	}

	if metrics.IsEnabled() || tracer.IsEnabled() {
		os = append(os, grpc.WithStatsHandler(telemetry.NewClientHandler()))
	}

	return grpc.NewClient(target, os...)
}

// UnaryClientInterceptors builds the unary client interceptor chain derived from opts.
//
// Order matters. Interceptors are appended in the following sequence:
//   - metadata propagation interceptor (user-agent and request-id)
//   - any custom interceptors provided via `WithClientUnaryInterceptors`
//   - a timeout interceptor
//   - optional limiter interceptor (when configured)
//   - optional circuit breaker interceptor (when enabled via `WithClientBreaker`)
//   - optional retry interceptor (when configured)
//   - optional logging interceptor (when configured)
//   - optional token injection interceptor (when configured)
//
// The limiter stays before the breaker so local quota denials are not counted as upstream failures.
func UnaryClientInterceptors(opts ...ClientOption) []grpc.UnaryClientInterceptor {
	os := options(opts...)
	unary := []grpc.UnaryClientInterceptor{}

	unary = append(unary, meta.UnaryClientInterceptor(os.userAgent, os.generator))
	unary = append(unary, os.unary...)
	unary = append(unary, grpc.TimeoutUnaryClientInterceptor(os.timeout))

	if os.limiter != nil {
		unary = append(unary, limiter.UnaryClientInterceptor(os.limiter))
	}

	if os.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor(os.breakerOptions...))
	}

	if os.retry != nil {
		unary = append(unary, retry.UnaryClientInterceptor(os.retry, os.retryPolicies...))
	}

	if os.logger != nil {
		unary = append(unary, logger.UnaryClientInterceptor(os.logger))
	}

	if os.gen != nil {
		unary = append(unary, token.UnaryClientInterceptor(os.id, os.gen))
	}

	return unary
}

func streamDialOption(opts *clientOpts) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{}
	stream = append(stream, meta.StreamClientInterceptor(opts.userAgent, opts.generator))
	stream = append(stream, opts.stream...)

	if opts.logger != nil {
		stream = append(stream, logger.StreamClientInterceptor(opts.logger))
	}

	if opts.gen != nil {
		stream = append(stream, token.StreamClientInterceptor(opts.id, opts.gen))
	}

	if opts.limiter != nil {
		stream = append(stream, limiter.StreamClientInterceptor(opts.limiter))
	}

	return grpc.WithChainStreamInterceptor(stream...)
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}
	if os.timeout <= 0 {
		os.timeout = time.DefaultTimeout
	}
	if os.generator == nil {
		os.generator = uuid.NewGenerator()
	}

	return os
}
