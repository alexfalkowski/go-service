package debug

import (
	"cmp"
	"net/http"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/errors"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
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

	c, err := config(params.FS, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	serv, err := sh.NewServer(svr, c)
	if err != nil {
		return nil, prefix(err)
	}

	server := &Server{
		Service: server.NewService("debug", serv, params.Logger, params.Shutdowner),
	}

	return server, nil
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

func config(fs *os.FS, cfg *Config) (*sh.Config, error) {
	config := &sh.Config{
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
