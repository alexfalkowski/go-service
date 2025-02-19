package grpc

import (
	"cmp"
	"context"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/limiter"
	sg "github.com/alexfalkowski/go-service/net/grpc"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token"
	gl "github.com/alexfalkowski/go-service/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	tl "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger"
	tm "github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	tt "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	tkn "github.com/alexfalkowski/go-service/transport/grpc/token"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// ServerParams for gRPC.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner

	Config    *Config
	Logger    *logger.Logger
	Tracer    *tracer.Tracer
	Meter     *metrics.Meter
	UserAgent env.UserAgent
	Version   env.Version
	ID        id.Generator
	Limiter   *limiter.Limiter               `optional:"true"`
	Verifier  token.Verifier                 `optional:"true"`
	Unary     []grpc.UnaryServerInterceptor  `optional:"true"`
	Stream    []grpc.StreamServerInterceptor `optional:"true"`
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
		return nil, errors.Prefix("grpc", err)
	}

	var meter *tm.Server
	if params.Meter != nil {
		meter = tm.NewServer(params.Meter)
	}

	timeout := timeout(params.Config)

	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             timeout,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     timeout,
			MaxConnectionAge:      timeout,
			MaxConnectionAgeGrace: timeout,
			Time:                  timeout,
			Timeout:               timeout,
		}),
		unaryServerOption(params, meter, params.Unary...),
		streamServerOption(params, meter, params.Stream...),
		opt,
	}

	svr := grpc.NewServer(opts...)
	reflection.Register(svr)

	serv, err := sg.NewServer(svr, config(params.Config))
	if err != nil {
		return nil, errors.Prefix("grpc", err)
	}

	server := &Server{
		srv:    server.NewServer("grpc", serv, params.Logger, params.Shutdowner),
		server: svr,
	}

	return server, nil
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

func unaryServerOption(params ServerParams, server *tm.Server, interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	uis := []grpc.UnaryServerInterceptor{meta.UnaryServerInterceptor(params.UserAgent, params.Version, params.ID)}

	if params.Tracer != nil {
		uis = append(uis, tt.UnaryServerInterceptor(params.Tracer))
	}

	if params.Logger != nil {
		uis = append(uis, tl.UnaryServerInterceptor(params.Logger))
	}

	if server != nil {
		uis = append(uis, server.UnaryInterceptor())
	}

	if params.Verifier != nil {
		uis = append(uis, tkn.UnaryServerInterceptor(params.Verifier))
	}

	if params.Limiter != nil {
		uis = append(uis, gl.UnaryServerInterceptor(params.Limiter))
	}

	uis = append(uis, interceptors...)

	return grpc.ChainUnaryInterceptor(uis...)
}

func streamServerOption(params ServerParams, server *tm.Server, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	sis := []grpc.StreamServerInterceptor{meta.StreamServerInterceptor(params.UserAgent, params.Version, params.ID)}

	if params.Tracer != nil {
		sis = append(sis, tt.StreamServerInterceptor(params.Tracer))
	}

	if params.Logger != nil {
		sis = append(sis, tl.StreamServerInterceptor(params.Logger))
	}

	if server != nil {
		sis = append(sis, server.StreamInterceptor())
	}

	if params.Verifier != nil {
		sis = append(sis, tkn.StreamServerInterceptor(params.Verifier))
	}

	sis = append(sis, interceptors...)

	return grpc.ChainStreamInterceptor(sis...)
}

func creds(cfg *Config) (grpc.ServerOption, error) {
	if !IsEnabled(cfg) || !tls.IsEnabled(cfg.TLS) {
		return grpc.EmptyServerOption{}, nil
	}

	conf, err := tls.NewConfig(cfg.TLS)
	if err != nil {
		return grpc.EmptyServerOption{}, err
	}

	creds := credentials.NewTLS(conf)

	return grpc.Creds(creds), nil
}

func config(cfg *Config) *sg.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	c := &sg.Config{
		Address: cmp.Or(cfg.Address, ":9090"),
	}

	return c
}

func timeout(cfg *Config) time.Duration {
	if !IsEnabled(cfg) {
		return time.Minute
	}

	return time.MustParseDuration(cfg.Timeout)
}
