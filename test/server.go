package test

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexfalkowski/go-service/security/meta"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport"
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
func NewHTTPServer(lc fx.Lifecycle, logger *zap.Logger, cfg *opentracing.Config, tcfg *transport.Config) *shttp.Server {
	tracer, _ := hopentracing.NewTracer(hopentracing.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})

	server := shttp.NewServer(shttp.ServerParams{
		Shutdowner: NewShutdowner(), Config: &tcfg.HTTP, Logger: logger,
		Tracer: tracer, Metrics: hprometheus.NewServerMetrics(lc, Version),
	})

	return server
}

// NewGRPCServer for test.
func NewGRPCServer(
	lc fx.Lifecycle, logger *zap.Logger, cfg *opentracing.Config, tcfg *transport.Config,
	verifyAuth bool, unary []grpc.UnaryServerInterceptor, stream []grpc.StreamServerInterceptor,
) *tgrpc.Server {
	tracer, _ := gopentracing.NewTracer(gopentracing.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})

	server := tgrpc.NewServer(tgrpc.ServerParams{
		Shutdowner: NewShutdowner(), Config: &tcfg.GRPC, Logger: logger,
		Tracer: tracer, Metrics: gprometheus.NewServerMetrics(lc, Version),
		Unary: unary, Stream: stream,
	})

	v1.RegisterGreeterServiceServer(server.Server, NewServer(verifyAuth))

	return server
}

// RegisterTransport for test.
func RegisterTransport(lc fx.Lifecycle, cfg *transport.Config, gs *tgrpc.Server, hs *shttp.Server) {
	transport.Register(transport.RegisterParams{Lifecycle: lc, Shutdowner: NewShutdowner(), Config: cfg, HTTP: hs, GRPC: gs})
}
