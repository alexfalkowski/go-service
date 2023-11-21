package grpc

import (
	"context"
	"net"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// ServerParams for gRPC.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Tracer     tracer.Tracer
	Meter      metric.Meter
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
	sh     fx.Shutdowner
	config *Config
	logger *zap.Logger
}

// NewServer for gRPC.
func NewServer(params ServerParams) (*Server, error) {
	metrics, err := metrics.NewServer(params.Meter)
	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Timeout,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Timeout,
			MaxConnectionAge:      time.Timeout,
			MaxConnectionAgeGrace: time.Timeout,
			Time:                  time.Timeout,
			Timeout:               time.Timeout,
		}),
		unaryServerOption(params.Logger, metrics, params.Tracer, params.Unary...),
		streamServerOption(params.Logger, metrics, params.Tracer, params.Stream...),
	}

	opt, err := creds(params)
	if err != nil {
		return nil, err
	}

	if opt != nil {
		opts = append(opts, opt)
	}

	s := grpc.NewServer(opts...)
	reflection.Register(s)

	server := &Server{
		Server: s,
		sh:     params.Shutdowner,
		config: params.Config,
		logger: params.Logger,
	}

	return server, nil
}

// Start the server.
func (s *Server) Start(listener net.Listener) {
	if !s.config.Enabled {
		listener.Close()

		return
	}

	s.logger.Info("starting grpc server", zap.String("addr", listener.Addr().String()))

	if err := s.Server.Serve(listener); err != nil {
		fields := []zapcore.Field{zap.String("addr", listener.Addr().String()), zap.Error(err)}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start grpc server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(_ context.Context) {
	if !s.config.Enabled {
		return
	}

	s.logger.Info("stopping grpc server")

	s.Server.GracefulStop()
}

func unaryServerOption(l *zap.Logger, m *metrics.Server, t tracer.Tracer, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(),
		szap.UnaryServerInterceptor(l),
		tracer.UnaryServerInterceptor(t),
		m.UnaryInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(defaultInterceptors...))
}

func streamServerOption(l *zap.Logger, m *metrics.Server, t tracer.Tracer, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(),
		szap.StreamServerInterceptor(l),
		tracer.StreamServerInterceptor(t),
		m.StreamInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(middleware.ChainStreamServer(defaultInterceptors...))
}

func creds(params ServerParams) (grpc.ServerOption, error) {
	if !params.Config.Security.IsEnabled() {
		return nil, nil
	}

	conf, err := security.NewTLSConfig(params.Config.Security)
	if err != nil {
		return nil, err
	}

	return grpc.Creds(credentials.NewTLS(conf)), nil
}
