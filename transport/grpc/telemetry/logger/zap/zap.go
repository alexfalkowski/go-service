package zap

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	grpcService   = "grpc.service"
	grpcMethod    = "grpc.method"
	grpcCode      = "grpc.code"
	grpcDuration  = "grpc.duration"
	grpcStartTime = "grpc.start_time"
	grpcDeadline  = "grpc.deadline"
	kind          = "kind"
	grpcKind      = "grpc"
	client        = "client"
	server        = "server"
)

// UnaryServerInterceptor for zap.
func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(ctx, req)
		}

		start := time.Now()
		method := path.Base(info.FullMethod)
		resp, err := handler(ctx, req)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, stime.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String("grpc.kind", server),
			zap.String(kind, grpcKind),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcDeadline, d.Format(time.RFC3339)))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		loggerLevel := codeToLevel(code, logger)
		loggerLevel(message, fields...)

		return resp, err
	}
}

// StreamServerInterceptor for zap.
func StreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(srv, stream)
		}

		start := time.Now()
		ctx := stream.Context()
		method := path.Base(info.FullMethod)
		err := handler(srv, stream)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, stime.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String("grpc.kind", server),
			zap.String(kind, grpcKind),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcDeadline, d.Format(time.RFC3339)))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		loggerLevel := codeToLevel(code, logger)
		loggerLevel(message, fields...)

		return err
	}
}

// UnaryClientInterceptor for zap.
func UnaryClientInterceptor(logger *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		start := time.Now()
		method := path.Base(fullMethod)
		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, stime.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String("grpc.kind", client),
			zap.String(kind, grpcKind),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcDeadline, d.Format(time.RFC3339)))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		loggerLevel := codeToLevel(code, logger)
		loggerLevel(message, fields...)

		return err
	}
}

// StreamClientInterceptor for zap.
func StreamClientInterceptor(logger *zap.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		start := time.Now()
		method := path.Base(fullMethod)
		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, stime.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String("grpc.kind", client),
			zap.String(kind, grpcKind),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcDeadline, d.Format(time.RFC3339)))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		loggerLevel := codeToLevel(code, logger)
		loggerLevel(message, fields...)

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
