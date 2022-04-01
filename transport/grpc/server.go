package grpc

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	smeta "github.com/alexfalkowski/go-service/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServerParams for gRPC.
type ServerParams struct {
	fx.In

	Config *Config
	Logger *zap.Logger
	Unary  []grpc.UnaryServerInterceptor
	Stream []grpc.StreamServerInterceptor
}

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor() []grpc.UnaryServerInterceptor {
	return nil
}

// StreamServerInterceptor for gRPC.
func StreamServerInterceptor() []grpc.StreamServerInterceptor {
	return nil
}

// NewServer for gRPC.
func NewServer(lc fx.Lifecycle, s fx.Shutdowner, params ServerParams, opts ...grpc.ServerOption) *grpc.Server {
	opts = append(opts, unaryServerOption(params.Logger, params.Unary...), streamServerOption(params.Logger, params.Stream...))

	server := grpc.NewServer(opts...)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", params.Config.Port))
			if err != nil {
				return err
			}

			go startServer(s, server, listener, params)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			stopServer(server, params)

			return nil
		},
	})

	return server
}

func startServer(s fx.Shutdowner, server *grpc.Server, listener net.Listener, params ServerParams) {
	params.Logger.Info("starting grpc server", zap.String("port", params.Config.Port))

	if err := server.Serve(listener); err != nil {
		fields := []zapcore.Field{zap.String("port", params.Config.Port), zap.Error(err)}

		if err := s.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		params.Logger.Error("could not start grpc server", fields...)
	}
}

func stopServer(server *grpc.Server, params ServerParams) {
	params.Logger.Info("stopping grpc server", zap.String("port", params.Config.Port))

	server.GracefulStop()
}

// nolint:ireturn
func unaryServerOption(logger *zap.Logger, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(),
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(recoveryHandler)),
		tags.UnaryServerInterceptor(),
		szap.UnaryServerInterceptor(logger),
		prometheus.UnaryServerInterceptor,
		opentracing.UnaryServerInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(defaultInterceptors...))
}

// nolint:ireturn
func streamServerOption(logger *zap.Logger, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(),
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(recoveryHandler)),
		tags.StreamServerInterceptor(),
		szap.StreamServerInterceptor(logger),
		prometheus.StreamServerInterceptor,
		opentracing.StreamServerInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(middleware.ChainStreamServer(defaultInterceptors...))
}

func recoveryHandler(ctx context.Context, p any) error {
	smeta.WithAttribute(ctx, "grpc.error", fmt.Sprintf("%v", p))
	smeta.WithAttribute(ctx, "grpc.stack", string(debug.Stack()))

	return status.Error(codes.Internal, "recovered from error")
}
