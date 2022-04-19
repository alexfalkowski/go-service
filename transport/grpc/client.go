package grpc

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	szap "github.com/alexfalkowski/go-service/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	gopentracing "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/opentracing/opentracing-go"
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
}

type clientOptionFunc func(*clientOptions)

func (f clientOptionFunc) apply(o *clientOptions) { f(o) }

// WithClientRetry for gRPC.
// nolint:ireturn
func WithClientRetry() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.retry = true
	})
}

// WithClientBreaker for gRPC.
// nolint:ireturn
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.breaker = true
	})
}

// WithClientSecure for gRPC.
// nolint:ireturn
func WithClientSecure() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.security = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	})
}

// WithClientDialOption for gRPC.
// nolint:ireturn
func WithClientDialOption(opts ...grpc.DialOption) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.opts = opts
	})
}

// WithClientUnaryInterceptors for gRPC.
// nolint:ireturn
func WithClientUnaryInterceptors(unary ...grpc.UnaryClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.unary = unary
	})
}

// WithClientUnaryInterceptors for gRPC.
// nolint:ireturn
func WithClientStreamInterceptors(stream ...grpc.StreamClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.stream = stream
	})
}

// WithClientLogger for gRPC.
// nolint:ireturn
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.logger = logger
	})
}

// WithClientConfig for gRPC.
// nolint:ireturn
func WithClientConfig(config *Config) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.config = config
	})
}

// WithClientConfig for gRPC.
// nolint:ireturn
func WithClientTracer(tracer opentracing.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
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

// nolint:ireturn
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
		meta.UnaryClientInterceptor(opts.config.UserAgent),
		szap.UnaryClientInterceptor(opts.logger),
		gopentracing.UnaryClientInterceptor(opts.tracer),
	)

	unary = append(unary, opts.unary...)

	return grpc.WithChainUnaryInterceptor(unary...)
}

// nolint:ireturn
func streamDialOption(opts *clientOptions) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{
		meta.StreamClientInterceptor(opts.config.UserAgent),
		szap.StreamClientInterceptor(opts.logger),
		gopentracing.StreamClientInterceptor(opts.tracer),
	}

	stream = append(stream, opts.stream...)

	return grpc.WithChainStreamInterceptor(stream...)
}
