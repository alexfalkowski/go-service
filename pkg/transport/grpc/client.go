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

// ClientParams for gRPC.
type ClientParams struct {
	Logger *zap.Logger
	Unary  []grpc.UnaryClientInterceptor
	Stream []grpc.StreamClientInterceptor
}

// NewClient to host for gRPC.
func NewClient(context context.Context, host string, params *ClientParams, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, unaryDialOption(params.Logger, params.Unary...), streamDialOption(params.Logger, params.Stream...))

	return grpc.DialContext(context, host, opts...)
}

func unaryDialOption(logger *zap.Logger, interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	defaultInterceptors := []grpc.UnaryClientInterceptor{
		grpcRetry.UnaryClientInterceptor(
			grpcRetry.WithCodes(codes.Unavailable, codes.DataLoss),
			grpcRetry.WithMax(5), // nolint:gomnd
			grpcRetry.WithBackoff(grpcRetry.BackoffLinear(50*time.Millisecond)), // nolint:gomnd
		),
		breaker.UnaryClientInterceptor(),
		meta.UnaryClientInterceptor(),
		pkgZap.UnaryClientInterceptor(logger),
		opentracing.UnaryClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.WithChainUnaryInterceptor(defaultInterceptors...)
}

func streamDialOption(logger *zap.Logger, interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	defaultInterceptors := []grpc.StreamClientInterceptor{
		meta.StreamClientInterceptor(),
		pkgZap.StreamClientInterceptor(logger),
		opentracing.StreamClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.WithChainStreamInterceptor(defaultInterceptors...)
}
