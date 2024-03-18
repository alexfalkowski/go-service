package debug

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/time"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ErrInvalidPort for HTTP.
var ErrInvalidPort = errors.New("invalid port")

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
func (s *Server) Start() error {
	if s.list == nil {
		return nil
	}

	go s.start()

	return nil
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) error {
	if s.list == nil {
		return nil
	}

	message := "stopping server"
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error(message, zap.Error(err), zap.String(tm.ServiceKey, "debug"))
	} else {
		s.logger.Info(message, zap.String(tm.ServiceKey, "debug"))
	}

	return err
}

func (s *Server) start() {
	s.logger.Info("starting server", zap.String("addr", s.list.Addr().String()), zap.String(tm.ServiceKey, "debug"))

	if err := s.serve(s.list); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.String("addr", s.list.Addr().String()), zap.Error(err), zap.String(tm.ServiceKey, "debug")}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start server", fields...)
	}
}

func (s *Server) serve(l net.Listener) error {
	se := s.config.Security
	if security.IsEnabled(se) {
		return s.server.ServeTLS(l, se.CertFile, se.KeyFile)
	}

	return s.server.Serve(l)
}

func listener(cfg *Config) (net.Listener, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	if cfg.Port == "" {
		return nil, ErrInvalidPort
	}

	return net.Listen("tcp", ":"+cfg.Port)
}
