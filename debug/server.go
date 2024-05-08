package debug

import (
	"net/http"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/errors"
	sh "github.com/alexfalkowski/go-service/net/http"
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

	c, err := config(params.Config)
	if err != nil {
		return nil, err
	}

	sv, err := sh.NewServer(s, c)
	if err != nil {
		return nil, errors.Prefix("new debug server", err)
	}

	svr := server.NewServer("debug", sv, params.Logger, params.Shutdowner)

	return &Server{Server: svr}, nil
}

//nolint:nilnil
func config(cfg *Config) (*sh.Config, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	c := &sh.Config{}

	c.Port = cfg.Port

	if !tls.IsEnabled(cfg.TLS) {
		return c, nil
	}

	t, err := tls.NewConfig(cfg.TLS)
	c.TLS = t

	return c, err
}
