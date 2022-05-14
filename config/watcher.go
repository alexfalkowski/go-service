package config

import (
	"context"
	"time"

	stime "github.com/alexfalkowski/go-service/time"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// WatchParams for config.
type WatchParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Logger     *zap.Logger
}

// Watch the configuration. If it changes terminate the application.
func Watch(params WatchParams) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go watch(params.Shutdowner, watcher, params.Logger)

			return watcher.Add(File())
		},
		OnStop: func(ctx context.Context) error {
			return watcher.Close()
		},
	})

	return nil
}

func watch(sh fx.Shutdowner, w *fsnotify.Watcher, logger *zap.Logger) {
	for {
		select {
		case e, ok := <-w.Events:
			if !ok {
				return
			}

			if e.Op&fsnotify.Write == fsnotify.Write {
				shutdown(sh, logger)

				return
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}

			if err != nil {
				logger.Error("watching configuration", zap.Error(err))

				return
			}
		}
	}
}

func shutdown(sh fx.Shutdowner, logger *zap.Logger) error {
	duration := stime.RandomWaitTime()

	logger.Info("configuration has been modified", zap.Duration("duration", duration))
	time.Sleep(duration)

	return sh.Shutdown()
}
