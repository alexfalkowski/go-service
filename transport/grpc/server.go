package grpc

import (
	"context"
	"net"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
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
	Tracer     trace.Tracer
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
	list   net.Listener
}

// NewServer for gRPC.
func NewServer(params ServerParams) (*Server, error) {
	metrics, err := metrics.NewServer(params.Meter)
	if err != nil {
		return nil, err
	}

	l, err := listener(params.Config)
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
		unaryServerOption(params, metrics, params.Unary...),
		streamServerOption(params, metrics, params.Stream...),
	}

	opt, err := creds(params.Config)
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
		list:   l,
	}

	return server, nil
}

// Start the server.
func (s *Server) Start() {
	if s.list == nil {
		return
	}

	go s.start()
}

func (s *Server) start() {
	s.logger.Info("starting server", zap.Stringer("addr", s.list.Addr()), zap.String(tm.ServiceKey, "grpc"))

	if err := s.Server.Serve(s.list); err != nil {
		fields := []zapcore.Field{zap.Stringer("addr", s.list.Addr()), zap.Error(err), zap.String(tm.ServiceKey, "grpc")}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(_ context.Context) {
	if s.list == nil {
		return
	}

	s.logger.Info("stopping server", zap.String(tm.ServiceKey, "grpc"))
	s.Server.GracefulStop()
}

func listener(cfg *Config) (net.Listener, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	return server.Listener(cfg.Port)
}

func unaryServerOption(params ServerParams, m *metrics.Server, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(UserAgent(params.Config)),
		tracer.UnaryServerInterceptor(params.Tracer),
		szap.UnaryServerInterceptor(params.Logger),
		m.UnaryInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(defaultInterceptors...))
}

func streamServerOption(params ServerParams, m *metrics.Server, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(UserAgent(params.Config)),
		tracer.StreamServerInterceptor(params.Tracer),
		szap.StreamServerInterceptor(params.Logger),
		m.StreamInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(middleware.ChainStreamServer(defaultInterceptors...))
}

func creds(cfg *Config) (grpc.ServerOption, error) {
	if !IsEnabled(cfg) || !security.IsEnabled(cfg.Security) {
		return nil, nil
	}

	conf, err := security.NewTLSConfig(cfg.Security)
	if err != nil {
		return nil, err
	}

	return grpc.Creds(credentials.NewTLS(conf)), nil
}
