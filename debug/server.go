package debug

import (
	"cmp"
	"net/http"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/errors"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
)

// ServerParams for debug.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *logger.Logger
}

// NewServer for debug.
func NewServer(params ServerParams) (*Server, error) {
	if !IsEnabled(params.Config) {
		return nil, nil
	}

	mux := http.NewServeMux()
	timeout := time.MustParseDuration(params.Config.Timeout)
	svr := &http.Server{
		Handler:     mux,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
	}

	c, err := config(params.Config)
	if err != nil {
		return nil, errors.Prefix("debug", err)
	}

	serv, err := sh.NewServer(svr, c)
	if err != nil {
		return nil, errors.Prefix("debug", err)
	}

	server := &Server{
		Server: server.NewServer("debug", serv, params.Logger, params.Shutdowner),
		mux:    mux,
	}

	return server, nil
}

// Server for debug.
type Server struct {
	mux *http.ServeMux
	*server.Server
}

// ServeMux for debug.
func (s *Server) ServeMux() *http.ServeMux {
	return s.mux
}

// GetServer returns the server, if defined.
func (s *Server) GetServer() *server.Server {
	if s == nil {
		return nil
	}

	return s.Server
}

func config(cfg *Config) (*sh.Config, error) {
	config := &sh.Config{
		Address: cmp.Or(cfg.Address, ":6060"),
	}

	if !tls.IsEnabled(cfg.TLS) {
		return config, nil
	}

	t, err := tls.NewConfig(cfg.TLS)
	config.TLS = t

	return config, err
}
