// Package meta provides HTTP-specific context metadata helpers and middleware for go-service.
//
// This package serves two related purposes for HTTP request handling:
//
//   - It exposes small convenience wrappers around the generic `meta` package for exporting
//     context-scoped attributes as string maps suitable for logging and header propagation
//     (for example CamelStrings).
//
//   - It provides a small context-backed store for request-scoped HTTP objects used by go-service
//     handlers and middleware, including:
//
//   - the incoming `*http.Request`
//
//   - the active `http.ResponseWriter`
//
//   - the negotiated `encoding.Encoder` (typically selected from the request Content-Type)
//
// # Safety and expectations
//
// Request, Response, and Encoder are intentionally strict helpers: they expect the corresponding values
// to have been stored in the context via WithRequest, WithResponse, and WithEncoder. Calling them without
// those values present will panic due to type assertions.
//
// These helpers are typically used in tightly controlled handler pipelines (for example those created by
// `net/http/content.NewHandler` / `NewRequestHandler`), which populate the context before invoking
// downstream logic.
//
// HTTP server metadata extraction intentionally prefers common forwarding headers
// such as X-Forwarded-For, X-Real-IP, CF-Connecting-IP, and True-Client-IP over
// RemoteAddr. Deployments that use the extracted IP for access logs, policy, or
// rate limiting should place the service behind trusted proxies that strip or
// overwrite spoofed forwarding headers.
//
// This package also provides HTTP metadata middleware via `NewHandler` and `NewRoundTripper`.
package meta
