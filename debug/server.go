package debug

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	debug "github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	"github.com/alexfalkowski/go-service/v2/net/http/server"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
)

// ServerParams for debug.
type ServerParams struct {
	di.In
	Shutdowner di.Shutdowner
	Mux        *debug.ServeMux
	Config     *Config
	Logger     *logger.Logger
	FS         *os.FS
}

// NewServer for debug.
func NewServer(params ServerParams) (*Server, error) {
	if !params.Config.IsEnabled() {
		return nil, nil
	}

	timeout := time.MustParseDuration(params.Config.Timeout)
	httpServer := http.NewServer(params.Config.Options, timeout, params.Mux)

	cfg, err := newConfig(params.FS, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	service, err := server.NewService("debug", httpServer, cfg, params.Logger, params.Shutdowner)
	if err != nil {
		return nil, prefix(err)
	}

	return &Server{service}, nil
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

func newConfig(fs *os.FS, cfg *Config) (*config.Config, error) {
	config := &config.Config{
		Address: cmp.Or(cfg.Address, net.DefaultAddress("6060")),
	}
	if !cfg.TLS.IsEnabled() {
		return config, nil
	}

	tls, err := tls.NewConfig(fs, cfg.TLS)
	if err != nil {
		return nil, prefix(err)
	}

	config.TLS = tls
	return config, nil
}

func prefix(err error) error {
	return errors.Prefix("debug", err)
}
