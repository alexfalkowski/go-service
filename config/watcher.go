package config

import (
	"context"

	"github.com/alexfalkowski/go-service/time"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/fx"
)

// WatchParams for config.
type WatchParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
}

// Watch the configuration. If it changes terminate the application.
func Watch(params WatchParams) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go watch(params.Shutdowner, watcher)

			return watcher.Add(File())
		},
		OnStop: func(ctx context.Context) error {
			return watcher.Close()
		},
	})

	return nil
}

func watch(sh fx.Shutdowner, w *fsnotify.Watcher) {
	for {
		select {
		case e, ok := <-w.Events:
			if !ok {
				shutdown(sh)

				return
			}

			if e.Op&fsnotify.Write == fsnotify.Write {
				shutdown(sh)

				return
			}
		case _, ok := <-w.Errors:
			if !ok {
				shutdown(sh)

				return
			}
		}
	}
}

func shutdown(sh fx.Shutdowner) error {
	time.SleepRandomWaitTime()

	return sh.Shutdown()
}
