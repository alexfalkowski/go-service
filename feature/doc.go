// Package feature provides OpenFeature wiring and helpers for go-service.
//
// This package integrates the OpenFeature SDK with Fx by registering an optional FeatureProvider
// into the application lifecycle. When a metrics provider is available, it also installs OpenTelemetry
// hooks for metrics and traces.
//
// Start with `Module`, `Register`, and `NewClient`.
package feature
