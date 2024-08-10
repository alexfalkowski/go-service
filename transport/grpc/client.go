package grpc

import (
	"time"

	st "github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security/token"
	t "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	logger "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	gm "github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	gt "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	tkn "github.com/alexfalkowski/go-service/transport/grpc/token"
	ri "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	ti "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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

var none = clientOptionFunc(func(_ *clientOpts) {
})

type clientOpts struct {
	tracer      trace.Tracer
	meter       metric.Meter
	security    grpc.DialOption
	gen         token.Generator
	logger      *zap.Logger
	retry       *retry.Config
	userAgent   env.UserAgent
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
		o.timeout = t.MustParseDuration(timeout)
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
func WithClientTLS(sec *st.Config) (ClientOption, error) {
	if !st.IsEnabled(sec) {
		return none, nil
	}

	conf, err := st.NewConfig(sec)
	if err != nil {
		return none, err
	}

	creds := credentials.NewTLS(conf)
	opt := clientOptionFunc(func(o *clientOpts) {
		o.security = grpc.WithTransportCredentials(creds)
	})

	return opt, nil
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
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// WithClientTracer for gRPC.
func WithClientTracer(tracer trace.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tracer = tracer
	})
}

// WithClientMetrics for gRPC.
func WithClientMetrics(meter metric.Meter) ClientOption {
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

// NewDialOptions for gRPC.
func NewDialOptions(opts ...ClientOption) []grpc.DialOption {
	cis := UnaryClientInterceptors(opts...)
	os := options(opts...)
	sto := streamDialOption(os)
	ops := []grpc.DialOption{
		grpc.WithUserAgent(string(os.userAgent)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                os.timeout,
			Timeout:             os.timeout,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(cis...), sto, os.security,
	}

	if os.compression {
		ops = append(ops, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	if os.gen != nil {
		ops = append(ops, grpc.WithPerRPCCredentials(tkn.NewPerRPCCredentials(os.gen)))
	}

	ops = append(ops, os.opts...)

	return ops
}

// NewClient for gRPC.
func NewClient(target string, opts ...ClientOption) (*grpc.ClientConn, error) {
	os := NewDialOptions(opts...)

	return grpc.NewClient(target, os...)
}

// UnaryClientInterceptors for gRPC.
func UnaryClientInterceptors(opts ...ClientOption) []grpc.UnaryClientInterceptor {
	os := options(opts...)
	unary := []grpc.UnaryClientInterceptor{}

	unary = append(unary, os.unary...)
	unary = append(unary, ti.UnaryClientInterceptor(os.timeout))

	if os.retry != nil {
		to := t.MustParseDuration(os.retry.Timeout)
		bo := t.MustParseDuration(os.retry.Backoff)

		unary = append(unary,
			ri.UnaryClientInterceptor(
				ri.WithCodes(codes.Unavailable, codes.DataLoss),
				ri.WithMax(os.retry.Attempts),
				ri.WithBackoff(ri.BackoffLinear(bo)),
				ri.WithPerRetryTimeout(to),
			),
		)
	}

	if os.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor())
	}

	if os.logger != nil {
		unary = append(unary, logger.UnaryClientInterceptor(os.logger))
	}

	if os.meter != nil {
		unary = append(unary, gm.NewClient(os.meter).UnaryInterceptor())
	}

	if os.tracer != nil {
		unary = append(unary, gt.UnaryClientInterceptor(os.tracer))
	}

	unary = append(unary, meta.UnaryClientInterceptor(os.userAgent))

	return unary
}

func streamDialOption(opts *clientOpts) grpc.DialOption {
	stream := []grpc.StreamClientInterceptor{}

	stream = append(stream, opts.stream...)

	if opts.logger != nil {
		stream = append(stream, logger.StreamClientInterceptor(opts.logger))
	}

	if opts.meter != nil {
		stream = append(stream, gm.NewClient(opts.meter).StreamInterceptor())
	}

	stream = append(stream, gt.StreamClientInterceptor(opts.tracer))
	stream = append(stream, meta.StreamClientInterceptor(opts.userAgent))

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

	if os.security == nil {
		os.security = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	return os
}
