package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/alexfalkowski/go-service/os"
	szap "github.com/alexfalkowski/go-service/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// ServerParams for gRPC.
type ServerParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Tracer     opentracing.Tracer
	Version    version.Version
	Unary      []grpc.UnaryServerInterceptor
	Stream     []grpc.StreamServerInterceptor
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
func NewServer(params ServerParams) *grpc.Server {
	labels := prom.Labels{"name": os.ExecutableName(), "version": string(params.Version)}
	metrics := prometheus.NewServerMetrics(prometheus.WithConstLabels(labels))
	metrics.EnableHandlingTimeHistogram(prometheus.WithHistogramBuckets(prom.DefBuckets))

	opts := []grpc.ServerOption{unaryServerOption(params, metrics, params.Unary...), streamServerOption(params, metrics, params.Stream...)}
	server := grpc.NewServer(opts...)

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", params.Config.Port))
			if err != nil {
				return err
			}

			go startServer(server, listener, params)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			stopServer(server, params)

			return nil
		},
	})

	return server
}

func startServer(server *grpc.Server, listener net.Listener, params ServerParams) {
	params.Logger.Info("starting grpc server", zap.String("port", params.Config.Port))

	if err := server.Serve(listener); err != nil {
		fields := []zapcore.Field{zap.String("port", params.Config.Port), zap.Error(err)}

		if err := params.Shutdowner.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		params.Logger.Error("could not start grpc server", fields...)
	}
}

func stopServer(server *grpc.Server, params ServerParams) {
	params.Logger.Info("stopping grpc server", zap.String("port", params.Config.Port))

	server.GracefulStop()
}

func unaryServerOption(params ServerParams, metrics *prometheus.ServerMetrics, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(),
		tags.UnaryServerInterceptor(),
		szap.UnaryServerInterceptor(params.Logger, params.Version),
		metrics.UnaryServerInterceptor(),
		opentracing.UnaryServerInterceptor(params.Tracer),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(defaultInterceptors...))
}

func streamServerOption(params ServerParams, metrics *prometheus.ServerMetrics, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(),
		tags.StreamServerInterceptor(),
		szap.StreamServerInterceptor(params.Logger, params.Version),
		metrics.StreamServerInterceptor(),
		opentracing.StreamServerInterceptor(params.Tracer),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(middleware.ChainStreamServer(defaultInterceptors...))
}
