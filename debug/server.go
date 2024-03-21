package debug

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/time"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
}

// Server for HTTP.
type Server struct {
	Mux    *http.ServeMux
	server *http.Server
	sh     fx.Shutdowner
	config *Config
	logger *zap.Logger
	list   net.Listener
}

// NewServer for HTTP.
func NewServer(params ServerParams) (*Server, error) {
	l, err := listener(params.Config)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	s := &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Timeout,
		WriteTimeout:      time.Timeout,
		IdleTimeout:       time.Timeout,
		ReadHeaderTimeout: time.Timeout,
	}

	server := &Server{
		Mux:    mux,
		server: s,
		sh:     params.Shutdowner,
		config: params.Config,
		logger: params.Logger,
		list:   l,
	}

	return server, nil
}

// Start the server.
func (s *Server) Start() {
	if s.list == nil {
		return
	}

	go s.start()
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if s.list == nil {
		return
	}

	message := "stopping server"
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error(message, zap.Error(err), zap.String(tm.ServiceKey, "debug"))
	} else {
		s.logger.Info(message, zap.String(tm.ServiceKey, "debug"))
	}
}

func (s *Server) start() {
	s.logger.Info("starting server", zap.Stringer("addr", s.list.Addr()), zap.String(tm.ServiceKey, "debug"))

	if err := s.serve(s.list); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.Stringer("addr", s.list.Addr()), zap.Error(err), zap.String(tm.ServiceKey, "debug")}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start server", fields...)
	}
}

func (s *Server) serve(l net.Listener) error {
	if IsEnabled(s.config) && security.IsEnabled(s.config.Security) {
		return s.server.ServeTLS(l, s.config.Security.CertFile, s.config.Security.KeyFile)
	}

	return s.server.Serve(l)
}

func listener(cfg *Config) (net.Listener, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	return server.Listener(cfg.Port)
}
