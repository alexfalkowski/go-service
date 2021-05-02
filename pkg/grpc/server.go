package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/alexfalkowski/go-service/pkg/config"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

const (
	server = "server"
)

// NewServerOptions is blank to satisfy the dep of []grpc.ServerOption
// This can be changed by implementers to provide options.
func NewServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

// NewServer for gRPC.
func NewServer(lc fx.Lifecycle, s fx.Shutdowner, cfg *config.Config, logger *zap.Logger, opts []grpc.ServerOption) *grpc.Server {
	allOpts := []grpc.ServerOption{
		unaryServerOption(logger),
		streamServerOption(logger),
	}
	allOpts = append(allOpts, opts...)

	server := grpc.NewServer(allOpts...)

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

func unaryServerOption(logger *zap.Logger) grpc.ServerOption {
	opt := grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
		metaUnaryServerInterceptor(),
		grpcRecovery.UnaryServerInterceptor(),
		grpcTags.UnaryServerInterceptor(),
		loggerUnaryServerInterceptor(logger),
		grpcPrometheus.UnaryServerInterceptor,
		traceUnaryServerInterceptor(),
	))

	return opt
}

func streamServerOption(logger *zap.Logger) grpc.ServerOption {
	opt := grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
		metaStreamServerInterceptor(),
		grpcRecovery.StreamServerInterceptor(),
		grpcTags.StreamServerInterceptor(),
		loggerStreamServerInterceptor(logger),
		grpcPrometheus.StreamServerInterceptor,
		traceStreamServerInterceptor(),
	))

	return opt
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
