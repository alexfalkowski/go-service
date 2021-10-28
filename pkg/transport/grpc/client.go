package grpc

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/pkg/transport/grpc/breaker"
	pkgZap "github.com/alexfalkowski/go-service/pkg/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc/trace/opentracing"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	backoffLinear = 50 * time.Millisecond
)

// ClientParams for gRPC.
type ClientParams struct {
	Host   string
	Config *Config
	Logger *zap.Logger
	Unary  []grpc.UnaryClientInterceptor
	Stream []grpc.StreamClientInterceptor
}

// NewClient to host for gRPC.
func NewClient(context context.Context, params *ClientParams, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, unaryDialOption(params), streamDialOption(params))

	return grpc.DialContext(context, params.Host, opts...)
}

func unaryDialOption(params *ClientParams) grpc.DialOption {
	defaultInterceptors := []grpc.UnaryClientInterceptor{
		grpcRetry.UnaryClientInterceptor(
			grpcRetry.WithCodes(codes.Unavailable, codes.DataLoss),
			grpcRetry.WithMax(params.Config.Retry.Attempts),
			grpcRetry.WithBackoff(grpcRetry.BackoffLinear(backoffLinear)),
			grpcRetry.WithPerRetryTimeout(time.Duration(params.Config.Retry.Timeout)*time.Second),
		),
		breaker.UnaryClientInterceptor(),
		meta.UnaryClientInterceptor(params.Config.UserAgent),
		pkgZap.UnaryClientInterceptor(params.Logger),
		opentracing.UnaryClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, params.Unary...)

	return grpc.WithChainUnaryInterceptor(defaultInterceptors...)
}

func streamDialOption(params *ClientParams) grpc.DialOption {
	defaultInterceptors := []grpc.StreamClientInterceptor{
		meta.StreamClientInterceptor(params.Config.UserAgent),
		pkgZap.StreamClientInterceptor(params.Logger),
		opentracing.StreamClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, params.Stream...)

	return grpc.WithChainStreamInterceptor(defaultInterceptors...)
}
