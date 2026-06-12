package logger

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
)

// UnaryClientInterceptor returns a gRPC unary client interceptor that logs the RPC outcome.
//
// Logged attributes include:
//   - system: "grpc"
//   - service/method: derived from the gRPC full method name
//   - duration: wall-clock elapsed time
//   - code: gRPC status code as a string
//
// Log level is derived from the status code (see [CodeToLevel]).
//
// The log message prefixes the target address and full method (for example, `conn.Target()+fullMethod`).
//
// Operator diagnostics:
// The raw error is intentionally attached to the log record for backend observability. Logs are expected
// to be protected operator telemetry.
//
// Target logging:
// The raw gRPC client target is intentionally included to identify the configured downstream endpoint in
// operator logs. Client targets are expected to be configuration-controlled service addresses and must not
// contain credentials, tokens, request data, or other secrets.
func UnaryClientInterceptor(log *Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service, method := grpc.ParseServiceMethod(fullMethod)
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

// StreamClientInterceptor returns a gRPC stream client interceptor that logs stream creation and terminal failures.
//
// It logs whether the client stream was opened successfully. When the stream opens successfully, the returned
// stream is wrapped so later RecvMsg and SendMsg failures are logged as terminal stream-operation failures.
//
// Logged attributes include:
//   - system: "grpc"
//   - service/method: derived from the gRPC full method name
//   - duration: wall-clock elapsed time
//   - code: gRPC status code as a string
//
// Log level is derived from the status code (see [CodeToLevel]).
//
// The log message prefixes the target address and full method (for example, `conn.Target()+fullMethod`).
//
// Operator diagnostics:
// The raw error is intentionally attached to the log record for backend observability. Logs are expected
// to be protected operator telemetry.
//
// Target logging:
// The raw gRPC client target is intentionally included to identify the configured downstream endpoint in
// operator logs. Client targets are expected to be configuration-controlled service addresses and must not
// contain credentials, tokens, request data, or other secrets.
func StreamClientInterceptor(log *Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service, method := grpc.ParseServiceMethod(fullMethod)
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

		if err != nil || stream == nil {
			return stream, err
		}

		return &clientStream{
			ClientStream: stream,
			ctx:          ctx,
			log:          log,
			message:      conn.Target() + fullMethod,
			service:      service,
			method:       method,
		}, nil
	}
}

type clientStream struct {
	grpc.ClientStream
	ctx     context.Context
	log     *Logger
	message string
	service string
	method  string
}

func (s *clientStream) RecvMsg(m any) error {
	start := time.Now()
	err := s.ClientStream.RecvMsg(m)
	s.logError("recv", start, err)

	return err
}

func (s *clientStream) SendMsg(m any) error {
	start := time.Now()
	err := s.ClientStream.SendMsg(m)
	s.logError("send", start, err)

	return err
}

func (s *clientStream) logError(operation string, start time.Time, err error) {
	if err == nil || errors.Is(err, io.EOF) {
		return
	}

	attrs := make([]logger.Attr, 0, 6)
	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))
	attrs = append(attrs, logger.String(meta.SystemKey, "grpc"))
	attrs = append(attrs, logger.String(meta.ServiceKey, s.service))
	attrs = append(attrs, logger.String(meta.MethodKey, s.method))
	attrs = append(attrs, logger.String("operation", operation))

	code := status.Code(err)
	attrs = append(attrs, logger.String(meta.CodeKey, code.String()))

	s.log.LogAttrs(s.ctx, CodeToLevel(code), logger.NewMessage(message(s.message), err), attrs...)
}
