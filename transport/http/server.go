package http

import (
	"cmp"
	"net/http"
	"time"

	ct "github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	lm "github.com/alexfalkowski/go-service/limiter"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/server"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token"
	hl "github.com/alexfalkowski/go-service/transport/http/limiter"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	logger "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	ht "github.com/alexfalkowski/go-service/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
	"github.com/sethvargo/go-limiter"
	"github.com/urfave/negroni/v3"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Mux        *http.ServeMux
	Config     *Config
	Logger     *zap.Logger
	Tracer     trace.Tracer
	Meter      metric.Meter
	UserAgent  env.UserAgent
	Version    env.Version
	Limiter    limiter.Store     `optional:"true"`
	Key        lm.KeyFunc        `optional:"true"`
	Verifier   token.Verifier    `optional:"true"`
	Handlers   []negroni.Handler `optional:"true"`
}

// Server for HTTP.
type Server struct {
	*server.Server
}

// NewServer for HTTP.
func NewServer(params ServerParams) (*Server, error) {
	timeout := timeout(params.Config)

	neg := negroni.New()
	neg.Use(meta.NewHandler(params.UserAgent, params.Version))

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
		neg.Use(ht.NewHandler(params.Verifier))
	}

	if params.Limiter != nil {
		neg.Use(hl.NewHandler(params.Limiter, params.Key))
	}

	neg.UseHandler(gzhttp.GzipHandler(params.Mux))

	svr := &http.Server{
		Handler:     neg,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
	}

	c, err := config(params.Config)
	if err != nil {
		return nil, err
	}

	serv, err := sh.NewServer(svr, c)
	if err != nil {
		return nil, errors.Prefix("http", err)
	}

	server := &Server{
		Server: server.NewServer("http", serv, params.Logger, params.Shutdowner),
	}

	return server, nil
}

//nolint:nilnil
func config(cfg *Config) (*sh.Config, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	config := &sh.Config{
		Address: cmp.Or(cfg.Address, ":8080"),
	}

	if !ct.IsEnabled(cfg.TLS) {
		return config, nil
	}

	tls, err := ct.NewConfig(cfg.TLS)
	config.TLS = tls

	return config, err
}

func timeout(cfg *Config) time.Duration {
	if !IsEnabled(cfg) {
		return time.Minute
	}

	return st.MustParseDuration(cfg.Timeout)
}
