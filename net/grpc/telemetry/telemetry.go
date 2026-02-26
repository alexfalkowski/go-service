package telemetry

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/stats"
)

// NewServerHandler returns a gRPC server stats handler instrumented with
// OpenTelemetry.
//
// The returned value is a google.golang.org/grpc/stats.Handler. When attached to
// a gRPC server (for example via grpc.StatsHandler), it observes RPC lifecycle
// events and emits spans and/or metrics based on the provided options and the
// OpenTelemetry SDK configuration in your application.
//
// This function delegates to otelgrpc.NewServerHandler. For the exact semantics
// and supported otelgrpc.Option values, consult the upstream otelgrpc
// documentation for the version you vend.
func NewServerHandler(opts ...otelgrpc.Option) stats.Handler {
	return otelgrpc.NewServerHandler(opts...)
}

// NewClientHandler returns a gRPC client stats handler instrumented with
// OpenTelemetry.
//
// The returned value is a google.golang.org/grpc/stats.Handler. When attached to
// a gRPC client connection (for example via grpc.WithStatsHandler), it observes
// outbound RPC lifecycle events and emits spans and/or metrics based on the
// provided options and the OpenTelemetry SDK configuration in your application.
//
// This function delegates to otelgrpc.NewClientHandler. For the exact semantics
// and supported otelgrpc.Option values, consult the upstream otelgrpc
// documentation for the version you vend.
func NewClientHandler(opts ...otelgrpc.Option) stats.Handler {
	return otelgrpc.NewClientHandler(opts...)
}
