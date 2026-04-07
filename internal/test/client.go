package test

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	grpcbreaker "github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	grpclimiter "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	transporthttp "github.com/alexfalkowski/go-service/v2/transport/http"
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

// HTTP returns an HTTP client configured with the world's logger, retry policy,
// token generator, limiter, tracing, and optional compression.
func (c *Client) HTTP(os ...httpbreaker.Option) (*http.Client, error) {
	opts := []transporthttp.ClientOption{
		transporthttp.WithClientLogger(c.Logger),
		transporthttp.WithClientRoundTripper(c.RoundTripper),
		transporthttp.WithClientBreaker(os...),
		transporthttp.WithClientRetry(c.Transport.HTTP.Retry),
		transporthttp.WithClientUserAgent(UserAgent),
		transporthttp.WithClientTokenGenerator(UserID, c.Generator),
		transporthttp.WithClientTimeout(time.Minute),
		transporthttp.WithClientTLS(c.TLS),
		transporthttp.WithClientID(c.ID),
		transporthttp.WithClientLimiter(c.HTTPLimiter),
	}

	if c.Compression {
		opts = append(opts, transporthttp.WithClientCompression())
	}

	return transporthttp.NewClient(opts...)
}

// NewHTTP returns an HTTP client configured with the world's logger, retry policy,
// token generator, limiter, tracing, and optional compression.
func (c *Client) NewHTTP(os ...httpbreaker.Option) (*http.Client, error) {
	return c.HTTP(os...)
}

// GRPC returns a gRPC client connection configured with the world's interceptors,
// retry policy, token generator, limiter, tracing, and optional compression.
func (c *Client) GRPC(os ...grpcbreaker.Option) (*grpc.ClientConn, error) {
	opts := []transportgrpc.ClientOption{
		transportgrpc.WithClientUnaryInterceptors(),
		transportgrpc.WithClientStreamInterceptors(),
		transportgrpc.WithClientLogger(c.Logger),
		transportgrpc.WithClientBreaker(os...),
		transportgrpc.WithClientRetry(c.Transport.GRPC.Retry),
		transportgrpc.WithClientUserAgent(UserAgent),
		transportgrpc.WithClientTokenGenerator(UserID, c.Generator),
		transportgrpc.WithClientTimeout(time.Minute),
		transportgrpc.WithClientKeepalive(time.Minute, time.Minute),
		transportgrpc.WithClientDialOption(),
		transportgrpc.WithClientTLS(c.TLS),
		transportgrpc.WithClientID(c.ID),
		transportgrpc.WithClientLimiter(c.GRPCLimiter),
	}

	if c.Compression {
		opts = append(opts, transportgrpc.WithClientCompression())
	}

	_, target := net.ListenNetworkAddress(c.Transport.GRPC.Address)

	return transportgrpc.NewClient(target, opts...)
}

// NewGRPC returns a gRPC client connection configured with the world's interceptors,
// retry policy, token generator, limiter, tracing, and optional compression.
func (c *Client) NewGRPC(os ...grpcbreaker.Option) (*grpc.ClientConn, error) {
	return c.GRPC(os...)
}
