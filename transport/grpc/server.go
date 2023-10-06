package grpc

import (
	"context"
	"net"

	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// ServerParams for gRPC.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Tracer     tracer.Tracer
	Metrics    *prometheus.ServerCollector
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

// Server for gRPC.
type Server struct {
	Server *grpc.Server
	params ServerParams
}

// NewServer for gRPC.
func NewServer(params ServerParams) *Server {
	opts := []grpc.ServerOption{unaryServerOption(params, params.Unary...), streamServerOption(params, params.Stream...)}
	server := &Server{Server: grpc.NewServer(opts...), params: params}

	return server
}

// Start the server.
func (s *Server) Start(listener net.Listener) {
	if !s.params.Config.Enabled {
		listener.Close()

		return
	}

	s.params.Logger.Info("starting grpc server", zap.String("addr", listener.Addr().String()))

	if err := s.Server.Serve(listener); err != nil {
		fields := []zapcore.Field{zap.String("addr", listener.Addr().String()), zap.Error(err)}

		if err := s.params.Shutdowner.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.params.Logger.Error("could not start grpc server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(_ context.Context) {
	if !s.params.Config.Enabled {
		return
	}

	s.params.Logger.Info("stopping grpc server")

	s.Server.GracefulStop()
}

func unaryServerOption(params ServerParams, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(),
		tags.UnaryServerInterceptor(),
		szap.UnaryServerInterceptor(params.Logger),
		params.Metrics.UnaryServerInterceptor(),
		tracer.UnaryServerInterceptor(params.Tracer),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(defaultInterceptors...))
}

func streamServerOption(params ServerParams, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(),
		tags.StreamServerInterceptor(),
		szap.StreamServerInterceptor(params.Logger),
		params.Metrics.StreamServerInterceptor(),
		tracer.StreamServerInterceptor(params.Tracer),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(middleware.ChainStreamServer(defaultInterceptors...))
}
