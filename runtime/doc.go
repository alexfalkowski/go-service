// Package runtime provides small runtime-oriented helpers used by go-service.
//
// The package is intentionally small and focused. It contains:
//   - strict helpers for startup/config paths (Must),
//   - panic recovery conversion helpers (ConvertRecover and ErrRecovered),
//   - build/runtime metadata helpers (Version),
//   - optional runtime tuning integration (RegisterMemLimit), and
//   - an Fx wiring module (Module) to register runtime integrations.
//
// # Error handling helpers
//
// Must is intended for code paths where failure is not meaningfully recoverable
// (for example, mandatory startup wiring). It panics when given a non-nil error
// and is commonly paired with constructors that already return an error.
//
// ConvertRecover is intended to be used with `recover()` in deferred functions.
// It converts an arbitrary recovered value into an error, wrapping it with
// ErrRecovered to provide a consistent sentinel that can be detected via errors.Is.
// The original panic value is preserved in the returned error string and, when it
// is already an error, is wrapped so it remains available via errors.As.
//
// Example:
//
//	func run() (err error) {
//		defer func() {
//			if v := recover(); v != nil {
//				err = runtime.ConvertRecover(v)
//			}
//		}()
//		// ...
//		return nil
//	}
//
// # Version reporting
//
// Version returns the build version of the running binary by reading build info
// via runtime/debug.ReadBuildInfo. When build info is unavailable, Version
// returns "development".
//
// Note: the exact value depends on how the binary is built. For example, builds
// produced from a module may embed a semantic version, a VCS-derived version, or
// "(devel)" depending on tooling and flags.
//
// # Memory limit integration (GOMEMLIMIT)
//
// RegisterMemLimit integrates with the upstream automemlimit library to set
// Go's memory limit (GOMEMLIMIT) based on container/cgroup constraints.
// This is typically useful in containerized deployments where the Go runtime
// otherwise cannot infer an appropriate heap limit.
//
// RegisterMemLimit is best-effort: it intentionally ignores returned values and
// errors, because failing to set a memory limit should not typically prevent a
// service from starting.
//
// # Dependency injection
//
// Module is an Fx module that registers RegisterMemLimit, allowing services that
// include the runtime module to apply the optional runtime tuning automatically.
//
// Start with Must, ConvertRecover, Version, and Module.
package runtime
