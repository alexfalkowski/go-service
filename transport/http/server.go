package http

import (
	"net/http"
	"time"

	ct "github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	lm "github.com/alexfalkowski/go-service/limiter"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/server"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/http/cors"
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
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	n := negroni.New()
	n.Use(meta.NewHandler(params.UserAgent, params.Version))
	n.Use(tracer.NewHandler(params.Tracer))
	n.Use(logger.NewHandler(params.Logger))
	n.Use(metrics.NewHandler(params.Meter))

	for _, hd := range params.Handlers {
		n.Use(hd)
	}

	if params.Verifier != nil {
		n.Use(ht.NewHandler(params.Verifier))
	}

	if params.Limiter != nil {
		n.Use(hl.NewHandler(params.Limiter, params.Key))
	}

	n.Use(cors.New())
	n.UseHandler(gzhttp.GzipHandler(params.Mux))

	s := &http.Server{
		Handler:     h2c.NewHandler(n, &http2.Server{}),
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
	}

	c, err := config(params.Config)
	if err != nil {
		return nil, err
	}

	sv, err := sh.NewServer(s, c)
	if err != nil {
		return nil, errors.Prefix("new http server", err)
	}

	svr := server.NewServer("http", sv, params.Logger, params.Shutdowner)

	return &Server{Server: svr}, nil
}

//nolint:nilnil
func config(cfg *Config) (*sh.Config, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	c := &sh.Config{
		Address: cfg.GetAddress(":8080"),
	}

	if !ct.IsEnabled(cfg.TLS) {
		return c, nil
	}

	tls, err := ct.NewConfig(cfg.TLS)
	c.TLS = tls

	return c, err
}

func timeout(cfg *Config) time.Duration {
	if !IsEnabled(cfg) {
		return time.Minute
	}

	return st.MustParseDuration(cfg.Timeout)
}
