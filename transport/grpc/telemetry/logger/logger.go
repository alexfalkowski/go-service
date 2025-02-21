package logger

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/transport/strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const service = "grpc"

// UnaryServerInterceptor for logger.
func UnaryServerInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		fields := []logger.Field{
			logger.Stringer(meta.DurationKey, time.Since(start)),
			logger.String(meta.ServiceKey, service),
			logger.String(meta.PathKey, info.FullMethod),
		}

		code := status.Code(err)
		fields = append(fields, logger.Any(meta.CodeKey, code))

		log.LogFunc(ctx, CodeToLogFunc(code, log), message(info.FullMethod), err, fields...)

		return resp, err
	}
}

// StreamServerInterceptor for logger.
func StreamServerInterceptor(log *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(p) {
			return handler(srv, stream)
		}

		start := time.Now()
		ctx := stream.Context()
		err := handler(srv, stream)
		fields := []logger.Field{
			logger.Stringer(meta.DurationKey, time.Since(start)),
			logger.String(meta.ServiceKey, service),
			logger.String(meta.PathKey, info.FullMethod),
		}

		code := status.Code(err)
		fields = append(fields, logger.Any(meta.CodeKey, code))

		log.LogFunc(ctx, CodeToLogFunc(code, log), message(info.FullMethod), err, fields...)

		return err
	}
}

// UnaryClientInterceptor for logger.
func UnaryClientInterceptor(log *logger.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		p := path.Dir(fullMethod)[1:]
		if strings.IsObservable(p) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		start := time.Now()
		err := invoker(ctx, fullMethod, req, resp, conn, opts...)
		fields := []logger.Field{
			logger.Stringer(meta.DurationKey, time.Since(start)),
			logger.String(meta.ServiceKey, service),
			logger.String(meta.PathKey, fullMethod),
		}

		code := status.Code(err)
		fields = append(fields, logger.Any(meta.CodeKey, code))

		log.LogFunc(ctx, CodeToLogFunc(code, log), message(conn.Target()+fullMethod), err, fields...)

		return err
	}
}

// StreamClientInterceptor for logger.
func StreamClientInterceptor(log *logger.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		p := path.Dir(fullMethod)[1:]
		if strings.IsObservable(p) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		start := time.Now()
		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)
		fields := []logger.Field{
			logger.Stringer(meta.DurationKey, time.Since(start)),
			logger.String(meta.ServiceKey, service),
			logger.String(meta.PathKey, fullMethod),
		}

		code := status.Code(err)
		fields = append(fields, logger.Any(meta.CodeKey, code))

		log.LogFunc(ctx, CodeToLogFunc(code, log), message(conn.Target()+fullMethod), err, fields...)

		return stream, err
	}
}

// CodeToLogFunc for logger.
//
//nolint:exhaustive
func CodeToLogFunc(code codes.Code, logger *logger.Logger) logger.LogFunc {
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
