package zap

import (
	"context"
	"fmt"
	"path"

	"github.com/alexfalkowski/go-service/pkg/grpc/encoder"
	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	grpcRequest         = "grpc.request"
	grpcResponse        = "grpc.response"
	grpcService         = "grpc.service"
	grpcMethod          = "grpc.method"
	grpcCode            = "grpc.code"
	grpcDuration        = "grpc.duration"
	grpcStartTime       = "grpc.start_time"
	grpcRequestDeadline = "grpc.request.deadline"
	component           = "component"
	grpcComponent       = "grpc"
	healthService       = "grpc.health.v1.Health"
	client              = "client"
	server              = "server"
)

// UnaryServerInterceptor for zap.
func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(ctx, req)
		}

		start := time.Now().UTC()
		method := path.Base(info.FullMethod)
		resp, err := handler(ctx, req)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, time.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String(grpcRequest, encoder.Message(req)),
			zap.String("span.kind", server),
			zap.String(component, grpcComponent),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcRequestDeadline, d.UTC().Format(time.RFC3339)))
		}

		tags := grpcTags.Extract(ctx)
		for k, v := range tags.Values() {
			fields = append(fields, zap.Any(k, v))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())
		loggerLevel := codeToLevel(code, logger)

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		if resp != nil {
			fields = append(fields, zap.String(grpcResponse, encoder.Message(resp)))
		}

		loggerLevel(message, fields...)

		return resp, err
	}
}

// StreamServerInterceptor for zap.
func StreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(srv, stream)
		}

		start := time.Now().UTC()
		ctx := stream.Context()
		method := path.Base(info.FullMethod)
		err := handler(srv, stream)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, time.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String("span.kind", server),
			zap.String(component, grpcComponent),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcRequestDeadline, d.UTC().Format(time.RFC3339)))
		}

		tags := grpcTags.Extract(ctx)
		for k, v := range tags.Values() {
			fields = append(fields, zap.Any(k, v))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())
		loggerLevel := codeToLevel(code, logger)

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		loggerLevel(message, fields...)

		return err
	}
}

// UnaryClientInterceptor for zap.
func UnaryClientInterceptor(logger *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if service == healthService {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		start := time.Now().UTC()
		method := path.Base(fullMethod)
		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, time.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String(grpcRequest, encoder.Message(req)),
			zap.String("span.kind", client),
			zap.String(component, grpcComponent),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcRequestDeadline, d.UTC().Format(time.RFC3339)))
		}

		tags := grpcTags.Extract(ctx)
		for k, v := range tags.Values() {
			fields = append(fields, zap.Any(k, v))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())
		loggerLevel := codeToLevel(code, logger)

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		fields = append(fields, zap.String(grpcResponse, encoder.Message(resp)))

		loggerLevel(message, fields...)

		return err
	}
}

// StreamClientInterceptor for zap.
func StreamClientInterceptor(logger *zap.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if service == healthService {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		start := time.Now().UTC()
		method := path.Base(fullMethod)
		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		fields := []zapcore.Field{
			zap.Int64(grpcDuration, time.ToMilliseconds(time.Since(start))),
			zap.String(grpcStartTime, start.Format(time.RFC3339)),
			zap.String(grpcService, service),
			zap.String(grpcMethod, method),
			zap.String("span.kind", client),
			zap.String(component, grpcComponent),
		}

		for k, v := range meta.Attributes(ctx) {
			fields = append(fields, zap.String(k, v))
		}

		if d, ok := ctx.Deadline(); ok {
			fields = append(fields, zap.String(grpcRequestDeadline, d.UTC().Format(time.RFC3339)))
		}

		tags := grpcTags.Extract(ctx)
		for k, v := range tags.Values() {
			fields = append(fields, zap.Any(k, v))
		}

		code := status.Code(err)
		message := fmt.Sprintf("finished call with code %s", code.String())
		loggerLevel := codeToLevel(code, logger)

		fields = append(fields, zap.Any(grpcCode, code))

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		loggerLevel(message, fields...)

		return stream, err
	}
}

func codeToLevel(code codes.Code, logger *zap.Logger) func(msg string, fields ...zapcore.Field) {
	switch code {
	case codes.OK, codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
		return logger.Info
	case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unavailable:
		return logger.Warn
	case codes.Unknown, codes.Unimplemented, codes.Internal, codes.DataLoss:
		return logger.Error
	default:
		return logger.Error
	}
}
