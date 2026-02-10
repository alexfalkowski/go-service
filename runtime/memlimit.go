package runtime

import (
	"log/slog"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// RegisterMemLimit configures Go's memory limit using automemlimit.
//
// This function attempts to set GOMEMLIMIT automatically based on the container/cgroup memory limit.
// Any returned values and errors are intentionally ignored.
func RegisterMemLimit(logger *slog.Logger) {
	_, _ = memlimit.SetGoMemLimitWithOpts(memlimit.WithLogger(logger))
}
