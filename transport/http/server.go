package http

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/http/cors"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/telemetry"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Tracer     telemetry.Tracer
	Metrics    *telemetry.ServerMetrics
}

// Server for HTTP.
type Server struct {
	Mux    *runtime.ServeMux
	server *http.Server
	params ServerParams
}

// NewServer for HTTP.
func NewServer(params ServerParams) *Server {
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(customMatcher))

	registerHandles(mux)

	var handler http.Handler = mux

	handler = cors.New().Handler(handler)
	handler = params.Metrics.Handler(handler)
	handler = telemetry.NewTracerHandler(params.Tracer, handler)
	handler = telemetry.NewLoggerHandler(telemetry.LoggerHandlerParams{Logger: params.Logger, Handler: handler})
	handler = meta.NewHandler(handler)

	server := &Server{
		Mux:    mux,
		server: &http.Server{Handler: handler, ReadHeaderTimeout: time.Timeout},
		params: params,
	}

	return server
}

// Start the server.
func (s *Server) Start(listener net.Listener) {
	s.params.Logger.Info("starting http server", zap.String("addr", listener.Addr().String()))

	if err := s.server.Serve(listener); err != nil && !s.ignoreError(err) {
		fields := []zapcore.Field{zap.String("addr", listener.Addr().String()), zap.Error(err)}

		if err := s.params.Shutdowner.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.params.Logger.Error("could not start http server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	s.params.Logger.Info("stopping http server", zap.Error(s.server.Shutdown(ctx)))
}

func (s *Server) ignoreError(err error) bool {
	return errors.Is(err, http.ErrServerClosed) || errors.Is(err, cmux.ErrListenerClosed) || errors.Is(err, cmux.ErrServerClosed)
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id", "Authorization", "Geolocation":
		return key, true
	case "User-Agent":
		return "ua", true
	case "X-Forwarded-For":
		return "forwarded-for", true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func registerHandles(mux *runtime.ServeMux) {
	ph := promhttp.Handler()

	mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		ph.ServeHTTP(w, r)
	})
}
