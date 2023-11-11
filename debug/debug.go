package debug

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/time"
	"github.com/arl/statsviz"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Register debug.
func Register(lc fx.Lifecycle, env env.Environment, logger *zap.Logger) {
	if !env.IsDevelopment() {
		return
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              "localhost:4200",
		Handler:           mux,
		ReadHeaderTimeout: time.Timeout,
	}

	statsviz.Register(mux)

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
}

func start(l *zap.Logger, s *http.Server) {
	l.Debug("starting debug server")

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Debug("could not start debug server", zap.Error(err))
	}
}

func stop(ctx context.Context, l *zap.Logger, s *http.Server) {
	l.Debug("stopping debug server", zap.Error(s.Shutdown(ctx)))
}
