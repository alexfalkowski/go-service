package config

import (
	"context"

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
				sh.Shutdown()

				return
			}

			if e.Op&fsnotify.Write == fsnotify.Write {
				sh.Shutdown()

				return
			}
		case _, ok := <-w.Errors:
			if !ok {
				sh.Shutdown()

				return
			}
		}
	}
}
