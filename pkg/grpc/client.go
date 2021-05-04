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

// NewClient to host for gRPC.
func NewClient(context context.Context, host string, logger *zap.Logger, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	allOpts := []grpc.DialOption{
		unaryDialOption(logger),
		streamDialOption(logger),
	}
	allOpts = append(allOpts, opts...)

	return grpc.DialContext(context, host, allOpts...)
}

func unaryDialOption(logger *zap.Logger) grpc.DialOption {
	opt := grpc.WithChainUnaryInterceptor(
		grpcRetry.UnaryClientInterceptor(
			grpcRetry.WithCodes(codes.Unavailable, codes.DataLoss),
			grpcRetry.WithMax(5), // nolint:gomnd
			grpcRetry.WithBackoff(grpcRetry.BackoffLinear(50*time.Millisecond)), // nolint:gomnd
		),
		meta.UnaryClientInterceptor(),
		pkgZap.UnaryClientInterceptor(logger),
		opentracing.UnaryClientInterceptor(),
	)

	return opt
}

func streamDialOption(logger *zap.Logger) grpc.DialOption {
	opt := grpc.WithChainStreamInterceptor(
		meta.StreamClientInterceptor(),
		pkgZap.StreamClientInterceptor(logger),
		opentracing.StreamClientInterceptor(),
	)

	return opt
}
