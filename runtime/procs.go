package runtime

import (
	"log/slog"

	"go.uber.org/automaxprocs/maxprocs"
)

// RegisterMaxProcs for runtime.
func RegisterMaxProcs(logger *slog.Logger) {
	var opts []maxprocs.Option
	if logger != nil {
		opts = append(opts, maxprocs.Logger(logger.Info))
	}

	_, _ = maxprocs.Set(opts...)
}
