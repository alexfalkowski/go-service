package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/transport/http/cors"
	szap "github.com/alexfalkowski/go-service/transport/http/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Tracer     opentracing.Tracer
	Metrics    *prometheus.ServerMetrics
}

// Server for HTTP.
type Server struct {
	Mux    *runtime.ServeMux
	server *http.Server
}

// NewServer for HTTP.
func NewServer(params ServerParams) *Server {
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(customMatcher))

	var handler http.Handler = mux

	handler = cors.New().Handler(handler)
	handler = params.Metrics.Handler(handler)
	handler = opentracing.NewHandler(params.Tracer, handler)
	handler = szap.NewHandler(szap.HandlerParams{Logger: params.Logger, Handler: handler})
	handler = meta.NewHandler(handler)

	addr := fmt.Sprintf(":%s", params.Config.Port)
	server := &Server{Mux: mux, server: &http.Server{Addr: addr, Handler: handler}}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			go server.start(listener, params)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.stop(ctx, params)
		},
	})

	return server
}

func (s *Server) start(listener net.Listener, params ServerParams) {
	params.Logger.Info("starting http server", zap.String("port", params.Config.Port))

	if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.String("port", params.Config.Port), zap.Error(err)}

		if err := params.Shutdowner.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		params.Logger.Error("could not start http server", fields...)
	}
}

func (s *Server) stop(ctx context.Context, params ServerParams) error {
	params.Logger.Info("stopping http server", zap.String("port", params.Config.Port))

	if err := s.server.Shutdown(ctx); err != nil {
		params.Logger.Error("could not stop http server", zap.String("port", params.Config.Port), zap.Error(err))

		return err
	}

	return nil
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
