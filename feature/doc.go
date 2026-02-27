// Package feature provides OpenFeature wiring and helpers for go-service.
//
// This package integrates the OpenFeature Go SDK into go-service's DI/lifecycle model by:
//   - optionally registering an OpenFeature FeatureProvider at application start, and
//   - shutting down the OpenFeature SDK at application stop.
//
// # Provider registration and enablement
//
// Provider registration is intentionally optional. The `Register` function is wired via `Module` and
// is a no-op when no `openfeature.FeatureProvider` is available in the DI graph (i.e. it is not
// provided by the consuming service).
//
// When a provider is present, `Register` appends lifecycle hooks that:
//   - call `openfeature.SetProviderAndWait` during application start, and
//   - call `openfeature.Shutdown` during application stop.
//
// # Telemetry hooks
//
// When a metrics provider is available, `Register` installs OpenTelemetry hooks for OpenFeature so
// evaluations can emit metrics and traces.
//
// # Clients
//
// `NewClient` constructs an OpenFeature client named after the service (via env.Name). Callers can
// use this client to evaluate feature flags.
//
// Start with `Module`, `Register`, and `NewClient`.
package feature
