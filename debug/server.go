package debug

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	debug "github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
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
	Mux        *debug.ServeMux
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
	svr := http.NewServer(timeout, params.Mux)

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

// GetService returns the service, if defined.
func (s *Server) GetService() *server.Service {
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
