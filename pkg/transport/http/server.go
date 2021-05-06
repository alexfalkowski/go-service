package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Register(lc fx.Lifecycle, s fx.Shutdowner, mux *runtime.ServeMux, cfg *config.Config, logger *zap.Logger) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.HTTPPort))
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
}

func startServer(s fx.Shutdowner, server *http.Server, listener net.Listener, cfg *config.Config, logger *zap.Logger) {
	logger.Info("starting http server", zap.String("port", cfg.HTTPPort))

	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		fields := []zapcore.Field{zap.String("port", cfg.HTTPPort), zap.Error(err)}

		if err := s.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		logger.Error("could not start http server", fields...)
	}
}

func stopServer(ctx context.Context, server *http.Server, cfg *config.Config, logger *zap.Logger) error {
	logger.Info("stopping http server", zap.String("port", cfg.HTTPPort))

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("could not stop http server", zap.String("port", cfg.HTTPPort), zap.Error(err))

		return err
	}

	return nil
}
