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
)

// Client bundles the dependencies needed to construct instrumented HTTP and gRPC clients for tests.
type Client struct {
	Lifecycle    di.Lifecycle
	Logger       *logger.Logger
	Tracer       *tracer.Config
	Transport    *transport.Config
	TLS          *tls.Config
	ID           id.Generator
	RoundTripper http.RoundTripper
	Meter        metrics.Meter
	Generator    token.Generator
	HTTPLimiter  *httplimiter.Client
	GRPCLimiter  *grpclimiter.Client
	Compression  bool
}

// NewHTTP returns an HTTP client configured with the world's logger, retry policy,
// token generator, limiter, tracing, and optional compression.
func (c *Client) NewHTTP(os ...httpbreaker.Option) (*http.Client, error) {
	opts := []http.ClientOption{
		http.WithClientLogger(c.Logger),
		http.WithClientRoundTripper(c.RoundTripper),
		http.WithClientBreaker(os...),
		http.WithClientRetry(c.Transport.HTTP.Retry),
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
// retry policy, token generator, limiter, tracing, and optional compression.
func (c *Client) NewGRPC(os ...grpcbreaker.Option) (*grpc.ClientConn, error) {
	opts := []grpc.ClientOption{
		grpc.WithClientUnaryInterceptors(),
		grpc.WithClientStreamInterceptors(),
		grpc.WithClientLogger(c.Logger),
		grpc.WithClientBreaker(os...),
		grpc.WithClientRetry(c.Transport.GRPC.Retry),
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
