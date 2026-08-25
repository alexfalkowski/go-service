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
// When automemlimit detects an unlimited cgroup and GOMEMLIMIT is not already
// configured, it sets Go's runtime memory limit to math.MaxInt64. This replaces
// any programmatic limit applied before startup.
//
// RegisterMemLimit is best-effort: any returned values and errors are intentionally
// ignored so that failure to set GOMEMLIMIT does not prevent a service from starting.
// Call the upstream automemlimit API directly to observe errors or choose a
// provider and limit policy.
func RegisterMemLimit(logger *slog.Logger) {
	_, _ = memlimit.Set(memlimit.WithLogger(logger))
}
