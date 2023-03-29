package grpc

import (
	"context"
	"errors"
	"net"

	szap "github.com/alexfalkowski/go-service/transport/grpc/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/grpc/otel"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/soheilhy/cmux"
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
	Tracer     otel.Tracer
	Metrics    *prometheus.ServerMetrics
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
	s.params.Logger.Info("starting grpc server", zap.String("addr", listener.Addr().String()))

	if err := s.Server.Serve(listener); err != nil && !s.ignoreError(err) {
		fields := []zapcore.Field{zap.String("addr", listener.Addr().String()), zap.Error(err)}

		if err := s.params.Shutdowner.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.params.Logger.Error("could not start grpc server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(_ context.Context) {
	s.params.Logger.Info("stopping grpc server")

	s.Server.GracefulStop()
}

func (s *Server) ignoreError(err error) bool {
	return errors.Is(err, cmux.ErrListenerClosed) || errors.Is(err, cmux.ErrServerClosed)
}

func unaryServerOption(params ServerParams, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(),
		tags.UnaryServerInterceptor(),
		szap.UnaryServerInterceptor(params.Logger),
		params.Metrics.UnaryServerInterceptor(),
		otel.UnaryServerInterceptor(params.Tracer),
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
		otel.StreamServerInterceptor(params.Tracer),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(middleware.ChainStreamServer(defaultInterceptors...))
}
