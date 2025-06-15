package http

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	"github.com/alexfalkowski/go-service/v2/net/http/server"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/meta"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
	"github.com/urfave/negroni/v3"
)

// ServerParams for HTTP.
type ServerParams struct {
	di.In

	Shutdowner di.Shutdowner
	Mux        *http.ServeMux
	Config     *Config
	Logger     *logger.Logger
	Tracer     *tracer.Tracer
	Meter      *metrics.Meter
	UserAgent  env.UserAgent
	Version    env.Version
	UserID     env.UserID
	ID         id.Generator
	FS         *os.FS
	Limiter    *limiter.Limiter
	Verifier   token.Verifier
	Handlers   []negroni.Handler `optional:"true"`
}

// NewServer for HTTP.
func NewServer(params ServerParams) (*Server, error) {
	if !IsEnabled(params.Config) {
		return nil, nil
	}

	neg := negroni.New()
	neg.Use(meta.NewHandler(params.UserAgent, params.Version, params.ID))

	if params.Tracer != nil {
		neg.Use(tracer.NewHandler(params.Tracer))
	}

	if params.Logger != nil {
		neg.Use(logger.NewHandler(params.Logger))
	}

	if params.Meter != nil {
		neg.Use(metrics.NewHandler(params.Meter))
	}

	for _, hd := range params.Handlers {
		neg.Use(hd)
	}

	if params.Verifier != nil {
		neg.Use(token.NewHandler(params.UserID, params.Verifier))
	}

	if params.Limiter != nil {
		neg.Use(limiter.NewHandler(params.Limiter))
	}

	neg.UseHandler(gzhttp.GzipHandler(params.Mux))

	timeout := time.MustParseDuration(params.Config.Timeout)
	svr := http.NewServer(timeout, neg)

	cfg, err := newConfig(params.FS, params.Config)
	if err != nil {
		return nil, prefix(err)
	}

	serv, err := server.NewService("http", svr, cfg, params.Logger, params.Shutdowner)
	if err != nil {
		return nil, prefix(err)
	}

	return &Server{serv}, nil
}

// Server for HTTP.
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
		Address: cmp.Or(cfg.Address, ":8080"),
	}

	if !tls.IsEnabled(cfg.TLS) {
		return config, nil
	}

	tls, err := tls.NewConfig(fs, cfg.TLS)
	config.TLS = tls

	return config, err
}

func prefix(err error) error {
	return errors.Prefix("http", err)
}
