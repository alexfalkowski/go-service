package http

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/http/cors"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	szap "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/urfave/negroni/v3"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
)

// ErrInvalidPort for HTTP.
var ErrInvalidPort = errors.New("invalid port")

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Tracer     tracer.Tracer
	Meter      metric.Meter
	Handlers   []negroni.Handler
}

// Server for HTTP.
type Server struct {
	Mux    *runtime.ServeMux
	server *http.Server
	sh     fx.Shutdowner
	config *Config
	logger *zap.Logger
	list   net.Listener
}

// ServerHandlers for HTTP.
func ServerHandlers() []negroni.Handler {
	return nil
}

// NewServer for HTTP.
func NewServer(params ServerParams) (*Server, error) {
	m, err := metrics.NewHandler(params.Meter)
	if err != nil {
		return nil, err
	}

	l, err := listener(params.Config)
	if err != nil {
		return nil, err
	}

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

	n := negroni.New()
	n.Use(meta.NewHandler())
	n.Use(tracer.NewHandler(params.Tracer))
	n.Use(szap.NewHandler(params.Logger))
	n.Use(m)

	for _, hd := range params.Handlers {
		n.Use(hd)
	}

	n.Use(cors.New())

	mux := runtime.NewServeMux(opts...)
	n.UseHandler(mux)

	s := &http.Server{
		Handler:           n,
		ReadTimeout:       time.Timeout,
		WriteTimeout:      time.Timeout,
		IdleTimeout:       time.Timeout,
		ReadHeaderTimeout: time.Timeout,
	}

	server := &Server{
		Mux:    mux,
		server: s,
		sh:     params.Shutdowner,
		config: params.Config,
		logger: params.Logger,
		list:   l,
	}

	return server, nil
}

// Start the server.
func (s *Server) Start() error {
	if s.list == nil {
		return nil
	}

	go s.start()

	return nil
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) error {
	if s.list == nil {
		return nil
	}

	message := "stopping http server"
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error(message, zap.Error(err), zap.String(tm.ServiceKey, "http"))
	} else {
		s.logger.Info(message, zap.String(tm.ServiceKey, "http"))
	}

	return err
}

func (s *Server) start() {
	s.logger.Info("starting server", zap.String("addr", s.list.Addr().String()), zap.String(tm.ServiceKey, "http"))

	if err := s.serve(s.list); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.String("addr", s.list.Addr().String()), zap.Error(err), zap.String(tm.ServiceKey, "http")}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start server", fields...)
	}
}

func (s *Server) serve(l net.Listener) error {
	se := s.config.Security
	if se != nil && se.Enabled {
		return s.server.ServeTLS(l, se.CertFile, se.KeyFile)
	}

	return s.server.Serve(l)
}

func listener(cfg *Config) (net.Listener, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}

	if cfg.Port == "" {
		return nil, ErrInvalidPort
	}

	return net.Listen("tcp", ":"+cfg.Port)
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id", "Geolocation", "X-Forwarded-For":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
