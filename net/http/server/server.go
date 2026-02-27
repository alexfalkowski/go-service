package server

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	"github.com/alexfalkowski/go-service/v2/net/http/errors"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// Service is an alias for server.Service.
type Service = server.Service

// NewService constructs a managed Service for an HTTP server.
//
// It adapts the provided *http.Server into this package's Server (see NewServer), then wraps it with the
// transport-agnostic lifecycle manager in net/server.
//
// The returned service will:
//   - start serving asynchronously (when Service.Start is called by upstream wiring),
//   - log start/stop events when a logger is provided, and
//   - trigger application shutdown via the provided di.Shutdowner if serving terminates unexpectedly.
//
// Errors returned are from listener creation performed by NewServer.
func NewService(name string, http *http.Server, cfg *config.Config, logger *logger.Logger, sh di.Shutdowner) (*Service, error) {
	serv, err := NewServer(http, cfg)
	if err != nil {
		return nil, err
	}

	return server.NewService(name, serv, logger, sh), nil
}

// NewServer constructs an HTTP Server adapter that owns a listener and serves HTTP (optionally with TLS).
//
// Address parsing:
// The cfg.Address is expected to use the go-service network address convention "<network>://<address>"
// (for example "tcp://:8080"). It is split using net.SplitNetworkAddress and then passed to net.Listen.
//
// Listener lifecycle:
// NewServer creates and stores the listener immediately. The listener is used by Serve/ServeTLS. If listener
// creation fails, the returned Server is still non-nil (with nil listener) alongside the error.
//
// TLS behavior:
// If cfg.TLS is non-nil, Serve will set the underlying http.Server.TLSConfig and call ServeTLS with empty
// cert/key file paths (because certificate material is expected to be present in TLSConfig). If cfg.TLS is nil,
// Serve will call Serve and serve plain HTTP.
func NewServer(server *http.Server, cfg *config.Config) (*Server, error) {
	srv := &Server{server: server, tls: cfg.TLS}
	n, a, _ := net.SplitNetworkAddress(cfg.Address)

	l, err := net.Listen(context.Background(), n, a)
	if err != nil {
		return srv, err
	}

	srv.listener = l
	return srv, nil
}

// Server adapts a *http.Server to the net/server.Server interface used by go-service.
//
// It holds:
//   - the underlying HTTP server,
//   - an optional TLS config, and
//   - the listener created during construction.
type Server struct {
	server   *http.Server
	tls      *tls.Config
	listener net.Listener
}

// Serve starts serving requests on the configured listener.
//
// Serve normalizes expected shutdown errors via net/http/errors.ServerError so callers can treat graceful
// termination consistently.
func (s *Server) Serve() error {
	return errors.ServerError(s.serve())
}

// Shutdown gracefully stops the underlying server.
//
// The provided context controls shutdown deadlines/cancellation.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// String returns the listener address as a string (used for logging/attribution).
func (s *Server) String() string {
	return s.listener.Addr().String()
}

func (s *Server) serve() error {
	if s.tls != nil {
		s.server.TLSConfig = s.tls

		return s.server.ServeTLS(s.listener, strings.Empty, strings.Empty)
	}

	return s.server.Serve(s.listener)
}
