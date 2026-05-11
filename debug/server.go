package debug

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/config/server"
	debug "github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	httpserver "github.com/alexfalkowski/go-service/v2/net/http/server"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// ServerParams defines dependencies for constructing the debug server.
//
// It is intended for dependency injection (Fx/Dig). `NewServer` uses these dependencies to build an
// HTTP debug service and register lifecycle/shutdown behavior.
//
// Fields:
//   - Shutdowner: used by the underlying service wiring to coordinate process shutdown.
//   - Mux: the debug HTTP mux where debug endpoints (pprof/fgprof/statsviz/psutil/etc.) are registered.
//   - Config: enables and configures the debug server (address/timeout/TLS/options).
//   - Logger: used by the underlying HTTP service wrapper.
//   - FS: filesystem used to resolve TLS source strings when TLS is enabled.
type ServerParams struct {
	di.In
	Shutdowner di.Shutdowner
	Mux        *debug.ServeMux
	Config     *Config
	Logger     *logger.Logger
	FS         *os.FS
}

// NewServer constructs the debug Server when enabled.
//
// Disabled behavior: if params.Config is nil/disabled, NewServer returns (nil, nil).
//
// Enabled behavior:
//   - constructs an HTTP server using the configured timeout and debug mux,
//   - builds the net/http server config (address and optional TLS), and
//   - wraps it in a managed service ("debug") that integrates with DI lifecycle/shutdown.
//
// Failure behavior:
//   - TLS configuration errors (when TLS is enabled) are returned, and
//   - underlying service construction errors are returned.
//
// Errors are prefixed with "debug" for easier attribution.
// If the configured address uses an ephemeral port such as `localhost:0`, params.Config.Address is updated
// to the actual bound listener address after construction.
func NewServer(params ServerParams) (*Server, error) {
	if !params.Config.IsEnabled() {
		return nil, nil
	}

	httpServer := http.NewServer(params.Config.Options, params.Config.Timeout, params.Mux)

	cfg, err := newConfig(params.FS, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	service, err := httpserver.NewService("debug", httpServer, cfg, params.Logger, params.Shutdowner)
	if err != nil {
		return nil, prefix(err)
	}

	return &Server{service}, nil
}

// Server wraps the managed debug HTTP service.
//
// The embedded *httpserver.Service provides lifecycle integration and start/stop behavior.
// This wrapper adds a nil-safe accessor via GetService.
type Server struct {
	*httpserver.Service
}

// GetService returns the underlying service.
//
// It is nil-safe: if the receiver is nil (e.g. debug server disabled and not constructed), GetService
// returns nil.
func (s *Server) GetService() *httpserver.Service {
	if s == nil {
		return nil
	}
	return s.Service
}

func newConfig(fs *os.FS, cfg *Config) (*config.Config, error) {
	config := &config.Config{
		Address: cmp.Or(cfg.Address, net.DefaultAddress("6060")),
	}
	if !cfg.TLS.IsEnabled() {
		return config, nil
	}

	tlsConfig, err := server.NewConfig(fs, cfg.TLS)
	if err != nil {
		return nil, prefix(err)
	}

	config.TLS = tlsConfig
	return config, nil
}

func prefix(err error) error {
	return errors.Prefix("debug", err)
}
