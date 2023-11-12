package debug

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/pprof"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/time"
	"github.com/arl/statsviz"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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
		Addr:              "localhost:6060",
		Handler:           mux,
		ReadHeaderTimeout: time.Timeout,
	}

	statsviz.Register(mux)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/debug/psutil", psutil)

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

func psutil(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := make(map[string]any)

	i, _ := cpu.InfoWithContext(ctx)
	t, _ := cpu.TimesWithContext(ctx, true)
	resp["cpu"] = map[string]any{
		"info":  i,
		"times": t,
	}

	s, _ := mem.SwapMemoryWithContext(ctx)
	v, _ := mem.VirtualMemoryWithContext(ctx)
	resp["mem"] = map[string]any{
		"swap":    s,
		"virtual": v,
	}

	w.Header().Add("Content-Type", "application/json")

	b, _ := json.Marshal(resp) //nolint:errchkjson
	w.Write(b)
}
