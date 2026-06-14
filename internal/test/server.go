package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	grpclimiter "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	httplimiter "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
)

// Server bundles the dependencies required to build HTTP, gRPC, and debug servers for tests.
type Server struct {
	// Lifecycle receives server start and stop hooks.
	Lifecycle di.Lifecycle
	// Meter provides metrics instrumentation for servers.
	Meter metrics.Meter
	// Verifier verifies inbound transport tokens.
	Verifier token.Verifier
	// Mux holds HTTP routes registered by tests.
	Mux *http.ServeMux
	// HTTPServer is populated when RegisterHTTP is true.
	HTTPServer *http.Server
	// GRPCServer is populated when RegisterGRPC is true.
	GRPCServer *grpc.Server
	// DebugServer is populated when RegisterDebug is true.
	DebugServer *debug.Server
	// TransportConfig configures HTTP and gRPC servers.
	TransportConfig *transport.Config
	// DebugConfig configures the debug server.
	DebugConfig *debug.Config
	// Tracer configures server tracing.
	Tracer *tracer.Config
	// GRPCLimiter limits inbound gRPC requests.
	GRPCLimiter *grpclimiter.Server
	// HTTPLimiter limits inbound HTTP requests.
	HTTPLimiter *httplimiter.Server
	// Logger logs server activity.
	Logger *logger.Logger
	// Generator generates request identifiers.
	Generator id.Generator
	// RegisterHTTP enables HTTP server construction.
	RegisterHTTP bool
	// RegisterGRPC enables gRPC server construction.
	RegisterGRPC bool
	// RegisterDebug enables debug server construction.
	RegisterDebug bool
}

// Register constructs and registers the enabled servers with the lifecycle.
//
// HTTP, gRPC, and debug servers are created only when their corresponding
// Register* flag is set. The method also wires the shared tracer registration
// and attaches common middleware, token handling, and test service handlers.
//
//nolint:funlen
func (s *Server) Register() error {
	RegisterTracer(s.Lifecycle, s.Tracer)

	sh := NewShutdowner()
	servers := []*server.Service{}

	if s.RegisterHTTP {
		params := http.ServerParams{
			Shutdowner: sh, Mux: s.Mux,
			Pool:     Pool,
			Config:   s.TransportConfig.HTTP,
			Logger:   s.Logger,
			Limiter:  s.HTTPLimiter,
			Handlers: []http.ChainedHandler{&EmptyHandler{}},
			Verifier: s.Verifier, ID: s.Generator, UserID: UserID,
			Name: Name, UserAgent: UserAgent, Version: Version,
		}

		httpServer, err := http.NewServer(params)
		if err != nil {
			return err
		}

		s.HTTPServer = httpServer
		s.TransportConfig.HTTP.Address = BoundAddress(s.TransportConfig.HTTP.Address, httpServer.GetService().String())
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
		s.TransportConfig.GRPC.Address = BoundAddress(s.TransportConfig.GRPC.Address, grpcServer.GetService().String())
		v1.RegisterGreeterServiceServer(grpcServer.ServiceRegistrar(), NewService())
		servers = append(servers, grpcServer.GetService())
	}

	if s.RegisterDebug {
		debugServer, err := NewDebugServer(s.Lifecycle, s.DebugConfig, s.Logger)
		if err != nil {
			return err
		}

		s.DebugServer = debugServer
		s.DebugConfig.Address = BoundAddress(s.DebugConfig.Address, debugServer.GetService().String())
		servers = append(servers, debugServer.GetService())
	}

	server.Register(s.Lifecycle, servers)

	return nil
}

// EmptyHandler is a no-op chained middleware used in server tests.
type EmptyHandler struct{}

// ServeHTTP implements [http.ChainedHandler] and just calls the next handler.
func (*EmptyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

// ErrServer is a [server.Server] test double whose Serve method fails with ErrFailed.
type ErrServer struct{}

// Serve implements [server.Server] and returns ErrFailed.
func (s *ErrServer) Serve() error {
	return ErrFailed
}

// Shutdown implements [server.Server] and always succeeds.
func (s *ErrServer) Shutdown(_ context.Context) error {
	return nil
}

// String returns a stable identifier for logs and assertions.
func (s *ErrServer) String() string {
	return "test"
}

// NoopServer is a [server.Server] test double whose lifecycle methods always succeed.
type NoopServer struct{}

// Serve implements [server.Server] and always succeeds.
func (s *NoopServer) Serve() error {
	return nil
}

// Shutdown implements [server.Server] and always succeeds.
func (s *NoopServer) Shutdown(_ context.Context) error {
	return nil
}

// String returns a stable identifier for logs and assertions.
func (s *NoopServer) String() string {
	return "test"
}

// NewObservableServer returns an ObservableServer with an initialized Done channel.
func NewObservableServer(err, shutdownErr error) *ObservableServer {
	return &ObservableServer{Err: err, ShutdownErr: shutdownErr, Done: make(chan struct{})}
}

// ObservableServer is a [server.Server] test double that records serve and shutdown activity.
type ObservableServer struct {
	Err         error
	ShutdownErr error
	Done        chan struct{}
	Shutdowns   int
}

// Serve closes Done and returns Err.
func (s *ObservableServer) Serve() error {
	close(s.Done)

	return s.Err
}

// Shutdown increments Shutdowns and returns ShutdownErr.
func (s *ObservableServer) Shutdown(_ context.Context) error {
	s.Shutdowns++

	return s.ShutdownErr
}

// String returns a stable identifier for logs and assertions.
func (*ObservableServer) String() string {
	return "test"
}
