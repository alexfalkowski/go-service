package grpc

import (
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	gm "github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	gt "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	r "github.com/grpc-ecosystem/go-grpc-middleware/retry"
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
type ClientOption interface{ apply(opts *clientOpts) }

type clientOpts struct {
	logger    *zap.Logger
	tracer    trace.Tracer
	meter     metric.Meter
	retry     *retry.Config
	breaker   bool
	userAgent string
	opts      []grpc.DialOption
	unary     []grpc.UnaryClientInterceptor
	stream    []grpc.StreamClientInterceptor
	security  grpc.DialOption
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) { f(o) }

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

// WithClientSecure for gRPC.
func WithClientSecure(sec *security.Config) (ClientOption, error) {
	var creds credentials.TransportCredentials

	if security.IsEnabled(sec) {
		conf, err := security.NewTLSConfig(sec)
		if err != nil {
			return nil, err
		}

		creds = credentials.NewTLS(conf)
	} else {
		creds = credentials.NewClientTLSFromCert(nil, "")
	}

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
func WithClientUserAgent(userAgent string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// NewDialOptions for gRPC.
func NewDialOptions(opts ...ClientOption) ([]grpc.DialOption, error) {
	cis, err := UnaryClientInterceptors(opts...)
	if err != nil {
		return nil, err
	}

	os := clientOptions(opts...)

	sto, err := streamDialOption(os)
	if err != nil {
		return nil, err
	}

	ops := []grpc.DialOption{
		grpc.WithUserAgent(os.userAgent),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Timeout,
			Timeout:             time.Timeout,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(cis...), sto, os.security,
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
func UnaryClientInterceptors(opts ...ClientOption) ([]grpc.UnaryClientInterceptor, error) {
	os := clientOptions(opts...)
	unary := []grpc.UnaryClientInterceptor{}

	unary = append(unary, os.unary...)

	if os.retry != nil {
		d := time.MustParseDuration(os.retry.Timeout)

		unary = append(unary,
			r.UnaryClientInterceptor(
				r.WithCodes(codes.Unavailable, codes.DataLoss),
				r.WithMax(os.retry.Attempts),
				r.WithBackoff(r.BackoffLinear(time.Backoff)),
				r.WithPerRetryTimeout(d),
			),
		)
	}

	if os.breaker {
		unary = append(unary, breaker.UnaryClientInterceptor())
	}

	if os.logger != nil {
		unary = append(unary, szap.UnaryClientInterceptor(os.logger))
	}

	if os.meter != nil {
		client, err := gm.NewClient(os.meter)
		if err != nil {
			return nil, err
		}

		unary = append(unary, client.UnaryInterceptor())
	}

	unary = append(unary, gt.UnaryClientInterceptor(os.tracer))
	unary = append(unary, meta.UnaryClientInterceptor(os.userAgent))

	return unary, nil
}

func streamDialOption(opts *clientOpts) (grpc.DialOption, error) {
	stream := []grpc.StreamClientInterceptor{}

	stream = append(stream, opts.stream...)

	if opts.logger != nil {
		stream = append(stream, szap.StreamClientInterceptor(opts.logger))
	}

	if opts.meter != nil {
		client, err := gm.NewClient(opts.meter)
		if err != nil {
			return nil, err
		}

		stream = append(stream, client.StreamInterceptor())
	}

	stream = append(stream, gt.StreamClientInterceptor(opts.tracer))
	stream = append(stream, meta.StreamClientInterceptor(opts.userAgent))

	return grpc.WithChainStreamInterceptor(stream...), nil
}

func clientOptions(opts ...ClientOption) *clientOpts {
	os := &clientOpts{
		security: grpc.WithTransportCredentials(insecure.NewCredentials()),
		tracer:   tracer.NewNoopTracer(),
	}
	for _, o := range opts {
		o.apply(os)
	}

	return os
}
