package grpc

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	gtracer "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	backoffLinear = 50 * time.Millisecond
)

// ClientOption for HTTP.
type ClientOption interface{ apply(*clientOptions) }

type clientOptions struct {
	logger   *zap.Logger
	tracer   gtracer.Tracer
	meter    metric.Meter
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
func WithClientSecure(sec security.Config) (ClientOption, error) {
	var creds credentials.TransportCredentials

	if sec.IsClientEnabled() {
		conf, err := security.ClientTLSConfig(sec)
		if err != nil {
			return nil, err
		}

		creds = credentials.NewTLS(conf)
	} else {
		creds = credentials.NewClientTLSFromCert(nil, "")
	}

	opt := clientOptionFunc(func(o *clientOptions) {
		o.security = grpc.WithTransportCredentials(creds)
	})

	return opt, nil
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

// WithClientTracer for gRPC.
func WithClientTracer(tracer gtracer.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// WithClientMetrics for gRPC.
func WithClientMetrics(meter metric.Meter) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.meter = meter
	})
}

// NewDialOptions for gRPC.
func NewDialOptions(cfg *Config, opts ...ClientOption) ([]grpc.DialOption, error) {
	defaultOptions := &clientOptions{
		security: grpc.WithTransportCredentials(insecure.NewCredentials()),
		tracer:   tracer.NewNoopTracer("grpc"),
	}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	udo, err := unaryDialOption(cfg, defaultOptions)
	if err != nil {
		return nil, err
	}

	sto, err := streamDialOption(cfg, defaultOptions)
	if err != nil {
		return nil, err
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithUserAgent(cfg.UserAgent),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	grpcOpts = append(grpcOpts, udo, sto, defaultOptions.security)
	grpcOpts = append(grpcOpts, defaultOptions.opts...)

	return grpcOpts, nil
}

// NewClient to host for gRPC.
func NewClient(ctx context.Context, host string, cfg *Config, opts ...ClientOption) (*grpc.ClientConn, error) {
	os, err := NewDialOptions(cfg, opts...)
	if err != nil {
		return nil, err
	}

	return grpc.DialContext(ctx, host, os...)
}

func unaryDialOption(cfg *Config, opts *clientOptions) (grpc.DialOption, error) {
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

	if opts.meter != nil {
		client, err := metrics.NewClient(opts.meter)
		if err != nil {
			return nil, err
		}

		unary = append(unary, client.UnaryInterceptor())
	}

	unary = append(unary, gtracer.UnaryClientInterceptor(opts.tracer))
	unary = append(unary, opts.unary...)

	return grpc.WithChainUnaryInterceptor(unary...), nil
}

func streamDialOption(cfg *Config, opts *clientOptions) (grpc.DialOption, error) {
	stream := []grpc.StreamClientInterceptor{meta.StreamClientInterceptor(cfg.UserAgent)}

	if opts.logger != nil {
		stream = append(stream, szap.StreamClientInterceptor(opts.logger))
	}

	if opts.meter != nil {
		client, err := metrics.NewClient(opts.meter)
		if err != nil {
			return nil, err
		}

		stream = append(stream, client.StreamInterceptor())
	}

	stream = append(stream, gtracer.StreamClientInterceptor(opts.tracer))
	stream = append(stream, opts.stream...)

	return grpc.WithChainStreamInterceptor(stream...), nil
}
