package grpc

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// ClientOption for gRPC.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	generator   id.Generator
	gen         token.Generator
	meter       *metrics.Meter
	security    *tls.Config
	logger      *logger.Logger
	retry       *retry.Config
	tracer      *tracer.Tracer
	limiter     *limiter.Client
	userAgent   env.UserAgent
	id          env.UserID
	opts        []grpc.DialOption
	unary       []grpc.UnaryClientInterceptor
	stream      []grpc.StreamClientInterceptor
	timeout     time.Duration
	breaker     bool
	compression bool
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithClientCompression for gRPC.
func WithClientCompression() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.compression = true
	})
}

// WithClientTokenGenerator for gRPC.
func WithClientTokenGenerator(id env.UserID, gen token.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.id = id
		o.gen = gen
	})
}

// WithClientTimeout for gRPC.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = time.MustParseDuration(timeout)
	})
}

// WithClientRetry for gRPC.
func WithClientRetry(cfg *retry.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.retry = cfg
	})
}

// WithClientBreaker for gRPC.
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.breaker = true
	})
}

// WithClientTLS for gRPC.
func WithClientTLS(sec *tls.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.security = sec
	})
}

// WithClientDialOption for gRPC.
func WithClientDialOption(opts ...grpc.DialOption) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.opts = opts
	})
}

// WithClientUnaryInterceptors for gRPC.
func WithClientUnaryInterceptors(unary ...grpc.UnaryClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.unary = unary
	})
}

// WithClientStreamInterceptors for gRPC.
func WithClientStreamInterceptors(stream ...grpc.StreamClientInterceptor) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.stream = stream
	})
}

// WithClientLogger for gRPC.
func WithClientLogger(logger *logger.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// WithClientTracer for gRPC.
func WithClientTracer(tracer *tracer.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tracer = tracer
	})
}

// WithClientMetrics for gRPC.
func WithClientMetrics(meter *metrics.Meter) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.meter = meter
	})
}

// WithClientUserAgent for gRPC.
func WithClientUserAgent(userAgent env.UserAgent) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// WithClientID for gRPC.
func WithClientID(generator id.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.generator = generator
	})
}

// WithClientLimiter for gRPC.
func WithClientLimiter(limiter *limiter.Client) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.limiter = limiter
	})
}

// NewDialOptions for gRPC.
func NewDialOptions(opts ...ClientOption) ([]grpc.DialOption, error) {
	os := options(opts...)

	var security grpc.DialOption
	if os.security.IsEnabled() {
		conf, err := tls.NewConfig(fs, os.security)
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
		grpc.WithKeepaliveParams(os.timeout),
		grpc.WithChainUnaryInterceptor(cis...), sto, security,
	}

	if os.compression {
		ops = append(ops, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	ops = append(ops, os.opts...)

	return ops, nil
}

// ClientConn is an alias for grpc.ClientConn.
type ClientConn = grpc.ClientConn

// NewClient for gRPC.
func NewClient(target string, opts ...ClientOption) (*ClientConn, error) {
	os, err := NewDialOptions(opts...)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(target, os...)
}

// UnaryClientInterceptors for gRPC.
func UnaryClientInterceptors(opts ...ClientOption) []grpc.UnaryClientInterceptor {
	os := options(opts...)
	unary := []grpc.UnaryClientInterceptor{}

	unary = append(unary, os.unary...)
	unary = append(unary, grpc.TimeoutUnaryClientInterceptor(os.timeout))

	if os.limiter != nil {
		unary = append(unary, limiter.UnaryClientInterceptor(os.limiter))
	}

	if os.retry != nil {
		unary = append(unary, retry.UnaryClientInterceptor(os.retry))
	}

	if os.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor())
	}

	if os.logger != nil {
		unary = append(unary, logger.UnaryClientInterceptor(os.logger))
	}

	if os.meter != nil {
		unary = append(unary, metrics.NewClient(os.meter).UnaryInterceptor())
	}

	if os.tracer != nil {
		unary = append(unary, tracer.UnaryClientInterceptor(os.tracer))
	}

	if os.gen != nil {
		unary = append(unary, token.UnaryClientInterceptor(os.id, os.gen))
	}

	unary = append(unary, meta.UnaryClientInterceptor(os.userAgent, os.generator))
	return unary
}

func streamDialOption(opts *clientOpts) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{}
	stream = append(stream, opts.stream...)

	if opts.logger != nil {
		stream = append(stream, logger.StreamClientInterceptor(opts.logger))
	}

	if opts.meter != nil {
		stream = append(stream, metrics.NewClient(opts.meter).StreamInterceptor())
	}

	if opts.tracer != nil {
		stream = append(stream, tracer.StreamClientInterceptor(opts.tracer))
	}

	if opts.gen != nil {
		stream = append(stream, token.StreamClientInterceptor(opts.id, opts.gen))
	}

	stream = append(stream, meta.StreamClientInterceptor(opts.userAgent, opts.generator))
	return grpc.WithChainStreamInterceptor(stream...)
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}
	if os.timeout == 0 {
		os.timeout = 30 * time.Second
	}
	if os.generator == nil {
		os.generator = &id.UUID{}
	}

	return os
}
