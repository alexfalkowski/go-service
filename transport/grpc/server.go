package grpc

import (
	"context"

	"github.com/alexfalkowski/go-service/limiter"
	sg "github.com/alexfalkowski/go-service/net/grpc"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/time"
	gl "github.com/alexfalkowski/go-service/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	logger "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	v3 "github.com/ulule/limiter/v3"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	Limiter    *v3.Limiter
	Key        limiter.KeyFunc
	Unary      []grpc.UnaryServerInterceptor  `optional:"true"`
	Stream     []grpc.StreamServerInterceptor `optional:"true"`
}

// Server for gRPC.
type Server struct {
	server *grpc.Server
	srv    *server.Server
}

// NewServer for gRPC.
func NewServer(params ServerParams) (*Server, error) {
	opt, err := creds(params.Config)
	if err != nil {
		return nil, err
	}

	metrics := metrics.NewServer(params.Meter)

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
		opt,
	}

	s := grpc.NewServer(opts...)
	reflection.Register(s)

	sv, err := sg.NewServer(s, config(params.Config))
	if err != nil {
		return nil, err
	}

	svr := server.NewServer("grpc", sv, params.Logger, params.Shutdowner)

	return &Server{srv: svr, server: s}, nil
}

// Start the server.
func (s *Server) Start() {
	s.srv.Start()
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	s.srv.Stop(ctx)
}

// Server for gRPC.
func (s *Server) Server() *grpc.Server {
	return s.server
}

func unaryServerOption(params ServerParams, m *metrics.Server, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(UserAgent(params.Config)),
		tracer.UnaryServerInterceptor(params.Tracer),
		logger.UnaryServerInterceptor(params.Logger),
		m.UnaryInterceptor(),
		gl.UnaryServerInterceptor(params.Limiter, params.Key),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(defaultInterceptors...))
}

func streamServerOption(params ServerParams, m *metrics.Server, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	defaultInterceptors := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(UserAgent(params.Config)),
		tracer.StreamServerInterceptor(params.Tracer),
		logger.StreamServerInterceptor(params.Logger),
		m.StreamInterceptor(),
		gl.StreamServerInterceptor(params.Limiter, params.Key),
	}

	defaultInterceptors = append(defaultInterceptors, interceptors...)

	return grpc.StreamInterceptor(middleware.ChainStreamServer(defaultInterceptors...))
}

func creds(cfg *Config) (grpc.ServerOption, error) {
	if !IsEnabled(cfg) || !security.IsEnabled(cfg.Security) {
		return grpc.EmptyServerOption{}, nil
	}

	var creds credentials.TransportCredentials

	if cfg.Security.HasKeyPair() {
		conf, err := security.NewTLSConfig(cfg.Security)
		if err != nil {
			return grpc.EmptyServerOption{}, err
		}

		creds = credentials.NewTLS(conf)
	} else {
		creds = credentials.NewClientTLSFromCert(nil, "")
	}

	return grpc.Creds(creds), nil
}

func config(cfg *Config) sg.Config {
	c := sg.Config{}

	if !IsEnabled(cfg) {
		return c
	}

	c.Enabled = true
	c.Port = cfg.Port

	return c
}
