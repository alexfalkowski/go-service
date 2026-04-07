package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/server"
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

// Server bundles the dependencies required to build HTTP, gRPC, and debug servers for tests.
type Server struct {
	Lifecycle       di.Lifecycle
	Meter           metrics.Meter
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

// Register constructs and registers the enabled servers with the lifecycle.
//
// HTTP, gRPC, and debug servers are created only when their corresponding
// Register* flag is set. The method also wires the shared tracer registration
// and attaches common middleware, token handling, and test service handlers.
func (s *Server) Register() error {
	RegisterTracer(s.Lifecycle, s.Tracer)

	sh := NewShutdowner()
	servers := []*server.Service{}

	if s.RegisterHTTP {
		params := th.ServerParams{
			Shutdowner: sh, Mux: s.Mux,
			Config:   s.TransportConfig.HTTP,
			Logger:   s.Logger,
			Limiter:  s.HTTPLimiter,
			Handlers: []negroni.Handler{&EmptyHandler{}},
			Verifier: s.Verifier, ID: s.Generator, UserID: UserID,
			UserAgent: UserAgent, Version: Version,
		}

		httpServer, err := th.NewServer(params)
		if err != nil {
			return err
		}

		s.HTTPServer = httpServer
		servers = append(servers, httpServer.GetService())
	}

	if s.RegisterGRPC {
		params := grpc.ServerParams{
			Shutdowner: sh,
			Config:     s.TransportConfig.GRPC,
			Logger:     s.Logger,
			Verifier:   s.Verifier, ID: s.Generator, UserID: UserID,
			UserAgent: UserAgent, Version: Version,
			Limiter: s.GRPCLimiter,
		}

		grpcServer, err := grpc.NewServer(params)
		if err != nil {
			return err
		}

		s.GRPCServer = grpcServer
		v1.RegisterGreeterServiceServer(grpcServer.ServiceRegistrar(), NewService())
		servers = append(servers, grpcServer.GetService())
	}

	if s.RegisterDebug {
		debugServer, err := NewDebugServer(s.DebugConfig, s.Logger)
		if err != nil {
			return err
		}

		s.DebugServer = debugServer
		servers = append(servers, debugServer.GetService())
	}

	server.Register(s.Lifecycle, servers)

	return nil
}

// EmptyHandler is a no-op Negroni middleware used in server tests.
type EmptyHandler struct{}

// ServeHTTP implements negroni.Handler and just calls the next handler.
func (*EmptyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

// ErrServer is a server.Server test double whose Serve method fails with ErrFailed.
type ErrServer struct{}

// Serve implements server.Server and returns ErrFailed.
func (s *ErrServer) Serve() error {
	return ErrFailed
}

// Shutdown implements server.Server and always succeeds.
func (s *ErrServer) Shutdown(_ context.Context) error {
	return nil
}

// String returns a stable identifier for logs and assertions.
func (s *ErrServer) String() string {
	return "test"
}

// NoopServer is a server.Server test double whose lifecycle methods always succeed.
type NoopServer struct{}

// Serve implements server.Server and always succeeds.
func (s *NoopServer) Serve() error {
	return nil
}

// Shutdown implements server.Server and always succeeds.
func (s *NoopServer) Shutdown(_ context.Context) error {
	return nil
}

// String returns a stable identifier for logs and assertions.
func (s *NoopServer) String() string {
	return "test"
}
