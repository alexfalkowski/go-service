package config

import (
	"context"
	"time"

	stime "github.com/alexfalkowski/go-service/time"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RandomWaitTime that will be used to shutdown the system.
var RandomWaitTime = stime.RandomWaitTime()

// WatchParams for config.
type WatchParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Logger     *zap.Logger
}

// Watch the configuration. If it changes terminate the application.
func Watch(params WatchParams) error {
	watcher := NewWatcher(params)

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return watcher.Start()
		},
		OnStop: func(ctx context.Context) error {
			return watcher.Stop()
		},
	})

	return nil
}

// NewWatcher of config changes.
func NewWatcher(params WatchParams) *Watcher {
	return &Watcher{sh: params.Shutdowner, logger: params.Logger}
}

// Watcher of config changes.
type Watcher struct {
	watcher *fsnotify.Watcher
	sh      fx.Shutdowner
	logger  *zap.Logger
}

// Start watching.
func (w *Watcher) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	w.watcher = watcher
	go w.watch()

	return watcher.Add(File())
}

// Stop watching.
func (w *Watcher) Stop() error {
	return w.watcher.Close()
}

func (w *Watcher) watch() {
	for {
		select {
		case e, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if e.Op&fsnotify.Write == fsnotify.Write {
				_ = w.shutdown()

				return
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			if err != nil {
				w.logger.Error("watching configuration", zap.Error(err))

				return
			}
		}
	}
}

func (w *Watcher) shutdown() error {
	w.logger.Info("configuration has been modified", zap.Duration("duration", RandomWaitTime))
	time.Sleep(RandomWaitTime)

	return w.sh.Shutdown()
}
