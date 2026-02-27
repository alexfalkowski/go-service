package http

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	"github.com/alexfalkowski/go-service/v2/net/http/server"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/meta"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
	"github.com/urfave/negroni/v3"
)

// ServerParams defines dependencies for constructing an HTTP transport `Server`.
//
// It is an Fx parameter struct (`di.In`) that collects the configuration and optional dependencies used
// to build and run the HTTP server.
//
// Optional fields:
//   - `Logger`: enables server-side request logging middleware when non-nil.
//   - `Limiter`: enables server-side rate limiting middleware when non-nil.
//   - `Verifier`: enables server-side token verification middleware when non-nil.
//   - `Handlers`: allows callers to inject additional Negroni middleware (and is optional in DI).
type ServerParams struct {
	di.In

	// Shutdowner is used by the underlying `*server.Service` to coordinate shutdown.
	Shutdowner di.Shutdowner

	// Mux is the HTTP request multiplexer that holds registered routes/handlers.
	Mux *http.ServeMux

	// Config controls HTTP server enablement, address, timeouts, TLS, and low-level HTTP options.
	Config *Config

	// Logger enables HTTP server logging middleware when non-nil.
	Logger *logger.Logger

	// UserAgent is the service user agent used by metadata middleware.
	UserAgent env.UserAgent

	// Version is the service version reported via response headers and/or request context metadata.
	Version env.Version

	// UserID identifies the metadata key used for injecting authenticated subjects into context.
	UserID env.UserID

	// ID generates request IDs when one is not already present.
	ID id.Generator

	// Limiter enables server-side rate limiting middleware when non-nil.
	Limiter *limiter.Server

	// Verifier enables server-side token verification middleware when non-nil.
	Verifier token.Verifier

	// Handlers are additional Negroni handlers to insert into the middleware chain.
	//
	// These handlers are applied in the order provided.
	Handlers []negroni.Handler `optional:"true"`
}

// NewServer constructs an HTTP transport `Server` when the transport is enabled.
//
// If `params.Config` is disabled, it returns (nil, nil) so that downstream wiring can treat the server
// as not configured.
//
// Middleware composition:
//
// The server is built using Negroni and composes middleware in this order (first listed runs first):
//   - metadata extraction/injection and response headers (`transport/http/meta`)
//   - optional logging (`transport/http/telemetry/logger`) when `params.Logger` is non-nil
//   - optional user-provided handlers (`params.Handlers`, in the order supplied)
//   - optional token verification (`transport/http/token`) when `params.Verifier` is non-nil
//   - optional rate limiting (`transport/http/limiter`) when `params.Limiter` is non-nil
//   - gzip compression wrapping the mux handler (`gzhttp.GzipHandler(params.Mux)`)
//
// Token verification and rate limiting middleware typically treat "ignorable" paths (health/metrics/etc.)
// as bypassable, so those endpoints do not require auth and do not consume limiter capacity by default.
//
// TLS:
//
// If TLS is enabled, TLS configuration is constructed using the package-registered filesystem dependency
// (see `Register` in this package) to resolve certificate/key "source strings" (for example `file:` or `env:`).
func NewServer(params ServerParams) (*Server, error) {
	if !params.Config.IsEnabled() {
		return nil, nil
	}

	neg := negroni.New()
	neg.Use(meta.NewHandler(params.UserAgent, params.Version, params.ID))

	if params.Logger != nil {
		neg.Use(logger.NewHandler(params.Logger))
	}

	for _, hd := range params.Handlers {
		neg.Use(hd)
	}

	if params.Verifier != nil {
		neg.Use(token.NewHandler(params.UserID, params.Verifier))
	}

	if params.Limiter != nil {
		neg.Use(limiter.NewHandler(params.Limiter))
	}

	neg.UseHandler(gzhttp.GzipHandler(params.Mux))

	timeout := time.MustParseDuration(params.Config.Timeout)
	httpServer := http.NewServer(params.Config.Options, timeout, neg)

	cfg, err := newConfig(fs, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	service, err := server.NewService("http", httpServer, cfg, params.Logger, params.Shutdowner)
	if err != nil {
		return nil, prefix(err)
	}

	return &Server{service}, nil
}

// Server wraps an HTTP server service wrapper.
//
// The embedded `*server.Service` provides start/stop orchestration and integrates with the application's
// lifecycle.
type Server struct {
	*server.Service
}

// GetService returns the runnable service wrapper.
//
// It returns nil if s is nil (for example, when the transport is disabled).
// This method is commonly used by higher-level wiring to collect enabled server services for lifecycle
// registration (see `transport.NewServers` and `transport.Register`).
func (s *Server) GetService() *server.Service {
	if s == nil {
		return nil
	}
	return s.Service
}

func newConfig(fs *os.FS, cfg *Config) (*config.Config, error) {
	config := &config.Config{
		Address: cmp.Or(cfg.Address, net.DefaultAddress("8080")),
	}
	if !cfg.TLS.IsEnabled() {
		return config, nil
	}

	tls, err := tls.NewConfig(fs, cfg.TLS)
	if err != nil {
		return nil, prefix(err)
	}

	config.TLS = tls
	return config, nil
}

func prefix(err error) error {
	return errors.Prefix("http", err)
}
