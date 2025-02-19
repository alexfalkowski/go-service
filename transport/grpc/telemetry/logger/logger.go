package logger

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/transport/strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	service = "grpc"
)

// UnaryServerInterceptor for logger.
func UnaryServerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		fields := []zapcore.Field{
			zap.Stringer(meta.DurationKey, time.Since(start)),
			zap.String(meta.ServiceKey, service),
			zap.String(meta.PathKey, info.FullMethod),
		}

		fields = append(fields, logger.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(meta.CodeKey, code))

		logger.LogWithFunc(CodeToLogFunc(code, log), message(info.FullMethod), err, fields...)

		return resp, err
	}
}

// StreamServerInterceptor for logger.
func StreamServerInterceptor(log *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(p) {
			return handler(srv, stream)
		}

		start := time.Now()
		ctx := stream.Context()
		err := handler(srv, stream)
		fields := []zapcore.Field{
			zap.Stringer(meta.DurationKey, time.Since(start)),
			zap.String(meta.ServiceKey, service),
			zap.String(meta.PathKey, info.FullMethod),
		}

		fields = append(fields, logger.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(meta.CodeKey, code))

		logger.LogWithFunc(CodeToLogFunc(code, log), message(info.FullMethod), err, fields...)

		return err
	}
}

// UnaryClientInterceptor for logger.
func UnaryClientInterceptor(log *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		p := path.Dir(fullMethod)[1:]
		if strings.IsObservable(p) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		start := time.Now()
		err := invoker(ctx, fullMethod, req, resp, conn, opts...)
		fields := []zapcore.Field{
			zap.Stringer(meta.DurationKey, time.Since(start)),
			zap.String(meta.ServiceKey, service),
			zap.String(meta.PathKey, fullMethod),
		}

		fields = append(fields, logger.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(meta.CodeKey, code))

		logger.LogWithFunc(CodeToLogFunc(code, log), message(conn.Target()+fullMethod), err, fields...)

		return err
	}
}

// StreamClientInterceptor for logger.
func StreamClientInterceptor(log *zap.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		p := path.Dir(fullMethod)[1:]
		if strings.IsObservable(p) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		start := time.Now()
		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)
		fields := []zapcore.Field{
			zap.Stringer(meta.DurationKey, time.Since(start)),
			zap.String(meta.ServiceKey, service),
			zap.String(meta.PathKey, fullMethod),
		}

		fields = append(fields, logger.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(meta.CodeKey, code))

		logger.LogWithFunc(CodeToLogFunc(code, log), message(conn.Target()+fullMethod), err, fields...)

		return stream, err
	}
}

// CodeToLogFunc for logger.
//
//nolint:exhaustive
func CodeToLogFunc(code codes.Code, logger *zap.Logger) logger.LogFunc {
	switch code {
	case codes.OK:
		return logger.Info
	case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated,
		codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange:
		return logger.Warn
	default:
		return logger.Error
	}
}

func message(msg string) string {
	return "grpc: get " + msg
}
