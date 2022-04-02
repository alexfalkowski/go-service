package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/alexfalkowski/go-service/transport/grpc/breaker"
	szap "github.com/alexfalkowski/go-service/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
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

// NewLocalClient to localhost for gRPC.
func NewLocalClient(context context.Context, params *ClientParams) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("127.0.0.1:%s", params.Config.Port)
	unary := grpc.WithChainUnaryInterceptor(
		meta.UnaryClientInterceptor(params.Config.UserAgent),
		szap.UnaryClientInterceptor(params.Logger),
		opentracing.UnaryClientInterceptor(),
	)
	stream := grpc.WithChainStreamInterceptor(
		meta.StreamClientInterceptor(params.Config.UserAgent),
		szap.StreamClientInterceptor(params.Logger),
		opentracing.StreamClientInterceptor(),
	)

	return grpc.DialContext(context, target, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()), unary, stream)
}

// nolint:ireturn
func unaryDialOption(params *ClientParams) grpc.DialOption {
	defaultInterceptors := []grpc.UnaryClientInterceptor{
		retry.UnaryClientInterceptor(
			retry.WithCodes(codes.Unavailable, codes.DataLoss),
			retry.WithMax(params.Config.Retry.Attempts),
			retry.WithBackoff(retry.BackoffLinear(backoffLinear)),
			retry.WithPerRetryTimeout(params.Config.Retry.Timeout),
		),
		breaker.UnaryClientInterceptor(),
		meta.UnaryClientInterceptor(params.Config.UserAgent),
		szap.UnaryClientInterceptor(params.Logger),
		opentracing.UnaryClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, params.Unary...)

	return grpc.WithChainUnaryInterceptor(defaultInterceptors...)
}

// nolint:ireturn
func streamDialOption(params *ClientParams) grpc.DialOption {
	defaultInterceptors := []grpc.StreamClientInterceptor{
		meta.StreamClientInterceptor(params.Config.UserAgent),
		szap.StreamClientInterceptor(params.Logger),
		opentracing.StreamClientInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, params.Stream...)

	return grpc.WithChainStreamInterceptor(defaultInterceptors...)
}
