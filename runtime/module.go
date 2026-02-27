package runtime

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires runtime integrations into Fx.
//
// Including this module in an Fx application enables optional runtime tuning
// provided by this package.
//
// Currently, Module registers RegisterMemLimit, which attempts to configure
// Go's memory limit (GOMEMLIMIT) using the automemlimit library based on
// container/cgroup constraints.
//
// Note: RegisterMemLimit is best-effort and intentionally ignores errors, so
// including this module will not fail application startup if a memory limit
// cannot be determined or applied.
var Module = di.Module(
	di.Register(RegisterMemLimit),
)
