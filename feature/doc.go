// Package feature provides OpenFeature wiring and helpers for go-service.
//
// This package integrates the OpenFeature Go SDK into go-service's DI/lifecycle model by:
//   - optionally registering an OpenFeature FeatureProvider at application start, and
//   - shutting down the OpenFeature SDK at application stop.
//
// # Provider registration and enablement
//
// Provider registration is intentionally optional. The [Register] function is wired via [Module] and
// is a no-op when no openfeature.FeatureProvider is available in the DI graph (i.e. it is not
// provided by the consuming service).
//
// This repository does not construct or connect a built-in OpenFeature provider from [Config].
// Consuming services that need a remote or custom provider should use that config in their own provider
// constructor and provide the resulting openfeature.FeatureProvider to the DI graph.
//
// When a provider is present, [Register] appends lifecycle hooks that:
//   - call openfeature.SetProviderWithContextAndWait during application start, and
//   - call openfeature.ShutdownWithContext during application stop.
//
// Register uses OpenFeature's package-level SDK APIs, so provider registration, hooks, and shutdown are
// process-global effects.
//
// # Telemetry hooks
//
// When a metrics provider is available, [Register] installs OpenTelemetry hooks for OpenFeature so
// evaluations can emit metrics and traces.
//
// Trace event attributes are produced by the upstream OpenFeature OpenTelemetry hook. They may include
// OpenFeature semantic-convention fields such as the targeting key and evaluated flag value. Treat
// feature trace data as operational telemetry and protect exporter/back-end access accordingly.
//
// # Clients
//
// [NewClient] constructs an OpenFeature client named after the service. Callers can
// use this client to evaluate feature flags.
//
// Start with [Module], [Register], and [NewClient].
package feature
