package runtime

import (
	"log/slog"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// RegisterMemLimit for runtime.
func RegisterMemLimit(logger *slog.Logger) {
	_, _ = memlimit.SetGoMemLimitWithOpts(memlimit.WithLogger(logger))
}
