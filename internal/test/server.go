package test

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/id"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	tg "github.com/alexfalkowski/go-service/v2/transport/grpc"
	th "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/urfave/negroni/v3"
	"go.uber.org/fx"
)

// Server for test.
type Server struct {
	Lifecycle       fx.Lifecycle
	Meter           *metrics.Meter
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
	RegisterHTTP    bool
	RegisterGRPC    bool
	RegisterDebug   bool
}

// Register server.
func (s *Server) Register() {
	sh := NewShutdowner()
	tracer := NewTracer(s.Lifecycle, s.Tracer)
	servers := []*server.Service{}

	if s.RegisterHTTP {
		httpServer, err := th.NewServer(th.ServerParams{
			Shutdowner: sh, Mux: s.Mux,
			Config: s.TransportConfig.HTTP, Logger: s.Logger,
			Tracer: tracer, Meter: s.Meter,
			Limiter: s.Limiter, Handlers: []negroni.Handler{&EmptyHandler{}},
			Verifier: s.Verifier, ID: s.ID,
			UserAgent: UserAgent, Version: Version,
			FS: FS,
		})
		runtime.Must(err)

		s.HTTPServer = httpServer
		servers = append(servers, httpServer.GetServer())
	}

	if s.RegisterGRPC {
		grpcServer, err := tg.NewServer(tg.ServerParams{
			Shutdowner: sh, Config: s.TransportConfig.GRPC, Logger: s.Logger,
			Tracer: tracer, Meter: s.Meter, Limiter: s.Limiter,
			Verifier: s.Verifier, ID: s.ID,
			UserAgent: UserAgent, Version: Version,
			FS: FS,
		})
		runtime.Must(err)

		s.GRPCServer = grpcServer
		v1.RegisterGreeterServiceServer(grpcServer.ServiceRegistrar(), NewService())
		servers = append(servers, grpcServer.GetServer())
	}

	if s.RegisterDebug {
		mux := debug.NewServeMux()
		debugServer, err := debug.NewServer(debug.ServerParams{
			Shutdowner: NewShutdowner(),
			Mux:        mux,
			Config:     s.DebugConfig,
			Logger:     s.Logger,
			FS:         FS,
		})
		runtime.Must(err)

		debug.RegisterPprof(mux)
		debug.RegisterFgprof(mux)
		debug.RegisterPsutil(mux, Content)

		err = debug.RegisterStatsviz(mux)
		runtime.Must(err)

		s.DebugServer = debugServer
		servers = append(servers, debugServer.GetServer())
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

func (s *NoopServer) Serve() error {
	return nil
}

func (s *NoopServer) Shutdown(_ context.Context) error {
	return nil
}

func (s *NoopServer) String() string {
	return "test"
}
