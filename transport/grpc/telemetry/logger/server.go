package logger

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
)

// UnaryServerInterceptor returns a gRPC unary server interceptor that logs the RPC outcome.
//
// Operation RPC methods (health/metrics/etc.) bypass logging (see [github.com/alexfalkowski/go-service/v2/net/grpc/strings.IsOperationMethod]).
//
// Logged attributes include:
//   - system: "grpc"
//   - service/method: derived from the gRPC full method name
//   - duration: wall-clock elapsed time
//   - code: gRPC status code as a string
//
// Log level is derived from the status code (see [CodeToLevel]). The log message includes the full
// method name and, when present, error details.
//
// Operator diagnostics:
// The raw error is intentionally attached to the log record for backend observability. Client-facing
// responses remain controlled by the gRPC status/error path; logs are expected to be protected operator
// telemetry.
func UnaryServerInterceptor(log *Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsOperationMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		service, method := grpc.ParseServiceMethod(info.FullMethod)
		start := time.Now()
		resp, err := handler(ctx, req)

		attrs := make([]logger.Attr, 0, 5)
		attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
		attrs = append(attrs, logger.String(meta.SystemKey, "grpc"))
		attrs = append(attrs, logger.String(meta.ServiceKey, service))
		attrs = append(attrs, logger.String(meta.MethodKey, method))

		code := status.Code(err)
		attrs = append(attrs, logger.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(info.FullMethod), err), attrs...)

		return resp, err
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that logs the RPC outcome.
//
// Operation RPC methods (health/metrics/etc.) bypass logging (see [github.com/alexfalkowski/go-service/v2/net/grpc/strings.IsOperationMethod]).
//
// Logged attributes include:
//   - system: "grpc"
//   - service/method: derived from the gRPC full method name
//   - duration: wall-clock elapsed time
//   - code: gRPC status code as a string
//
// Log level is derived from the status code (see [CodeToLevel]). The log message includes the full
// method name and, when present, error details.
//
// Operator diagnostics:
// The raw error is intentionally attached to the log record for backend observability. Client-facing
// responses remain controlled by the gRPC status/error path; logs are expected to be protected operator
// telemetry.
func StreamServerInterceptor(log *Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if strings.IsOperationMethod(info.FullMethod) {
			return handler(srv, stream)
		}

		service, method := grpc.ParseServiceMethod(info.FullMethod)
		start := time.Now()
		ctx := stream.Context()
		err := handler(srv, stream)

		attrs := make([]logger.Attr, 0, 5)
		attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
		attrs = append(attrs, logger.String(meta.SystemKey, "grpc"))
		attrs = append(attrs, logger.String(meta.ServiceKey, service))
		attrs = append(attrs, logger.String(meta.MethodKey, method))

		code := status.Code(err)
		attrs = append(attrs, logger.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(info.FullMethod), err), attrs...)

		return err
	}
}
