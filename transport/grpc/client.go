package grpc

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics/prometheus"
	gtracer "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
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
	tracer   gtracer.Tracer
	metrics  *prometheus.ClientCollector
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
func WithClientTracer(tracer gtracer.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// WithClientConfig for gRPC.
func WithClientMetrics(metrics *prometheus.ClientCollector) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.metrics = metrics
	})
}

// NewClient to host for gRPC.
func NewClient(ctx context.Context, host string, cfg *Config, opts ...ClientOption) (*grpc.ClientConn, error) {
	defaultOptions := &clientOptions{
		security: grpc.WithTransportCredentials(insecure.NewCredentials()),
		tracer:   tracer.NewNoopTracer("grpc"),
	}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	grpcOpts := []grpc.DialOption{}
	grpcOpts = append(grpcOpts, unaryDialOption(cfg, defaultOptions), streamDialOption(cfg, defaultOptions), defaultOptions.security)
	grpcOpts = append(grpcOpts, defaultOptions.opts...)

	return grpc.DialContext(ctx, host, grpcOpts...)
}

func unaryDialOption(cfg *Config, opts *clientOptions) grpc.DialOption {
	unary := []grpc.UnaryClientInterceptor{}

	if opts.retry {
		unary = append(unary,
			retry.UnaryClientInterceptor(
				retry.WithCodes(codes.Unavailable, codes.DataLoss),
				retry.WithMax(cfg.Retry.Attempts),
				retry.WithBackoff(retry.BackoffLinear(backoffLinear)),
				retry.WithPerRetryTimeout(cfg.Retry.Timeout),
			),
		)
	}

	if opts.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor())
	}

	unary = append(unary, meta.UnaryClientInterceptor(cfg.UserAgent))

	if opts.logger != nil {
		unary = append(unary, szap.UnaryClientInterceptor(opts.logger))
	}

	if opts.metrics != nil {
		unary = append(unary, opts.metrics.UnaryClientInterceptor())
	}

	unary = append(unary, gtracer.UnaryClientInterceptor(opts.tracer))
	unary = append(unary, opts.unary...)

	return grpc.WithChainUnaryInterceptor(unary...)
}

func streamDialOption(cfg *Config, opts *clientOptions) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{meta.StreamClientInterceptor(cfg.UserAgent)}

	if opts.logger != nil {
		stream = append(stream, szap.StreamClientInterceptor(opts.logger))
	}

	if opts.metrics != nil {
		stream = append(stream, opts.metrics.StreamClientInterceptor())
	}

	stream = append(stream, gtracer.StreamClientInterceptor(opts.tracer))
	stream = append(stream, opts.stream...)

	return grpc.WithChainStreamInterceptor(stream...)
}
