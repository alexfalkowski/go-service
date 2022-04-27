package grpc

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	szap "github.com/alexfalkowski/go-service/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	backoffLinear = 50 * time.Millisecond
)

// ClientOption for HTTP.
type ClientOption interface{ apply(*clientOptions) }

type clientOptions struct {
	config   *Config
	logger   *zap.Logger
	tracer   opentracing.Tracer
	retry    bool
	breaker  bool
	opts     []grpc.DialOption
	unary    []grpc.UnaryClientInterceptor
	stream   []grpc.StreamClientInterceptor
	security grpc.DialOption
	version  version.Version
}

type clientOptionFunc func(*clientOptions)

func (f clientOptionFunc) apply(o *clientOptions) { f(o) }

// WithClientRetry for gRPC.
func WithClientRetry() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.retry = true
	})
}

// WithClientBreaker for gRPC.
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.breaker = true
	})
}

// WithClientSecure for gRPC.
func WithClientSecure() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.security = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	})
}

// WithClientDialOption for gRPC.
func WithClientDialOption(opts ...grpc.DialOption) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.opts = opts
	})
}

// WithClientUnaryInterceptors for gRPC.
func WithClientUnaryInterceptors(unary ...grpc.UnaryClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.unary = unary
	})
}

// WithClientUnaryInterceptors for gRPC.
func WithClientStreamInterceptors(stream ...grpc.StreamClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.stream = stream
	})
}

// WithClientLogger for gRPC.
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.logger = logger
	})
}

// WithClientConfig for gRPC.
func WithClientConfig(config *Config) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.config = config
	})
}

// WithClientConfig for gRPC.
func WithClientTracer(tracer opentracing.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// WithClientVersion for gRPC.
func WithClientVersion(version version.Version) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.version = version
	})
}

// NewClient to host for gRPC.
func NewClient(context context.Context, host string, opts ...ClientOption) (*grpc.ClientConn, error) {
	defaultOptions := &clientOptions{security: grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	grpcOpts := []grpc.DialOption{}
	grpcOpts = append(grpcOpts, unaryDialOption(defaultOptions), streamDialOption(defaultOptions), defaultOptions.security)
	grpcOpts = append(grpcOpts, defaultOptions.opts...)

	return grpc.DialContext(context, host, grpcOpts...)
}

func unaryDialOption(opts *clientOptions) grpc.DialOption {
	unary := []grpc.UnaryClientInterceptor{}

	if opts.retry {
		unary = append(unary,
			retry.UnaryClientInterceptor(
				retry.WithCodes(codes.Unavailable, codes.DataLoss),
				retry.WithMax(opts.config.Retry.Attempts),
				retry.WithBackoff(retry.BackoffLinear(backoffLinear)),
				retry.WithPerRetryTimeout(opts.config.Retry.Timeout),
			),
		)
	}

	if opts.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor())
	}

	unary = append(unary,
		meta.UnaryClientInterceptor(opts.config.UserAgent, opts.version),
		szap.UnaryClientInterceptor(opts.logger),
		opentracing.UnaryClientInterceptor(opts.tracer),
	)

	unary = append(unary, opts.unary...)

	return grpc.WithChainUnaryInterceptor(unary...)
}

func streamDialOption(opts *clientOptions) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{
		meta.StreamClientInterceptor(opts.config.UserAgent, opts.version),
		szap.StreamClientInterceptor(opts.logger),
		opentracing.StreamClientInterceptor(opts.tracer),
	}

	stream = append(stream, opts.stream...)

	return grpc.WithChainStreamInterceptor(stream...)
}
