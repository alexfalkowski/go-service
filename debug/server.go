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
	*sn.Server
}

// NewServer for debug.
func NewServer(params ServerParams) (*Server, error) {
	l, err := listener(params.Config)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Handler:           params.Mux,
		ReadTimeout:       time.Timeout,
		WriteTimeout:      time.Timeout,
		IdleTimeout:       time.Timeout,
		ReadHeaderTimeout: time.Timeout,
	}

	c, k := files(params.Config)
	svr := sn.NewServer("debug", sh.NewServer(s, l, c, k), l, params.Logger, params.Shutdowner)

	return &Server{Server: svr}, nil
}

func files(cfg *Config) (string, string) {
	if IsEnabled(cfg) && security.IsEnabled(cfg.Security) {
		return cfg.Security.CertFile, cfg.Security.KeyFile
	}

	return "", ""
}

func listener(cfg *Config) (net.Listener, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	return server.Listener(cfg.Port)
}
