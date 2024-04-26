package zap

import (
	"context"
	"path"
	"time"

	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	st "github.com/alexfalkowski/go-service/time"
	tm "github.com/alexfalkowski/go-service/transport/meta"
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

// UnaryServerInterceptor for zap.
func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(p) {
			return handler(ctx, req)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		fields := []zapcore.Field{
			zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
			zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
			zap.String(tm.ServiceKey, service),
			zap.String(tm.PathKey, info.FullMethod),
		}

		fields = append(fields, tz.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(tm.CodeKey, code))

		tz.LogWithFunc(message(info.FullMethod), err, codeToLevel(code, logger), fields...)

		return resp, err
	}
}

// StreamServerInterceptor for zap.
func StreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(p) {
			return handler(srv, stream)
		}

		start := time.Now()
		ctx := stream.Context()
		err := handler(srv, stream)
		fields := []zapcore.Field{
			zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
			zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
			zap.String(tm.ServiceKey, service),
			zap.String(tm.PathKey, info.FullMethod),
		}

		fields = append(fields, tz.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(tm.CodeKey, code))

		tz.LogWithFunc(message(info.FullMethod), err, codeToLevel(code, logger), fields...)

		return err
	}
}

// UnaryClientInterceptor for zap.
func UnaryClientInterceptor(logger *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		p := path.Dir(fullMethod)[1:]
		if strings.IsHealth(p) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		start := time.Now()
		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		fields := []zapcore.Field{
			zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
			zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
			zap.String(tm.ServiceKey, service),
			zap.String(tm.PathKey, fullMethod),
		}

		fields = append(fields, tz.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(tm.CodeKey, code))

		tz.LogWithFunc(message(cc.Target()+fullMethod), err, codeToLevel(code, logger), fields...)

		return err
	}
}

// StreamClientInterceptor for zap.
func StreamClientInterceptor(logger *zap.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		p := path.Dir(fullMethod)[1:]
		if strings.IsHealth(p) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		start := time.Now()
		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		fields := []zapcore.Field{
			zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
			zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
			zap.String(tm.ServiceKey, service),
			zap.String(tm.PathKey, fullMethod),
		}

		fields = append(fields, tz.Meta(ctx)...)

		code := status.Code(err)
		fields = append(fields, zap.Any(tm.CodeKey, code))

		tz.LogWithFunc(message(cc.Target()+fullMethod), err, codeToLevel(code, logger), fields...)

		return stream, err
	}
}

func codeToLevel(code codes.Code, logger *zap.Logger) func(msg string, fields ...zapcore.Field) {
	switch code {
	case codes.OK:
		return logger.Info
	case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated,
		codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange:
		return logger.Warn
	case codes.Unknown, codes.DeadlineExceeded, codes.Unimplemented, codes.Internal, codes.Unavailable, codes.DataLoss:
		return logger.Error
	default:
		return logger.Error
	}
}

func message(msg string) string {
	return "grpc: get " + msg
}
