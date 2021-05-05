package grpc

import (
	"context"
	"time"

	pkgZap "github.com/alexfalkowski/go-service/pkg/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/grpc/meta"
	"github.com/alexfalkowski/go-service/pkg/grpc/trace/opentracing"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// UnaryDialOption for gRPC.
func UnaryDialOption(logger *zap.Logger, interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	defaultInterceptors := []grpc.UnaryClientInterceptor{
		grpcRetry.UnaryClientInterceptor(
			grpcRetry.WithCodes(codes.Unavailable, codes.DataLoss),
			grpcRetry.WithMax(5), // nolint:gomnd
			grpcRetry.WithBackoff(grpcRetry.BackoffLinear(50*time.Millisecond)), // nolint:gomnd
		),
		meta.UnaryClientInterceptor(),
		pkgZap.UnaryClientInterceptor(logger),
		opentracing.UnaryClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.WithChainUnaryInterceptor(defaultInterceptors...)
}

// StreamDialOption for gRPC.
func StreamDialOption(logger *zap.Logger, interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	defaultInterceptors := []grpc.StreamClientInterceptor{
		meta.StreamClientInterceptor(),
		pkgZap.StreamClientInterceptor(logger),
		opentracing.StreamClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.WithChainStreamInterceptor(defaultInterceptors...)
}

// NewClient to host for gRPC.
func NewClient(context context.Context, host string, logger *zap.Logger, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if len(opts) == 0 {
		opts = append(opts, grpc.WithInsecure(), UnaryDialOption(logger), StreamDialOption(logger))
	}

	return grpc.DialContext(context, host, opts...)
}
