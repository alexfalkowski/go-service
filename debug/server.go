package debug

import (
	"cmp"
	"net/http"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	"github.com/alexfalkowski/go-service/v2/net/http/server"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"go.uber.org/fx"
)

// ServerParams for debug.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Mux        *ServeMux
	Config     *Config
	Logger     *logger.Logger
	FS         *os.FS
}

// NewServer for debug.
func NewServer(params ServerParams) (*Server, error) {
	if !IsEnabled(params.Config) {
		return nil, nil
	}

	timeout := time.MustParseDuration(params.Config.Timeout)
	svr := &http.Server{
		Handler:     params.Mux,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
	}

	cfg, err := conf(params.FS, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	serv, err := server.NewService("debug", svr, cfg, params.Logger, params.Shutdowner)
	if err != nil {
		return nil, prefix(err)
	}

	return &Server{serv}, nil
}

// Server for debug.
type Server struct {
	*server.Service
}

// GetServer returns the server, if defined.
func (s *Server) GetServer() *server.Service {
	if s == nil {
		return nil
	}

	return s.Service
}

func conf(fs *os.FS, cfg *Config) (*config.Config, error) {
	config := &config.Config{
		Address: cmp.Or(cfg.Address, ":6060"),
	}

	if !tls.IsEnabled(cfg.TLS) {
		return config, nil
	}

	t, err := tls.NewConfig(fs, cfg.TLS)
	config.TLS = t

	return config, err
}

func prefix(err error) error {
	return errors.Prefix("debug", err)
}
