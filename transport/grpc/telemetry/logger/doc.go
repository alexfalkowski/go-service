// Package logger provides gRPC logging interceptors and wiring for go-service.
//
// This package integrates request/response logging into gRPC servers and clients via interceptors.
// Server interceptors skip operation methods such as gRPC health checks. Client interceptors log client
// RPC outcomes unless callers add their own filtering in the client interceptor chain.
//
// Logged attributes include system ("grpc"), service/method (derived from the full method name),
// duration, and gRPC status code. Log level is derived from the status code (see CodeToLevel).
//
// Use [NewServer] to construct server-side logging interceptors ([Server.UnaryInterceptor] /
// [Server.StreamInterceptor]) and [NewClient] to construct client-side logging interceptors
// ([Client.UnaryInterceptor] / [Client.StreamInterceptor]).
package logger
