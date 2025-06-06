package logger

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// Logger is an alias of logger.Logger.
type Logger = logger.Logger

const service = "grpc"

// UnaryServerInterceptor for logger.
func UnaryServerInterceptor(log *Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		attrs := []slog.Attr{
			slog.String(meta.DurationKey, time.Since(start).String()),
			slog.String(meta.ServiceKey, service),
			slog.String(meta.PathKey, p),
		}

		code := status.Code(err)
		attrs = append(attrs, slog.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(info.FullMethod), err), attrs...)

		return resp, err
	}
}

// StreamServerInterceptor for logger.
func StreamServerInterceptor(log *Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(srv, stream)
		}

		start := time.Now()
		ctx := stream.Context()
		err := handler(srv, stream)
		attrs := []slog.Attr{
			slog.String(meta.DurationKey, time.Since(start).String()),
			slog.String(meta.ServiceKey, service),
			slog.String(meta.PathKey, p),
		}

		code := status.Code(err)
		attrs = append(attrs, slog.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(info.FullMethod), err), attrs...)

		return err
	}
}

// UnaryClientInterceptor for logger.
func UnaryClientInterceptor(log *Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		p := fullMethod[1:]
		if strings.IsObservable(p) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		start := time.Now()
		err := invoker(ctx, fullMethod, req, resp, conn, opts...)
		attrs := []slog.Attr{
			slog.String(meta.DurationKey, time.Since(start).String()),
			slog.String(meta.ServiceKey, service),
			slog.String(meta.PathKey, fullMethod),
		}

		code := status.Code(err)
		attrs = append(attrs, slog.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(conn.Target()+fullMethod), err), attrs...)

		return err
	}
}

// StreamClientInterceptor for logger.
func StreamClientInterceptor(log *Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		p := fullMethod[1:]
		if strings.IsObservable(p) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		start := time.Now()
		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)
		attrs := []slog.Attr{
			slog.String(meta.DurationKey, time.Since(start).String()),
			slog.String(meta.ServiceKey, service),
			slog.String(meta.PathKey, fullMethod),
		}

		code := status.Code(err)
		attrs = append(attrs, slog.String(meta.CodeKey, code.String()))

		log.LogAttrs(ctx, CodeToLevel(code), logger.NewMessage(message(conn.Target()+fullMethod), err), attrs...)

		return stream, err
	}
}

// CodeToLevel translates a gRPC status code to a slog.Level.
func CodeToLevel(code codes.Code) slog.Level {
	switch code {
	case codes.OK:
		return slog.LevelInfo
	case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated,
		codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

func message(msg string) string {
	return "grpc: get " + msg
}
