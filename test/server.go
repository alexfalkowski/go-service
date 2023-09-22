package test

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexfalkowski/go-service/security/jwt/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics/prometheus"
	gtracer "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/telemetry/metrics/prometheus"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
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
	c, _ := meta.RegisteredClaims(ctx)
	if s.verifyAuth && c == nil {
		return nil, ErrInvalidToken
	}

	return &v1.SayHelloResponse{Message: fmt.Sprintf("Hello %s", req.GetName())}, nil
}

// SayStreamHello ...
func (s *Server) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	c, _ := meta.RegisteredClaims(stream.Context())
	if s.verifyAuth && c == nil {
		return ErrInvalidToken
	}

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&v1.SayStreamHelloResponse{Message: fmt.Sprintf("Hello %s", req.GetName())})
}

// NewHTTPServer for test.
func NewHTTPServer(lc fx.Lifecycle, logger *zap.Logger, cfg *tracer.Config, tcfg *transport.Config) *shttp.Server {
	tracer, _ := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: cfg, Version: Version})

	server := shttp.NewServer(shttp.ServerParams{
		Shutdowner: NewShutdowner(), Config: &tcfg.HTTP, Logger: logger,
		Tracer: tracer, Metrics: hprometheus.NewServerCollector(lc, Version),
	})

	return server
}

// NewGRPCServer for test.
func NewGRPCServer(
	lc fx.Lifecycle, logger *zap.Logger, cfg *tracer.Config, tcfg *transport.Config,
	verifyAuth bool, unary []grpc.UnaryServerInterceptor, stream []grpc.StreamServerInterceptor,
) *tgrpc.Server {
	tracer, _ := gtracer.NewTracer(gtracer.Params{Lifecycle: lc, Config: cfg, Version: Version})

	server := tgrpc.NewServer(tgrpc.ServerParams{
		Shutdowner: NewShutdowner(), Config: &tcfg.GRPC, Logger: logger,
		Tracer: tracer, Metrics: gprometheus.NewServerCollector(lc, Version),
		Unary: unary, Stream: stream,
	})

	v1.RegisterGreeterServiceServer(server.Server, NewServer(verifyAuth))

	return server
}

// RegisterTransport for test.
func RegisterTransport(lc fx.Lifecycle, cfg *transport.Config, gs *tgrpc.Server, hs *shttp.Server) {
	transport.Register(transport.RegisterParams{Lifecycle: lc, Shutdowner: NewShutdowner(), Config: cfg, HTTP: hs, GRPC: gs})
}
