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
//   - Define parameter structs embedding [In] to declare injected dependencies.
//   - Use [Module] to compose multiple [Option] values into a single module.
//   - Use [Constructor] to provide constructors (Fx Provide).
//   - Use [Decorate] to wrap/modify provided values (Fx Decorate).
//   - Use [Register] to run registration/invocation hooks at startup (Fx Invoke).
//   - Use [Recover] to recover constructor, decorator, and invocation panics.
//   - Use [RootCause] to unwrap Fx/Dig errors to the underlying cause for reporting.
//
// Start with [Module], [Constructor], [Register], [Decorate], [Recover], and [RootCause].
package di
