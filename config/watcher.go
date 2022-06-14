package config

import (
	"context"
	"time"

	stime "github.com/alexfalkowski/go-service/time"
	"github.com/radovskyb/watcher"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// WaitTime for shutting down.
type WaitTime time.Duration

// NewWaitTime for shutting down.
func NewWaitTime() WaitTime {
	return WaitTime(stime.RandomWaitTime())
}

// WatchParams for config.
type WatchParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Logger     *zap.Logger
	WaitTime   WaitTime
	Config     Configurator
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
	return &Watcher{sh: params.Shutdowner, logger: params.Logger, waitTime: time.Duration(params.WaitTime), cfg: params.Config}
}

// Watcher of config changes.
type Watcher struct {
	poll     *watcher.Watcher
	sh       fx.Shutdowner
	logger   *zap.Logger
	waitTime time.Duration
	cfg      Configurator
}

// Start watching.
func (w *Watcher) Start() error {
	poll := watcher.New()

	poll.FilterOps(watcher.Write)

	w.poll = poll
	go w.pollWatch()

	if err := poll.Add(File()); err != nil {
		return err
	}

	go poll.Start(time.Second)

	return nil
}

// Stop watching.
func (w *Watcher) Stop() error {
	w.poll.Close()

	return nil
}

func (w *Watcher) pollWatch() {
	for {
		select {
		case _ = <-w.poll.Event:
			_ = w.shutdown()
		case err := <-w.poll.Error:
			w.logger.Error("watching configuration", zap.Error(err))
		case <-w.poll.Closed:
			return
		}
	}
}

func (w *Watcher) shutdown() error {
	w.logger.Info("configuration has been modified", zap.Duration("duration", w.waitTime))
	time.Sleep(w.waitTime)

	return w.sh.Shutdown()
}
