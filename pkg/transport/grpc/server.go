package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgZap "github.com/alexfalkowski/go-service/pkg/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc/trace/opentracing"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// UnaryServerOption for gRPC.
func UnaryServerOption(logger *zap.Logger, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(),
		grpcRecovery.UnaryServerInterceptor(),
		grpcTags.UnaryServerInterceptor(),
		pkgZap.UnaryServerInterceptor(logger),
		grpcPrometheus.UnaryServerInterceptor,
		opentracing.UnaryServerInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(defaultInterceptors...))
}

// StreamServerOption for gRPC.
func StreamServerOption(logger *zap.Logger, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(),
		grpcRecovery.StreamServerInterceptor(),
		grpcTags.StreamServerInterceptor(),
		pkgZap.StreamServerInterceptor(logger),
		grpcPrometheus.StreamServerInterceptor,
		opentracing.StreamServerInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(defaultInterceptors...))
}

// NewServer for gRPC.
func NewServer(lc fx.Lifecycle, s fx.Shutdowner, cfg *config.Config, logger *zap.Logger, opts ...grpc.ServerOption) *grpc.Server {
	if len(opts) == 0 {
		opts = append(opts, UnaryServerOption(logger), StreamServerOption(logger))
	}

	server := grpc.NewServer(opts...)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
			if err != nil {
				return err
			}

			go startServer(s, server, listener, cfg, logger)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			stopServer(server, cfg, logger)

			return nil
		},
	})

	return server
}

func startServer(s fx.Shutdowner, server *grpc.Server, listener net.Listener, cfg *config.Config, logger *zap.Logger) {
	logger.Info("starting grpc server", zap.String("port", cfg.GRPCPort))

	if err := server.Serve(listener); err != nil {
		fields := []zapcore.Field{zap.String("port", cfg.GRPCPort), zap.Error(err)}

		if err := s.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		logger.Error("could not start grpc server", fields...)
	}
}

func stopServer(server *grpc.Server, cfg *config.Config, logger *zap.Logger) {
	logger.Info("stopping grpc server", zap.String("port", cfg.GRPCPort))

	server.GracefulStop()
}
