package http

import (
	"cmp"
	"net/http"

	ct "github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/limiter"
	sh "github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	hl "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/meta"
	tl "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/logger"
	tm "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	tt "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/tracer"
	ht "github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
	"github.com/urfave/negroni/v3"
	"go.uber.org/fx"
)

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Mux        *http.ServeMux
	Config     *Config
	Logger     *logger.Logger
	Tracer     *tracer.Tracer
	Meter      *metrics.Meter
	UserAgent  env.UserAgent
	Version    env.Version
	ID         id.Generator
	FS         *os.FS
	Limiter    *limiter.Limiter  `optional:"true"`
	Verifier   token.Verifier    `optional:"true"`
	Handlers   []negroni.Handler `optional:"true"`
}

// NewServer for HTTP.
func NewServer(params ServerParams) (*Server, error) {
	if !IsEnabled(params.Config) {
		return nil, nil
	}

	timeout := time.MustParseDuration(params.Config.Timeout)

	neg := negroni.New()
	neg.Use(meta.NewHandler(params.UserAgent, params.Version, params.ID))

	if params.Tracer != nil {
		neg.Use(tt.NewHandler(params.Tracer))
	}

	if params.Logger != nil {
		neg.Use(tl.NewHandler(params.Logger))
	}

	if params.Meter != nil {
		neg.Use(tm.NewHandler(params.Meter))
	}

	for _, hd := range params.Handlers {
		neg.Use(hd)
	}

	if params.Verifier != nil {
		neg.Use(ht.NewHandler(params.Verifier))
	}

	if params.Limiter != nil {
		neg.Use(hl.NewHandler(params.Limiter))
	}

	neg.UseHandler(gzhttp.GzipHandler(params.Mux))

	svr := &http.Server{
		Handler:     neg,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
		Protocols: sh.Protocols(),
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
		Service: server.NewService("http", serv, params.Logger, params.Shutdowner),
	}

	return server, nil
}

// Server for HTTP.
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
		Address: cmp.Or(cfg.Address, ":8080"),
	}

	if !ct.IsEnabled(cfg.TLS) {
		return config, nil
	}

	tls, err := ct.NewConfig(fs, cfg.TLS)
	config.TLS = tls

	return config, err
}

func prefix(err error) error {
	return errors.Prefix("http", err)
}
