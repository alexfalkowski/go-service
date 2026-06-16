// Package logger provides gRPC logging interceptors and wiring for go-service.
//
// This package integrates request/response logging into gRPC servers and clients via interceptors.
// Server interceptors skip operation methods such as gRPC health checks. Client interceptors log client
// RPC outcomes unless callers add their own filtering in the client interceptor chain.
//
// Logged attributes include system ("grpc"), service/method (derived from the full method name),
// duration, and gRPC status code. Log level is derived from the status code (see CodeToLevel).
//
// Start with [UnaryServerInterceptor] / [StreamServerInterceptor] for server-side logging and
// [UnaryClientInterceptor] / [StreamClientInterceptor] for client-side logging.
package logger
