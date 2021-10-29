package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewServer for HTTP.
func NewServer(lc fx.Lifecycle, s fx.Shutdowner, cfg *Config, logger *zap.Logger) *http.Server {
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(customMatcher))
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{Addr: addr, Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			go startServer(s, server, listener, cfg, logger)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return stopServer(ctx, server, cfg, logger)
		},
	})

	return server
}

func startServer(s fx.Shutdowner, server *http.Server, listener net.Listener, cfg *Config, logger *zap.Logger) {
	logger.Info("starting http server", zap.String("port", cfg.Port))

	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		fields := []zapcore.Field{zap.String("port", cfg.Port), zap.Error(err)}

		if err := s.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		logger.Error("could not start http server", fields...)
	}
}

func stopServer(ctx context.Context, server *http.Server, cfg *Config, logger *zap.Logger) error {
	logger.Info("stopping http server", zap.String("port", cfg.Port))

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("could not stop http server", zap.String("port", cfg.Port), zap.Error(err))

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
