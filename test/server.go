package test

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexfalkowski/go-service/security/meta"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	gopentracing "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	hopentracing "github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ErrInvalidToken ...
var ErrInvalidToken = errors.New("invalid token")

// NewServer ...
func NewServer(verifyAuth bool) *Server {
	return &Server{verifyAuth: verifyAuth}
}

// Server ...
type Server struct {
	verifyAuth bool
	v1.UnimplementedGreeterServiceServer
}

// SayHello ...
func (s *Server) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	if s.verifyAuth && meta.AuthorizedParty(ctx) == "" {
		return nil, ErrInvalidToken
	}

	return &v1.SayHelloResponse{Message: fmt.Sprintf("Hello %s", req.GetName())}, nil
}

// SayStreamHello ...
func (s *Server) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	if s.verifyAuth && meta.AuthorizedParty(stream.Context()) == "" {
		return ErrInvalidToken
	}

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&v1.SayStreamHelloResponse{Message: fmt.Sprintf("Hello %s", req.GetName())})
}

// NewHTTPServer for test.
func NewHTTPServer(lc fx.Lifecycle, logger *zap.Logger, cfg *opentracing.Config) (*shttp.Server, string) {
	tracer, _ := hopentracing.NewTracer(hopentracing.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})
	httpConfig := NewHTTPConfig()

	server := shttp.NewServer(shttp.ServerParams{
		Lifecycle: lc, Shutdowner: NewShutdowner(), Config: httpConfig, Logger: logger,
		Tracer: tracer, Version: Version, Metrics: hprometheus.NewServerMetrics(lc, Version),
	})

	return server, httpConfig.Port
}

// NewGRPCServer for test.
func NewGRPCServer(
	lc fx.Lifecycle, logger *zap.Logger, cfg *opentracing.Config, verifyAuth bool,
	unary []grpc.UnaryServerInterceptor, stream []grpc.StreamServerInterceptor,
) (*grpc.Server, *tgrpc.Config) {
	tracer, _ := gopentracing.NewTracer(gopentracing.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})
	grpcConfig := NewGRPCConfig()

	server := tgrpc.NewServer(tgrpc.ServerParams{
		Lifecycle: lc, Shutdowner: NewShutdowner(), Config: grpcConfig, Logger: logger,
		Tracer: tracer, Version: Version, Metrics: gprometheus.NewServerMetrics(lc, Version),
		Unary: unary, Stream: stream,
	})

	v1.RegisterGreeterServiceServer(server, NewServer(verifyAuth))

	return server, grpcConfig
}
