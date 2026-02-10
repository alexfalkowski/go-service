// Package meta provides HTTP-specific metadata helpers for go-service.
//
// This package contains small helpers used by HTTP transports and middleware to:
//
//   - export meta attributes as string maps (for example CamelStrings), and
//   - store and retrieve request-scoped HTTP objects (request/response/encoder) in context.
//
// Note: Request, Response, and Encoder expect the corresponding values to have been stored in the context
// via WithRequest, WithResponse, and WithEncoder. Calling them without those values present will panic
// due to type assertions.
package meta
