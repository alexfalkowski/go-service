package server

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc/config"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"google.golang.org/grpc"
)

// Service is an alias for server.Service.
//
// A Service is a managed, runnable unit that participates in the application's
// lifecycle (start/stop) and integrates with go-service logging and shutdown
// wiring.
type Service = server.Service

// NewService constructs a managed go-service Service for a gRPC server.
//
// It:
//   - Binds the provided *grpc.Server to the listener described by cfg.Address
//     (via NewServer).
//   - Wraps the resulting Server in a generic go-service server runner
//     (server.NewService) that handles logging and coordinates shutdown via sh.
//
// Parameters:
//   - name: a logical service name used for logging/identification.
//   - grpc: the already-configured gRPC server instance (services registered,
//     interceptors/credentials/stats handlers configured, etc.).
//   - cfg: gRPC bind configuration; cfg.Address is expected to be in the
//     go-service network address format (for example "tcp://:9090").
//   - logger: the go-service logger used by the generic server runner.
//   - sh: a shutdown coordinator used to signal application shutdown.
//
// Returns an error if the listener cannot be created.
func NewService(name string, grpc *grpc.Server, cfg *config.Config, logger *logger.Logger, sh di.Shutdowner) (*Service, error) {
	serv, err := NewServer(grpc, cfg)
	if err != nil {
		return nil, err
	}

	return server.NewService(name, serv, logger, sh), nil
}

// NewServer constructs a Server that binds the provided *grpc.Server to cfg.Address.
//
// The address is expected to be in the go-service network address format (for
// example "tcp://:9090"). Internally, this is split into a network and address
// via net.SplitNetworkAddress and then bound using net.Listen.
//
// On listen error, the returned *Server is still non-nil (and contains the
// underlying *grpc.Server) to preserve existing behavior, but it will not have a
// listener configured.
func NewServer(server *grpc.Server, cfg *config.Config) (*Server, error) {
	srv := &Server{server: server}
	n, a, _ := net.SplitNetworkAddress(cfg.Address)

	l, err := net.Listen(context.Background(), n, a)
	if err != nil {
		return srv, err
	}

	srv.listener = l
	return srv, nil
}

// Server wraps a *grpc.Server together with the net.Listener it serves on.
//
// It adapts gRPC's Serve/GracefulStop API to the go-service generic server
// runner expectations (Serve, Shutdown, and String).
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// Serve starts serving requests on the configured listener.
//
// This delegates to (*grpc.Server).Serve. It will return an error if the server
// cannot begin serving, or if the underlying listener encounters an error.
func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

// Shutdown gracefully stops the gRPC server.
//
// This delegates to (*grpc.Server).GracefulStop, which stops accepting new
// connections and blocks until in-flight RPCs complete (subject to gRPC's
// internal semantics). The provided context is currently ignored.
func (s *Server) Shutdown(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

// String returns the bound listener address.
//
// This is typically used for logging/diagnostics.
func (s *Server) String() string {
	return s.listener.Addr().String()
}
