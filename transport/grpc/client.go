package grpc

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	gtel "github.com/alexfalkowski/go-service/transport/grpc/telemetry"
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
	logger   *zap.Logger
	tracer   gtel.Tracer
	metrics  *gtel.ClientMetrics
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
func WithClientTracer(tracer gtel.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// WithClientConfig for gRPC.
func WithClientMetrics(metrics *gtel.ClientMetrics) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.metrics = metrics
	})
}

// ClientParams for gRPC.
type ClientParams struct {
	Context context.Context
	Host    string
	Config  *Config
}

// NewClient to host for gRPC.
func NewClient(params ClientParams, opts ...ClientOption) (*grpc.ClientConn, error) {
	defaultOptions := &clientOptions{
		security: grpc.WithTransportCredentials(insecure.NewCredentials()),
		tracer:   telemetry.NewNoopTracer("grpc"),
	}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	grpcOpts := []grpc.DialOption{}
	grpcOpts = append(grpcOpts, unaryDialOption(params, defaultOptions), streamDialOption(params, defaultOptions), defaultOptions.security)
	grpcOpts = append(grpcOpts, defaultOptions.opts...)

	return grpc.DialContext(params.Context, params.Host, grpcOpts...)
}

func unaryDialOption(params ClientParams, opts *clientOptions) grpc.DialOption {
	unary := []grpc.UnaryClientInterceptor{}

	if opts.retry {
		unary = append(unary,
			retry.UnaryClientInterceptor(
				retry.WithCodes(codes.Unavailable, codes.DataLoss),
				retry.WithMax(params.Config.Retry.Attempts),
				retry.WithBackoff(retry.BackoffLinear(backoffLinear)),
				retry.WithPerRetryTimeout(params.Config.Retry.Timeout),
			),
		)
	}

	if opts.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor())
	}

	unary = append(unary, meta.UnaryClientInterceptor(params.Config.UserAgent))

	if opts.logger != nil {
		unary = append(unary, gtel.LoggerUnaryClientInterceptor(opts.logger))
	}

	if opts.metrics != nil {
		unary = append(unary, opts.metrics.UnaryClientInterceptor())
	}

	unary = append(unary, gtel.TracerUnaryClientInterceptor(opts.tracer))
	unary = append(unary, opts.unary...)

	return grpc.WithChainUnaryInterceptor(unary...)
}

func streamDialOption(params ClientParams, opts *clientOptions) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{meta.StreamClientInterceptor(params.Config.UserAgent)}

	if opts.logger != nil {
		stream = append(stream, gtel.LoggerStreamClientInterceptor(opts.logger))
	}

	if opts.metrics != nil {
		stream = append(stream, opts.metrics.StreamClientInterceptor())
	}

	stream = append(stream, gtel.TracerStreamClientInterceptor(opts.tracer))
	stream = append(stream, opts.stream...)

	return grpc.WithChainStreamInterceptor(stream...)
}
