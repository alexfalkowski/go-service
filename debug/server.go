package debug

import (
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewServeMux for debug.
func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}

// ServerParams for debug.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Mux        *http.ServeMux
	Config     *Config
	Logger     *zap.Logger
}

// Server for debug.
type Server struct {
	*server.Server
}

// NewServer for debug.
func NewServer(params ServerParams) (*Server, error) {
	s := &http.Server{
		Handler:     params.Mux,
		ReadTimeout: time.Timeout, WriteTimeout: time.Timeout,
		IdleTimeout: time.Timeout, ReadHeaderTimeout: time.Timeout,
	}

	sv, err := sh.NewServer(s, config(params.Config))
	if err != nil {
		return nil, errors.Prefix("new debug server", err)
	}

	svr := server.NewServer("debug", sv, params.Logger, params.Shutdowner)

	return &Server{Server: svr}, nil
}

func config(cfg *Config) sh.Config {
	c := sh.Config{}

	if !IsEnabled(cfg) {
		return c
	}

	c.Enabled = true
	c.Port = cfg.Port

	if !security.IsEnabled(cfg.Security) {
		return c
	}

	c.Security.Enabled = true
	c.Security.CertFile = cfg.Security.CertFile
	c.Security.KeyFile = cfg.Security.KeyFile

	return c
}
