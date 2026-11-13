package runtime

import "github.com/alexfalkowski/go-signal/v2"

// Module wires runtime integrations into [github.com/alexfalkowski/go-signal/v2].
//
// Including this module in a go-signal application enables optional runtime tuning
// provided by this package.
//
// Currently, Module registers RegisterMemLimit, which attempts to configure
// Go's memory limit (GOMEMLIMIT) using the automemlimit library based on
// container/cgroup constraints.
// See [RegisterMemLimit] for its unlimited-cgroup behavior.
//
// Note: RegisterMemLimit is best-effort and intentionally ignores errors, so
// including this module will not fail application startup if a memory limit
// cannot be determined or applied.
var Module = signal.Module(
	signal.Register(RegisterMemLimit),
)
