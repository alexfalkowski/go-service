package test

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/id"
	v1 "github.com/alexfalkowski/go-service/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	tg "github.com/alexfalkowski/go-service/transport/grpc"
	th "github.com/alexfalkowski/go-service/transport/http"
	"github.com/urfave/negroni/v3"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
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
	Logger          *logger.Logger
	ID              id.Generator
	VerifyAuth      bool
	RegisterHTTP    bool
	RegisterGRPC    bool
	RegisterDebug   bool
}

// Register server.
func (s *Server) Register() {
	sh := NewShutdowner()
	tracer := NewTracer(s.Lifecycle, s.Tracer, s.Logger)

	httpServer, err := th.NewServer(th.ServerParams{
		Shutdowner: sh, Mux: s.Mux,
		Config: s.TransportConfig.HTTP, Logger: s.Logger,
		Tracer: tracer, Meter: s.Meter,
		Limiter: s.Limiter, Handlers: []negroni.Handler{&EmptyHandler{}},
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

	servers := []transport.Server{}

	if s.RegisterHTTP {
		servers = append(servers, httpServer)
	}

	if s.RegisterGRPC {
		v1.RegisterGreeterServiceServer(grpcServer.Server(), NewService(s.VerifyAuth))
		servers = append(servers, grpcServer)
	}

	if s.RegisterDebug {
		servers = append(servers, debugServer)
	}

	transport.Register(s.Lifecycle, servers)
}

// EmptyHandler for test.
type EmptyHandler struct{}

func (*EmptyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

// ErrServer for test.
type ErrServer struct{}

func (s *ErrServer) IsEnabled() bool {
	return true
}

func (s *ErrServer) Serve() error {
	return ErrFailed
}

func (s *ErrServer) Shutdown(_ context.Context) error {
	return nil
}

func (s *ErrServer) String() string {
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
