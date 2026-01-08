package logger

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// LevelError is an alias of logger.LevelError.
const LevelError = logger.LevelError

// Logger is an alias of logger.Logger.
type Logger = logger.Logger

// UnaryServerInterceptor for logger.
func UnaryServerInterceptor(log *Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsIgnorable(info.FullMethod) {
			return handler(ctx, req)
		}

		service, method, _ := strings.SplitServiceMethod(info.FullMethod)
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

// StreamServerInterceptor for logger.
func StreamServerInterceptor(log *Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if strings.IsIgnorable(info.FullMethod) {
			return handler(srv, stream)
		}

		service, method, _ := strings.SplitServiceMethod(info.FullMethod)
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

// UnaryClientInterceptor for logger.
func UnaryClientInterceptor(log *Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if strings.IsIgnorable(fullMethod) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		service, method, _ := strings.SplitServiceMethod(fullMethod)
		start := time.Now()
		err := invoker(ctx, fullMethod, req, resp, conn, opts...)

		attrs := make([]logger.Attr, 0, 5)
		attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
		attrs = append(attrs, logger.String(meta.SystemKey, "grpc"))
		attrs = append(attrs, logger.String(meta.ServiceKey, service))
		attrs = append(attrs, logger.String(meta.MethodKey, method))

		code := status.Code(err)
		attrs = append(attrs, logger.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(conn.Target()+fullMethod), err), attrs...)

		return err
	}
}

// StreamClientInterceptor for logger.
func StreamClientInterceptor(log *Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if strings.IsIgnorable(fullMethod) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		service, method, _ := strings.SplitServiceMethod(fullMethod)
		start := time.Now()
		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)

		attrs := make([]logger.Attr, 0, 5)
		attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
		attrs = append(attrs, logger.String(meta.SystemKey, "grpc"))
		attrs = append(attrs, logger.String(meta.ServiceKey, service))
		attrs = append(attrs, logger.String(meta.MethodKey, method))

		code := status.Code(err)
		attrs = append(attrs, logger.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(conn.Target()+fullMethod), err), attrs...)

		return stream, err
	}
}

// CodeToLevel translates a gRPC status code to a logger.Level.
func CodeToLevel(code codes.Code) logger.Level {
	switch code {
	case codes.OK:
		return logger.LevelInfo
	case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated,
		codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange:
		return logger.LevelWarn
	default:
		return logger.LevelError
	}
}

func message(msg string) string {
	return "grpc: get " + msg
}
