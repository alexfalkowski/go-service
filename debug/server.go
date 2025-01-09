package debug

import (
	"cmp"
	"net/http"
	"time"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/errors"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/server"
	t "github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServerParams for debug.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
}

// Server for debug.
type Server struct {
	mux *http.ServeMux
	*server.Server
}

// NewServer for debug.
func NewServer(params ServerParams) (*Server, error) {
	mux := http.NewServeMux()
	timeout := timeout(params.Config)
	svr := &http.Server{
		Handler:     mux,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
	}

	c, err := config(params.Config)
	if err != nil {
		return nil, err
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

// ServeMux for debug.
func (s *Server) ServeMux() *http.ServeMux {
	return s.mux
}

//nolint:nilnil
func config(cfg *Config) (*sh.Config, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

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

func timeout(cfg *Config) time.Duration {
	if !IsEnabled(cfg) {
		return time.Minute
	}

	return t.MustParseDuration(cfg.Timeout)
}
