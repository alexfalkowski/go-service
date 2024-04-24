package test

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/transport"
	tg "github.com/alexfalkowski/go-service/transport/grpc"
	th "github.com/alexfalkowski/go-service/transport/http"
	"github.com/urfave/negroni/v3"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Mux ...
var Mux = th.NewServeMux()

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
	if s.verifyAuth && !meta.IsEqual(Test(ctx), "auth") {
		return nil, ErrInvalidToken
	}

	return &v1.SayHelloResponse{Message: "Hello " + req.GetName()}, nil
}

// SayStreamHello ...
func (s *Server) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	if s.verifyAuth && !meta.IsEqual(Test(stream.Context()), "auth") {
		return ErrInvalidToken
	}

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&v1.SayStreamHelloResponse{Message: "Hello " + req.GetName()})
}

// NewHTTPServer for test.
func NewHTTPServer(lc fx.Lifecycle, logger *zap.Logger, cfg *tracer.Config, tcfg *transport.Config, meter metric.Meter, handlers []negroni.Handler) *th.Server {
	tracer, err := tracer.NewTracer(lc, Environment, Version, cfg, logger)
	if err != nil {
		panic(err)
	}

	server, err := th.NewServer(th.ServerParams{
		Shutdowner: NewShutdowner(), Mux: Mux,
		Config: tcfg.HTTP, Logger: logger,
		Tracer: tracer, Meter: meter, Handlers: handlers,
	})
	if err != nil {
		panic(err)
	}

	return server
}

// NewGRPCServer for test.
func NewGRPCServer(
	lc fx.Lifecycle, logger *zap.Logger, cfg *tracer.Config, tcfg *transport.Config,
	verifyAuth bool, meter metric.Meter,
	unary []grpc.UnaryServerInterceptor, stream []grpc.StreamServerInterceptor,
) *tg.Server {
	tracer, err := tracer.NewTracer(lc, Environment, Version, cfg, logger)
	if err != nil {
		panic(err)
	}

	server, err := tg.NewServer(tg.ServerParams{
		Shutdowner: NewShutdowner(), Config: tcfg.GRPC, Logger: logger,
		Tracer: tracer, Meter: meter,
		Unary: unary, Stream: stream,
	})
	if err != nil {
		panic(err)
	}

	v1.RegisterGreeterServiceServer(server.Server(), NewServer(verifyAuth))

	return server
}

// RegisterTransport for test.
func RegisterTransport(lc fx.Lifecycle, gs *tg.Server, hs *th.Server) {
	transport.Register(transport.RegisterParams{Lifecycle: lc, Servers: []transport.Server{gs, hs}})
}
