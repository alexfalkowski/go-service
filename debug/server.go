package debug

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newServeMux(lc fx.Lifecycle, logger *zap.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              "localhost:6060",
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
	l.Debug("starting debug server")

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Debug("could not start debug server", zap.Error(err))
	}
}

func stop(ctx context.Context, l *zap.Logger, s *http.Server) {
	l.Debug("stopping debug server", zap.Error(s.Shutdown(ctx)))
}
