package grpc

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	rpc "github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/config"
	"github.com/alexfalkowski/go-service/v2/net/grpc/server"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor.
)

// ServerParams for gRPC.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *logger.Logger
	Tracer     *tracer.Tracer
	Meter      *metrics.Meter
	UserAgent  env.UserAgent
	Version    env.Version
	ID         id.Generator
	FS         *os.FS
	Limiter    *limiter.Limiter               `optional:"true"`
	Verifier   token.Verifier                 `optional:"true"`
	Unary      []grpc.UnaryServerInterceptor  `optional:"true"`
	Stream     []grpc.StreamServerInterceptor `optional:"true"`
}

// NewServer for gRPC.
func NewServer(params ServerParams) (*Server, error) {
	if !IsEnabled(params.Config) {
		return nil, nil
	}

	opt, err := creds(params.FS, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	var meter *metrics.Server
	if params.Meter != nil {
		meter = metrics.NewServer(params.Meter)
	}

	timeout := time.MustParseDuration(params.Config.Timeout)
	svr := rpc.NewServer(timeout,
		unaryServerOption(params, meter, params.Unary...),
		streamServerOption(params, meter, params.Stream...),
		opt,
	)
	cfg := &config.Config{Address: cmp.Or(params.Config.Address, ":9090")}

	serv, err := server.NewService("grpc", svr, cfg, params.Logger, params.Shutdowner)
	if err != nil {
		return nil, prefix(err)
	}

	return &Server{server: svr, Service: serv}, nil
}

// Server for gRPC.
type Server struct {
	server *grpc.Server
	*server.Service
}

// ServiceRegistrar for service registration.
func (s *Server) ServiceRegistrar() grpc.ServiceRegistrar {
	if s == nil {
		return nil
	}

	return s.server
}

// GetServer returns the server, if defined.
func (s *Server) GetServer() *server.Service {
	if s == nil {
		return nil
	}

	return s.Service
}

func unaryServerOption(params ServerParams, server *metrics.Server, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	uis := []grpc.UnaryServerInterceptor{meta.UnaryServerInterceptor(params.UserAgent, params.Version, params.ID)}

	if params.Tracer != nil {
		uis = append(uis, tracer.UnaryServerInterceptor(params.Tracer))
	}

	if params.Logger != nil {
		uis = append(uis, logger.UnaryServerInterceptor(params.Logger))
	}

	if server != nil {
		uis = append(uis, server.UnaryInterceptor())
	}

	if params.Verifier != nil {
		uis = append(uis, token.UnaryServerInterceptor(params.Verifier))
	}

	if params.Limiter != nil {
		uis = append(uis, limiter.UnaryServerInterceptor(params.Limiter))
	}

	uis = append(uis, interceptors...)

	return grpc.ChainUnaryInterceptor(uis...)
}

func streamServerOption(params ServerParams, server *metrics.Server, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	sis := []grpc.StreamServerInterceptor{meta.StreamServerInterceptor(params.UserAgent, params.Version, params.ID)}

	if params.Tracer != nil {
		sis = append(sis, tracer.StreamServerInterceptor(params.Tracer))
	}

	if params.Logger != nil {
		sis = append(sis, logger.StreamServerInterceptor(params.Logger))
	}

	if server != nil {
		sis = append(sis, server.StreamInterceptor())
	}

	if params.Verifier != nil {
		sis = append(sis, token.StreamServerInterceptor(params.Verifier))
	}

	sis = append(sis, interceptors...)

	return grpc.ChainStreamInterceptor(sis...)
}

func creds(fs *os.FS, cfg *Config) (grpc.ServerOption, error) {
	if !tls.IsEnabled(cfg.TLS) {
		return grpc.EmptyServerOption{}, nil
	}

	conf, err := tls.NewConfig(fs, cfg.TLS)
	if err != nil {
		return grpc.EmptyServerOption{}, err
	}

	creds := credentials.NewTLS(conf)

	return grpc.Creds(creds), nil
}

func provide(server *Server) grpc.ServiceRegistrar {
	return server.ServiceRegistrar()
}

func prefix(err error) error {
	return errors.Prefix("grpc", err)
}
