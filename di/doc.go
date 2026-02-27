// Package di provides small wrappers around Uber Fx/Dig to standardize dependency injection wiring.
//
// This package is a thin convenience layer over:
//   - go.uber.org/fx for application lifecycle and dependency injection wiring, and
//   - go.uber.org/dig for error introspection (RootCause).
//
// Most identifiers in this package are type aliases or small wrappers that re-export Fx/Dig concepts
// behind a stable go-service import path. The goal is to keep wiring consistent across the module and
// reduce direct Fx/Dig imports throughout the codebase.
//
// # Common patterns
//
//   - Define parameter structs embedding `di.In` to declare injected dependencies.
//   - Use `di.Module` to compose multiple `di.Option` values into a single module.
//   - Use `di.Constructor` to provide constructors (Fx Provide).
//   - Use `di.Decorate` to wrap/modify provided values (Fx Decorate).
//   - Use `di.Register` to run registration/invocation hooks at startup (Fx Invoke).
//   - Use `di.RootCause` to unwrap Fx/Dig errors to the underlying cause for reporting.
//
// Start with `Module`, `Constructor`, `Register`, `Decorate`, and `RootCause`.
package di
