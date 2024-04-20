package http

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
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
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	Handlers   []negroni.Handler
}

// Server for HTTP.
type Server struct {
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

	n := negroni.New()
	n.Use(meta.NewHandler(UserAgent(params.Config)))
	n.Use(tracer.NewHandler(params.Tracer))
	n.Use(szap.NewHandler(params.Logger))
	n.Use(m)

	for _, hd := range params.Handlers {
		n.Use(hd)
	}

	n.Use(cors.New())
	n.UseHandler(params.Mux)

	s := &http.Server{
		Handler:           n,
		ReadTimeout:       time.Timeout,
		WriteTimeout:      time.Timeout,
		IdleTimeout:       time.Timeout,
		ReadHeaderTimeout: time.Timeout,
	}

	server := &Server{
		server: s,
		sh:     params.Shutdowner,
		config: params.Config,
		logger: params.Logger,
		list:   l,
	}

	return server, nil
}

// Start the server.
func (s *Server) Start() {
	if s.list == nil {
		return
	}

	go s.start()
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if s.list == nil {
		return
	}

	message := "stopping server"
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error(message, zap.Error(err), zap.String(tm.ServiceKey, "http"))
	} else {
		s.logger.Info(message, zap.String(tm.ServiceKey, "http"))
	}
}

func (s *Server) start() {
	s.logger.Info("starting server", zap.Stringer("addr", s.list.Addr()), zap.String(tm.ServiceKey, "http"))

	if err := s.serve(s.list); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.Stringer("addr", s.list.Addr()), zap.Error(err), zap.String(tm.ServiceKey, "http")}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start server", fields...)
	}
}

func (s *Server) serve(l net.Listener) error {
	if IsEnabled(s.config) && security.IsEnabled(s.config.Security) {
		return s.server.ServeTLS(l, s.config.Security.CertFile, s.config.Security.KeyFile)
	}

	return s.server.Serve(l)
}

func listener(cfg *Config) (net.Listener, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	return server.Listener(cfg.Port)
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id", "Geolocation", "X-Forwarded-For":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
