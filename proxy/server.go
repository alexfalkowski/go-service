package proxy

import (
	"cmp"
	"net/http"
	"time"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/errors"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/server"
	t "github.com/alexfalkowski/go-service/time"
	proxy "github.com/elazarl/goproxy"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServerParams for proxy.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Server     *proxy.ProxyHttpServer
}

// Server for proxy.
type Server struct {
	*server.Server
}

// NewServer for proxy.
func NewServer(params ServerParams) (*Server, error) {
	timeout := timeout(params.Config)
	svr := &http.Server{
		Handler:     params.Server,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
	}

	c, err := config(params.Config)
	if err != nil {
		return nil, errors.Prefix("proxy", err)
	}

	serv, err := sh.NewServer(svr, c)
	if err != nil {
		return nil, errors.Prefix("proxy", err)
	}

	server := &Server{
		Server: server.NewServer("proxy", serv, params.Logger, params.Shutdowner),
	}

	return server, nil
}

//nolint:nilnil
func config(cfg *Config) (*sh.Config, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	config := &sh.Config{
		Address: cmp.Or(cfg.Address, ":7070"),
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
