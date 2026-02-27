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

// ServerParams defines dependencies for constructing a gRPC transport `Server`.
//
// It is an Fx parameter struct (`di.In`) that collects the configuration and optional
// dependencies used to build and run the gRPC server.
//
// Optional fields:
//   - `Logger`: enables server-side RPC outcome logging when non-nil.
//   - `Limiter`: enables server-side unary rate limiting when non-nil.
//   - `Verifier`: enables server-side token verification when non-nil.
//   - `Unary`/`Stream`: allow callers to inject additional interceptors (and are optional in DI).
type ServerParams struct {
	di.In

	// Shutdowner is used by the underlying `*server.Service` to coordinate shutdown.
	Shutdowner di.Shutdowner

	// Config controls gRPC server enablement, address, timeouts, TLS, and low-level gRPC options.
	Config *Config

	// Logger enables gRPC server logging interceptors when non-nil.
	Logger *logger.Logger

	// UserAgent is the service user agent used by metadata interceptors.
	UserAgent env.UserAgent

	// Version is the service version reported via metadata interceptors.
	Version env.Version

	// UserID identifies the metadata key used for injecting authenticated subjects into context.
	UserID env.UserID

	// ID generates request IDs when one is not already present.
	ID id.Generator

	// Limiter enables server-side unary rate limiting when non-nil.
	Limiter *limiter.Server

	// Verifier enables server-side token verification when non-nil.
	Verifier token.Verifier

	// Unary are additional unary server interceptors to append after the standard chain.
	Unary []grpc.UnaryServerInterceptor `optional:"true"`

	// Stream are additional stream server interceptors to append after the standard chain.
	Stream []grpc.StreamServerInterceptor `optional:"true"`
}

// NewServer constructs a gRPC transport `Server` when the transport is enabled.
//
// If `params.Config` is disabled, it returns (nil, nil) so that downstream wiring can treat the server
// as not configured.
//
// The constructed server includes:
//   - OpenTelemetry stats handling (server-side RPC instrumentation).
//   - A unary interceptor chain that performs metadata extraction/injection, and optionally logging,
//     token verification, and rate limiting, followed by any user-provided interceptors.
//   - A stream interceptor chain that performs metadata extraction/injection, and optionally logging
//     and token verification, followed by any user-provided interceptors.
//
// TLS:
// If TLS is enabled in config, server credentials are built from `params.Config.TLS`. Certificate/key
// sources may be resolved via the package-registered filesystem (see `Register` in this package).
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

// Server wraps a configured gRPC server and its runnable service wrapper.
//
// The embedded `*server.Service` provides start/stop orchestration and integrates with the application's
// lifecycle, while the underlying `*grpc.Server` is used as the `grpc.ServiceRegistrar` for registering
// service implementations.
type Server struct {
	server *grpc.Server
	*server.Service
}

// ServiceRegistrar returns the underlying gRPC service registrar.
//
// This is primarily used for registering generated gRPC services against the server.
// It returns nil if s is nil (for example, when the transport is disabled).
func (s *Server) ServiceRegistrar() grpc.ServiceRegistrar {
	if s == nil {
		return nil
	}
	return s.server
}

// GetService returns the runnable service wrapper.
//
// It returns nil if s is nil (for example, when the transport is disabled).
// This method is commonly used by higher-level wiring to collect enabled server services for lifecycle
// registration (see `transport.NewServers` and `transport.Register`).
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
