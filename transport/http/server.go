package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/limiter"
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/http/cors"
	hl "github.com/alexfalkowski/go-service/transport/http/limiter"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	logger "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v3 "github.com/ulule/limiter/v3"
	"github.com/urfave/negroni/v3"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewServeMux for HTTP.
func NewServeMux() *runtime.ServeMux {
	opts := []runtime.ServeMuxOption{
		runtime.WithIncomingHeaderMatcher(customMatcher),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}

	return runtime.NewServeMux(opts...)
}

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Mux        *runtime.ServeMux
	Config     *Config
	Logger     *zap.Logger
	Tracer     trace.Tracer
	Meter      metric.Meter
	Limiter    *v3.Limiter
	Key        limiter.KeyFunc
	Handlers   []negroni.Handler `optional:"true"`
}

// Server for HTTP.
type Server struct {
	*server.Server
}

// NewServer for HTTP.
func NewServer(params ServerParams) (*Server, error) {
	n := negroni.New()
	n.Use(meta.NewHandler(UserAgent(params.Config)))
	n.Use(tracer.NewHandler(params.Tracer))
	n.Use(logger.NewHandler(params.Logger))
	n.Use(metrics.NewHandler(params.Meter))

	for _, hd := range params.Handlers {
		n.Use(hd)
	}

	n.Use(cors.New())
	n.Use(hl.NewHandler(params.Limiter, params.Key))
	n.UseHandler(params.Mux)

	s := &http.Server{
		Handler:     n,
		ReadTimeout: time.Timeout, WriteTimeout: time.Timeout,
		IdleTimeout: time.Timeout, ReadHeaderTimeout: time.Timeout,
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

	c := &sh.Config{}

	c.Port = cfg.Port

	if !security.IsEnabled(cfg.Security) {
		return c, nil
	}

	tls, err := security.NewServerTLSConfig(cfg.Security)
	c.TLS = tls

	return c, err
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id", "Geolocation", "X-Forwarded-For":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
