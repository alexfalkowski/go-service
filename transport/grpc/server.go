package grpc

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/config"
	"github.com/alexfalkowski/go-service/v2/net/grpc/server"
	"github.com/alexfalkowski/go-service/v2/net/grpc/telemetry"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// ServerParams defines dependencies for constructing a gRPC Server.
type ServerParams struct {
	di.In
	Shutdowner di.Shutdowner
	Config     *Config
	Logger     *logger.Logger
	UserAgent  env.UserAgent
	Version    env.Version
	UserID     env.UserID
	ID         id.Generator
	Limiter    *limiter.Server
	Verifier   token.Verifier
	Unary      []grpc.UnaryServerInterceptor  `optional:"true"`
	Stream     []grpc.StreamServerInterceptor `optional:"true"`
}

// NewServer constructs a gRPC Server when the transport is enabled.
//
// If params.Config is disabled, it returns (nil, nil).
//
// The server is instrumented with OpenTelemetry stats handling and composes server-side interceptors for:
// metadata extraction, optional logging, optional token verification, optional rate limiting, plus any
// user-provided interceptors. If TLS is enabled, credentials are built using the registered filesystem
// (see Register in this package).
func NewServer(params ServerParams) (*Server, error) {
	if !params.Config.IsEnabled() {
		return nil, nil
	}

	opt, err := credsServerOption(fs, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	timeout := time.MustParseDuration(params.Config.Timeout)
	grpcServer := grpc.NewServer(params.Config.Options, timeout,
		grpc.StatsHandler(telemetry.NewServerHandler()),
		unaryServerOption(params, params.Unary...),
		streamServerOption(params, params.Stream...),
		opt,
	)
	cfg := &config.Config{Address: cmp.Or(params.Config.Address, net.DefaultAddress("9090"))}

	service, err := server.NewService("grpc", grpcServer, cfg, params.Logger, params.Shutdowner)
	if err != nil {
		return nil, prefix(err)
	}

	return &Server{server: grpcServer, Service: service}, nil
}

// Server wraps a gRPC server and its runnable service.
type Server struct {
	server *grpc.Server
	*server.Service
}

// ServiceRegistrar returns the underlying gRPC service registrar.
//
// It returns nil if s is nil.
func (s *Server) ServiceRegistrar() grpc.ServiceRegistrar {
	if s == nil {
		return nil
	}
	return s.server
}

// GetService returns the runnable service, if defined.
//
// It returns nil if s is nil.
func (s *Server) GetService() *server.Service {
	if s == nil {
		return nil
	}
	return s.Service
}

func unaryServerOption(params ServerParams, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	uis := []grpc.UnaryServerInterceptor{meta.UnaryServerInterceptor(params.UserAgent, params.Version, params.ID)}

	if params.Logger != nil {
		uis = append(uis, logger.UnaryServerInterceptor(params.Logger))
	}

	if params.Verifier != nil {
		uis = append(uis, token.UnaryServerInterceptor(params.UserID, params.Verifier))
	}

	if params.Limiter != nil {
		uis = append(uis, limiter.UnaryServerInterceptor(params.Limiter))
	}

	uis = append(uis, interceptors...)

	return grpc.ChainUnaryInterceptor(uis...)
}

func streamServerOption(params ServerParams, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	sis := []grpc.StreamServerInterceptor{meta.StreamServerInterceptor(params.UserAgent, params.Version, params.ID)}

	if params.Logger != nil {
		sis = append(sis, logger.StreamServerInterceptor(params.Logger))
	}

	if params.Verifier != nil {
		sis = append(sis, token.StreamServerInterceptor(params.UserID, params.Verifier))
	}

	sis = append(sis, interceptors...)

	return grpc.ChainStreamInterceptor(sis...)
}

func credsServerOption(fs *os.FS, cfg *Config) (grpc.ServerOption, error) {
	if !cfg.TLS.IsEnabled() {
		return grpc.EmptyServerOption{}, nil
	}

	conf, err := tls.NewConfig(fs, cfg.TLS)
	if err != nil {
		return grpc.EmptyServerOption{}, prefix(err)
	}

	return grpc.Creds(grpc.NewTLS(conf)), nil
}

func registrar(server *Server) grpc.ServiceRegistrar {
	return server.ServiceRegistrar()
}

func prefix(err error) error {
	return errors.Prefix("grpc", err)
}
