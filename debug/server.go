package debug

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newServer(lc fx.Lifecycle, cfg *Config, json *marshaller.JSON, logger *zap.Logger) *server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:              "localhost:" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: time.Timeout,
	}

	s := &server{
		mux:    mux,
		srv:    srv,
		json:   json,
		logger: logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.stop(ctx)

			return nil
		},
	})

	return s
}

type server struct {
	mux    *http.ServeMux
	srv    *http.Server
	json   *marshaller.JSON
	logger *zap.Logger
}

func (s *server) start() {
	s.logger.Info("starting debug server")

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("could not start debug server", zap.Error(err))
	}
}

func (s *server) stop(ctx context.Context) {
	message := "stopping debug server"
	err := s.srv.Shutdown(ctx)

	if err != nil {
		s.logger.Error(message, zap.Error(err))
	} else {
		s.logger.Info(message)
	}
}
