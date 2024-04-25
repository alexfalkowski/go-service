package debug

import (
	"net"
	"net/http"

	sn "github.com/alexfalkowski/go-service/net"
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
	l, err := listener(params.Config)
	if err != nil {
		return nil, err
	}

	s := &http.Server{Handler: params.Mux, ReadTimeout: time.Timeout, WriteTimeout: time.Timeout, IdleTimeout: time.Timeout, ReadHeaderTimeout: time.Timeout}
	sv := sh.NewServer(s, config(params.Config, l))
	svr := server.NewServer("debug", sv, params.Logger, params.Shutdowner)

	return &Server{Server: svr}, nil
}

func config(cfg *Config, l net.Listener) sh.Config {
	c := sh.Config{
		Listener: l,
	}

	if !IsEnabled(cfg) || !security.IsEnabled(cfg.Security) {
		return c
	}

	c.Security.Enabled = true
	c.Security.CertFile = cfg.Security.CertFile
	c.Security.KeyFile = cfg.Security.KeyFile

	return c
}

func listener(cfg *Config) (net.Listener, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	return sn.Listener(cfg.Port)
}
