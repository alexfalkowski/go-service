package debug

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func mux(lc fx.Lifecycle, cfg *Config, logger *zap.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              "localhost:" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: time.Timeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go start(logger, server)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			stop(ctx, logger, server)

			return nil
		},
	})

	return mux
}

func start(l *zap.Logger, s *http.Server) {
	l.Info("starting debug server")

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Error("could not start debug server", zap.Error(err))
	}
}

func stop(ctx context.Context, l *zap.Logger, s *http.Server) {
	message := "stopping debug server"
	err := s.Shutdown(ctx)

	if err != nil {
		l.Error(message, zap.Error(err))
	} else {
		l.Info(message)
	}
}
