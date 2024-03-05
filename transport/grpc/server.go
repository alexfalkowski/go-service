package grpc

import (
	"context"
	"errors"
	"net"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	szap "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
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

// ErrInvalidPort for gRPC.
var ErrInvalidPort = errors.New("invalid port")

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
		list:   l,
	}

	return server, nil
}

// Start the server.
func (s *Server) Start() error {
	if s.list == nil {
		return nil
	}

	go s.start()

	return nil
}

func (s *Server) start() {
	s.logger.Info("starting server", zap.String("addr", s.list.Addr().String()), zap.String(tm.ServiceKey, "grpc"))

	if err := s.Server.Serve(s.list); err != nil {
		fields := []zapcore.Field{zap.String("addr", s.list.Addr().String()), zap.Error(err), zap.String(tm.ServiceKey, "grpc")}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(_ context.Context) error {
	if s.list == nil {
		return nil
	}

	s.logger.Info("stopping grpc server", zap.String(tm.ServiceKey, "grpc"))

	s.Server.GracefulStop()

	return nil
}

func listener(cfg *Config) (net.Listener, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.Port == "" {
		return nil, ErrInvalidPort
	}

	return net.Listen("tcp", ":"+cfg.Port)
}

func unaryServerOption(l *zap.Logger, m *metrics.Server, t tracer.Tracer, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(),
		tracer.UnaryServerInterceptor(t),
		szap.UnaryServerInterceptor(l),
		m.UnaryInterceptor(),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(defaultInterceptors...))
}

func streamServerOption(l *zap.Logger, m *metrics.Server, t tracer.Tracer, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(),
		tracer.StreamServerInterceptor(t),
		szap.StreamServerInterceptor(l),
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
