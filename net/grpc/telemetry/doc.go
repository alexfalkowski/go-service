// Package telemetry provides minimal, stable helpers for wiring OpenTelemetry
// instrumentation into gRPC clients and servers using gRPC stats handlers.
//
// This package is intentionally thin and delegates to the upstream OpenTelemetry
// gRPC instrumentation implementation:
//
//	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
//
// The primary goal is to standardize how services attach gRPC telemetry without
// re-exporting the full upstream API surface.
//
// # Attaching telemetry
//
// gRPC telemetry is attached at construction time via a
// google.golang.org/grpc/stats.Handler.
//
// Server-side:
//
//	s := grpc.NewServer(
//		grpc.StatsHandler(telemetry.NewServerHandler(/* otelgrpc options... */)),
//	)
//
// Client-side (Dial):
//
//	cc, err := grpc.DialContext(
//		ctx,
//		target,
//		grpc.WithStatsHandler(telemetry.NewClientHandler(/* otelgrpc options... */)),
//	)
//
// The returned stats handler observes RPC lifecycle events and produces
// OpenTelemetry spans and/or metrics according to the provided options and the
// OpenTelemetry SDK configuration in your application.
//
// # Configuration
//
// The helper constructors accept otelgrpc.Option values. Those options are
// defined by the upstream otelgrpc package and may control (depending on the
// upstream version):
//
//   - which RPCs are instrumented (filters),
//   - how spans are named,
//   - whether and how metrics are produced,
//   - propagation/context behaviors supported by the upstream library.
//
// This package does not configure OpenTelemetry SDK state (exporters, sampling,
// resources, propagators, or global providers). Configure those elsewhere in your
// application.
//
// # Relationship to upstream otelgrpc
//
// NewServerHandler and NewClientHandler are thin wrappers around
// otelgrpc.NewServerHandler and otelgrpc.NewClientHandler. For the exact
// semantics, supported options, and version-specific behavior, consult the
// upstream otelgrpc documentation for the version you vend.
package telemetry
