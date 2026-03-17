// Package meta provides gRPC metadata interceptors and helpers for go-service.
//
// This package extracts incoming request metadata into the context on the server side and injects
// outgoing metadata on the client side. It also provides small helpers around gRPC metadata maps so
// higher-level transport code can build on consistent metadata behavior.
//
// Metadata keys used by this package include: "user-agent", "request-id", "authorization", and "geolocation".
// Server interceptors also set response header metadata such as "service-version" and "request-id".
//
// Start with `UnaryServerInterceptor` / `StreamServerInterceptor` for server-side extraction and
// `UnaryClientInterceptor` / `StreamClientInterceptor` for client-side injection. Use
// `ExtractIncoming` and `ExtractOutgoing` when you need mutable copies of metadata maps.
package meta
