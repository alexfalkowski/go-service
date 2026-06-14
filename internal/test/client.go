package test

import (
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	grpcbreaker "github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	grpclimiter "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	httpbreaker "github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	httplimiter "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/retry"
)

// Client bundles the dependencies needed to construct instrumented HTTP and gRPC clients for tests.
type Client struct {
	// Lifecycle receives client-related lifecycle hooks.
	Lifecycle di.Lifecycle
	// Logger logs client activity.
	Logger *logger.Logger
	// Tracer is retained for world configuration symmetry. Client tracing follows
	// the globally registered telemetry providers.
	Tracer *tracer.Config
	// Transport configures HTTP and gRPC client targets.
	Transport *transport.Config
	// TLS configures client TLS.
	TLS *tls.Config
	// ID generates request identifiers.
	ID id.Generator
	// RoundTripper is the base HTTP transport for world clients.
	RoundTripper http.RoundTripper
	// Meter provides metrics instrumentation for clients.
	Meter metrics.Meter
	// Generator generates outbound transport tokens.
	Generator token.Generator
	// Retry configures client retries.
	Retry *retry.Config
	// HTTPLimiter limits outbound HTTP requests.
	HTTPLimiter *httplimiter.Client
	// GRPCLimiter limits outbound gRPC requests.
	GRPCLimiter *grpclimiter.Client
	// Compression enables transport compression for world clients.
	Compression bool
}

// NewHTTP returns an HTTP client configured with the world's logger, retry policy,
// token generator, limiter, and optional compression.
func (c *Client) NewHTTP(os ...httpbreaker.Option) (*http.Client, error) {
	opts := []http.ClientOption{
		http.WithClientLogger(c.Logger),
		http.WithClientRoundTripper(c.RoundTripper),
		http.WithClientBreaker(os...),
		http.WithClientRetry(c.Retry),
		http.WithClientUserAgent(UserAgent),
		http.WithClientTokenGenerator(UserID, c.Generator),
		http.WithClientTimeout(time.Minute),
		http.WithClientTLS(c.TLS),
		http.WithClientID(c.ID),
		http.WithClientLimiter(c.HTTPLimiter),
	}

	if c.Compression {
		opts = append(opts, http.WithClientCompression())
	}

	return http.NewClient(opts...)
}

// NewGRPC returns a gRPC client connection configured with the world's interceptors,
// retry policy, token generator, limiter, and optional compression.
func (c *Client) NewGRPC(os ...grpcbreaker.Option) (*grpc.ClientConn, error) {
	opts := []grpc.ClientOption{
		grpc.WithClientUnaryInterceptors(),
		grpc.WithClientStreamInterceptors(),
		grpc.WithClientLogger(c.Logger),
		grpc.WithClientBreaker(os...),
		grpc.WithClientRetry(c.Retry),
		grpc.WithClientUserAgent(UserAgent),
		grpc.WithClientTokenGenerator(UserID, c.Generator),
		grpc.WithClientTimeout(time.Minute),
		grpc.WithClientKeepalive(time.Minute, time.Minute),
		grpc.WithClientDialOption(),
		grpc.WithClientTLS(c.TLS),
		grpc.WithClientID(c.ID),
		grpc.WithClientLimiter(c.GRPCLimiter),
	}

	if c.Compression {
		opts = append(opts, grpc.WithClientCompression())
	}

	_, target := net.ListenNetworkAddress(c.Transport.GRPC.Address)

	return grpc.NewClient(target, opts...)
}
