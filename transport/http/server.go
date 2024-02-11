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
	n.Use(szap.NewHandler(params.Logger))
	n.Use(tracer.NewHandler(params.Tracer))
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
	}

	return server, nil
}

// Start the server.
func (s *Server) Start() error {
	l, err := s.listener(s.config.Port)
	if err != nil {
		return err
	}

	go s.start(l)

	return nil
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) error {
	message := "stopping http server"
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error(message, zap.Error(err))
	} else {
		s.logger.Info(message)
	}

	return err
}

func (s *Server) start(l net.Listener) {
	s.logger.Info("starting http server", zap.String("addr", l.Addr().String()))

	if err := s.serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.String("addr", l.Addr().String()), zap.Error(err)}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start http server", fields...)
	}
}

func (s *Server) serve(l net.Listener) error {
	if s.config.Security.IsEnabled() {
		return s.server.ServeTLS(l, s.config.Security.CertFile, s.config.Security.KeyFile)
	}

	return s.server.Serve(l)
}

func (s *Server) listener(port string) (net.Listener, error) {
	if port == "" {
		return nil, ErrInvalidPort
	}

	return net.Listen("tcp", ":"+port)
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id", "Geolocation", "X-Forwarded-For":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
