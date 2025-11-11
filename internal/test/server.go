package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	gl "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	th "github.com/alexfalkowski/go-service/v2/transport/http"
	hl "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/urfave/negroni/v3"
)

// Server for test.
type Server struct {
	Lifecycle       di.Lifecycle
	Meter           *metrics.Meter
	Verifier        token.Verifier
	Mux             *http.ServeMux
	HTTPServer      *th.Server
	GRPCServer      *grpc.Server
	DebugServer     *debug.Server
	TransportConfig *transport.Config
	DebugConfig     *debug.Config
	Tracer          *tracer.Config
	GRPCLimiter     *gl.Server
	HTTPLimiter     *hl.Server
	Logger          *logger.Logger
	Generator       id.Generator
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
		params := th.ServerParams{
			Shutdowner: sh, Mux: s.Mux,
			Config: s.TransportConfig.HTTP, Logger: s.Logger,
			Tracer: tracer, Meter: s.Meter, Limiter: s.HTTPLimiter,
			Handlers: []negroni.Handler{&EmptyHandler{}},
			Verifier: s.Verifier, ID: s.Generator, UserID: UserID,
			UserAgent: UserAgent, Version: Version,
			FS: FS,
		}

		httpServer, err := th.NewServer(params)
		runtime.Must(err)

		s.HTTPServer = httpServer
		servers = append(servers, httpServer.GetService())
	}

	if s.RegisterGRPC {
		params := grpc.ServerParams{
			Shutdowner: sh, Config: s.TransportConfig.GRPC,
			Logger: s.Logger, Tracer: tracer, Meter: s.Meter,
			Verifier: s.Verifier, ID: s.Generator, UserID: UserID,
			UserAgent: UserAgent, Version: Version,
			FS: FS, Limiter: s.GRPCLimiter,
		}

		grpcServer, err := grpc.NewServer(params)
		runtime.Must(err)

		s.GRPCServer = grpcServer
		v1.RegisterGreeterServiceServer(grpcServer.ServiceRegistrar(), NewService())
		servers = append(servers, grpcServer.GetService())
	}

	if s.RegisterDebug {
		debugServer := NewDebugServer(s.DebugConfig, s.Logger)

		s.DebugServer = debugServer
		servers = append(servers, debugServer.GetService())
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
