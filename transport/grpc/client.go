package grpc

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	tl "github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	tm "github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/metrics"
	tt "github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/tracer"
	tkn "github.com/alexfalkowski/go-service/v2/transport/grpc/token"
	ri "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	ti "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ClientOption for gRPC.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	tracer      *tracer.Tracer
	meter       *metrics.Meter
	security    *tls.Config
	gen         token.Generator
	logger      *logger.Logger
	retry       *retry.Config
	userAgent   env.UserAgent
	generator   id.Generator
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

// WithClientBreaker for gRPC.
func WithClientCompression() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.compression = true
	})
}

// WithClientTokenGenerator for gRPC.
func WithClientTokenGenerator(gen token.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
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

// WithClientUnaryInterceptors for gRPC.
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

// WithUserAgent for gRPC.
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

// NewDialOptions for gRPC.
func NewDialOptions(opts ...ClientOption) ([]grpc.DialOption, error) {
	os := options(opts...)

	var security grpc.DialOption

	if tls.IsEnabled(os.security) {
		conf, err := tls.NewConfig(fs, os.security)
		if err != nil {
			return nil, err
		}

		creds := credentials.NewTLS(conf)

		security = grpc.WithTransportCredentials(creds)
	} else {
		security = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	cis := UnaryClientInterceptors(opts...)
	sto := streamDialOption(os)
	ops := []grpc.DialOption{
		grpc.WithUserAgent(os.userAgent.String()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                os.timeout,
			Timeout:             os.timeout,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(cis...), sto, security,
	}

	if os.compression {
		ops = append(ops, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	ops = append(ops, os.opts...)

	return ops, nil
}

// NewClient for gRPC.
func NewClient(target string, opts ...ClientOption) (*grpc.ClientConn, error) {
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
	unary = append(unary, ti.UnaryClientInterceptor(os.timeout))

	if os.retry != nil {
		timeout := time.MustParseDuration(os.retry.Timeout)
		backoff := time.MustParseDuration(os.retry.Backoff)

		unary = append(unary,
			ri.UnaryClientInterceptor(
				ri.WithCodes(codes.Unavailable, codes.DataLoss),
				ri.WithMax(uint(os.retry.Attempts)),
				ri.WithBackoff(ri.BackoffLinear(backoff)),
				ri.WithPerRetryTimeout(timeout),
			),
		)
	}

	if os.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor())
	}

	if os.logger != nil {
		unary = append(unary, tl.UnaryClientInterceptor(os.logger))
	}

	if os.meter != nil {
		unary = append(unary, tm.NewClient(os.meter).UnaryInterceptor())
	}

	if os.tracer != nil {
		unary = append(unary, tt.UnaryClientInterceptor(os.tracer))
	}

	if os.gen != nil {
		unary = append(unary, tkn.UnaryClientInterceptor(os.gen))
	}

	unary = append(unary, meta.UnaryClientInterceptor(os.userAgent, os.generator))

	return unary
}

func streamDialOption(opts *clientOpts) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{}
	stream = append(stream, opts.stream...)

	if opts.logger != nil {
		stream = append(stream, tl.StreamClientInterceptor(opts.logger))
	}

	if opts.meter != nil {
		stream = append(stream, tm.NewClient(opts.meter).StreamInterceptor())
	}

	if opts.tracer != nil {
		stream = append(stream, tt.StreamClientInterceptor(opts.tracer))
	}

	if opts.gen != nil {
		stream = append(stream, tkn.StreamClientInterceptor(opts.gen))
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
