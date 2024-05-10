package test

import (
	"github.com/alexfalkowski/go-service/runtime"
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

// Mux for test.
var Mux = th.NewServeMux()

// Server for test.
type Server struct {
	Lifecycle  fx.Lifecycle
	Meter      metric.Meter
	Logger     *zap.Logger
	Tracer     *tracer.Config
	Transport  *transport.Config
	GRPC       *tg.Server
	HTTP       *th.Server
	Handlers   []negroni.Handler
	Unary      []grpc.UnaryServerInterceptor
	Stream     []grpc.StreamServerInterceptor
	VerifyAuth bool
}

// Register server.
func (s *Server) Register() {
	sh := NewShutdowner()
	tracer := tracer.NewTracer(s.Lifecycle, Environment, Version, s.Tracer, s.Logger)
	h, err := th.NewServer(th.ServerParams{
		Shutdowner: sh, Mux: Mux,
		Config: s.Transport.HTTP, Logger: s.Logger,
		Tracer: tracer, Meter: s.Meter, Handlers: s.Handlers,
	})
	runtime.Must(err)

	s.HTTP = h

	g, err := tg.NewServer(tg.ServerParams{
		Shutdowner: sh, Config: s.Transport.GRPC, Logger: s.Logger,
		Tracer: tracer, Meter: s.Meter,
		Unary: s.Unary, Stream: s.Stream,
	})
	runtime.Must(err)

	s.GRPC = g

	v1.RegisterGreeterServiceServer(g.Server(), NewService(s.VerifyAuth))
	transport.Register(transport.RegisterParams{Lifecycle: s.Lifecycle, Servers: []transport.Server{h, g}})
}
