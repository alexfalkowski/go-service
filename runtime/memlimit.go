package runtime

import (
	"log/slog"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// RegisterMemLimit configures Go's memory limit (GOMEMLIMIT) using automemlimit.
//
// In containerized environments, the Go runtime may not automatically infer an
// appropriate memory limit from cgroup constraints. This helper delegates to the
// upstream automemlimit library to set a Go memory limit based on the detected
// container/cgroup memory limit.
//
// The provided logger is passed through to automemlimit and may be used to emit
// diagnostic messages about detection and the chosen limit.
//
// RegisterMemLimit is best-effort: any returned values and errors are intentionally
// ignored so that failure to set GOMEMLIMIT does not prevent a service from starting.
// If you need to enforce that a limit is set or handle errors, call the upstream
// automemlimit API directly.
func RegisterMemLimit(logger *slog.Logger) {
	_, _ = memlimit.SetGoMemLimitWithOpts(memlimit.WithLogger(logger))
}
