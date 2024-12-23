package grpc

import (
	"cmp"
	"context"
	"time"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	lm "github.com/alexfalkowski/go-service/limiter"
	sg "github.com/alexfalkowski/go-service/net/grpc"
	"github.com/alexfalkowski/go-service/server"
	t "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token"
	gl "github.com/alexfalkowski/go-service/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	logger "github.com/alexfalkowski/go-service/transport/grpc/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	tkn "github.com/alexfalkowski/go-service/transport/grpc/token"
	"github.com/sethvargo/go-limiter"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	Logger    *zap.Logger
	Tracer    trace.Tracer
	Meter     metric.Meter
	UserAgent env.UserAgent
	Version   env.Version
	Limiter   limiter.Store                  `optional:"true"`
	Key       lm.KeyFunc                     `optional:"true"`
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
		return nil, err
	}

	var ms *metrics.Server
	if params.Meter != nil {
		ms = metrics.NewServer(params.Meter)
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
		unaryServerOption(params, ms, params.Unary...),
		streamServerOption(params, ms, params.Stream...),
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
	uis := []grpc.UnaryServerInterceptor{meta.UnaryServerInterceptor(params.UserAgent, params.Version)}

	if params.Tracer != nil {
		uis = append(uis, tracer.UnaryServerInterceptor(params.Tracer))
	}

	if params.Logger != nil {
		uis = append(uis, logger.UnaryServerInterceptor(params.Logger))
	}

	if m != nil {
		uis = append(uis, m.UnaryInterceptor())
	}

	if params.Verifier != nil {
		uis = append(uis, tkn.UnaryServerInterceptor(params.Verifier))
	}

	if params.Limiter != nil {
		uis = append(uis, gl.UnaryServerInterceptor(params.Limiter, params.Key))
	}

	uis = append(uis, interceptors...)

	return grpc.ChainUnaryInterceptor(uis...)
}

func streamServerOption(params ServerParams, m *metrics.Server, interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	sis := []grpc.StreamServerInterceptor{meta.StreamServerInterceptor(params.UserAgent, params.Version)}

	if params.Tracer != nil {
		sis = append(sis, tracer.StreamServerInterceptor(params.Tracer))
	}

	if params.Logger != nil {
		sis = append(sis, logger.StreamServerInterceptor(params.Logger))
	}

	if m != nil {
		sis = append(sis, m.StreamInterceptor())
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

	return t.MustParseDuration(cfg.Timeout)
}
