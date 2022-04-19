package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	szap "github.com/alexfalkowski/go-service/transport/http/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	hopentracing "github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
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

	handler = hopentracing.NewHandler(params.Tracer, handler)
	handler = szap.NewHandler(params.Logger, handler)
	handler = meta.NewHandler(handler)

	addr := fmt.Sprintf(":%s", params.Config.Port)
	server := &Server{Mux: mux, server: &http.Server{Addr: addr, Handler: handler}}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			go startServer(server.server, listener, params)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return stopServer(ctx, server.server, params)
		},
	})

	return server
}

func startServer(server *http.Server, listener net.Listener, params ServerParams) {
	params.Logger.Info("starting http server", zap.String("port", params.Config.Port))

	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.String("port", params.Config.Port), zap.Error(err)}

		if err := params.Shutdowner.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		params.Logger.Error("could not start http server", fields...)
	}
}

func stopServer(ctx context.Context, server *http.Server, params ServerParams) error {
	params.Logger.Info("stopping http server", zap.String("port", params.Config.Port))

	if err := server.Shutdown(ctx); err != nil {
		params.Logger.Error("could not stop http server", zap.String("port", params.Config.Port), zap.Error(err))

		return err
	}

	return nil
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id", "Authorization":
		return key, true
	case "User-Agent":
		return "ua", true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
