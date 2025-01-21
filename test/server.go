package test

import (
	"context"
	"net/http"
	"os"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	tg "github.com/alexfalkowski/go-service/transport/grpc"
	th "github.com/alexfalkowski/go-service/transport/http"
	"github.com/urfave/negroni/v3"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server for test.
type Server struct {
	Lifecycle       fx.Lifecycle
	Meter           metric.Meter
	Verifier        token.Verifier
	Mux             *http.ServeMux
	HTTPServer      *th.Server
	GRPCServer      *tg.Server
	DebugServer     *debug.Server
	TransportConfig *transport.Config
	DebugConfig     *debug.Config
	Tracer          *tracer.Config
	Limiter         *limiter.Limiter
	Logger          *zap.Logger
	ID              id.Generator
	VerifyAuth      bool
}

// Register server.
func (s *Server) Register() {
	sh := NewShutdowner()
	tracer := NewTracer(s.Lifecycle, s.Tracer, s.Logger)

	httpServer, err := th.NewServer(th.ServerParams{
		Shutdowner: sh, Mux: s.Mux,
		Config: s.TransportConfig.HTTP, Logger: s.Logger,
		Tracer: tracer, Meter: s.Meter,
		Limiter: s.Limiter, Handlers: []negroni.Handler{&none{}},
		Verifier: s.Verifier, ID: s.ID,
		UserAgent: UserAgent, Version: Version,
	})
	runtime.Must(err)

	s.HTTPServer = httpServer

	grpcServer, err := tg.NewServer(tg.ServerParams{
		Shutdowner: sh, Config: s.TransportConfig.GRPC, Logger: s.Logger,
		Tracer: tracer, Meter: s.Meter, Limiter: s.Limiter,
		Verifier: s.Verifier, ID: s.ID,
		UserAgent: UserAgent, Version: Version,
	})
	runtime.Must(err)

	s.GRPCServer = grpcServer

	debugServer, err := debug.NewServer(debug.ServerParams{
		Shutdowner: NewShutdowner(),
		Config:     s.DebugConfig,
		Logger:     s.Logger,
	})
	runtime.Must(err)

	debug.RegisterPprof(debugServer)
	debug.RegisterFgprof(debugServer)
	debug.RegisterPsutil(debugServer, Content)

	err = debug.RegisterStatsviz(debugServer)
	runtime.Must(err)

	s.DebugServer = debugServer

	v1.RegisterGreeterServiceServer(grpcServer.Server(), NewService(s.VerifyAuth))
	transport.Register(transport.RegisterParams{Lifecycle: s.Lifecycle, Servers: []transport.Server{httpServer, grpcServer, debugServer}})
}

type none struct{}

func (*none) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

// BadServer for test.
type BadServer struct{}

func (s *BadServer) IsEnabled() bool {
	return true
}

func (s *BadServer) Serve() error {
	return os.ErrNotExist
}

func (s *BadServer) Shutdown(_ context.Context) error {
	return nil
}

func (s *BadServer) String() string {
	return "test"
}

// NoopServer for test.
type NoopServer struct{}

func (s *NoopServer) IsEnabled() bool {
	return true
}

func (s *NoopServer) Serve() error {
	return nil
}

func (s *NoopServer) Shutdown(_ context.Context) error {
	return nil
}

func (s *NoopServer) String() string {
	return "test"
}
