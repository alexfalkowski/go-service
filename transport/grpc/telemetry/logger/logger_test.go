package logger_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	grpclogger "github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	"github.com/stretchr/testify/require"
)

func TestUnaryServerInterceptorLogs(t *testing.T) {
	var logs bytes.Buffer
	interceptor := grpclogger.UnaryServerInterceptor(newLogger(&logs))

	resp, err := interceptor(t.Context(), nil, &grpc.UnaryServerInfo{FullMethod: "/greet.v1.GreeterService/SayHello"}, func(context.Context, any) (any, error) {
		return "ok", nil
	})

	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.Contains(t, logs.String(), `"level":"INFO"`)
	require.Contains(t, logs.String(), `"msg":"grpc: /greet.v1.GreeterService/SayHello"`)
	require.Contains(t, logs.String(), `"system":"grpc"`)
	require.Contains(t, logs.String(), `"service":"greet.v1.GreeterService"`)
	require.Contains(t, logs.String(), `"method":"SayHello"`)
	require.Contains(t, logs.String(), `"code":"OK"`)
}

func TestUnaryServerInterceptorSkipsOperationMethod(t *testing.T) {
	var logs bytes.Buffer
	interceptor := grpclogger.UnaryServerInterceptor(newLogger(&logs))
	called := false

	resp, err := interceptor(t.Context(), nil, &grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"}, func(context.Context, any) (any, error) {
		called = true
		return "ok", nil
	})

	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.True(t, called)
	require.Empty(t, logs.String())
}

func TestStreamServerInterceptorLogs(t *testing.T) {
	var logs bytes.Buffer
	interceptor := grpclogger.StreamServerInterceptor(newLogger(&logs))
	stream := &test.MetaServerStream{Ctx: t.Context()}

	err := interceptor(nil, stream, &grpc.StreamServerInfo{FullMethod: "/greet.v1.GreeterService/SayStreamHello"}, func(any, grpc.ServerStream) error {
		return status.Error(codes.NotFound, "missing")
	})

	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
	require.Contains(t, logs.String(), `"level":"WARN"`)
	require.Contains(t, logs.String(), `"msg":"grpc: /greet.v1.GreeterService/SayStreamHello"`)
	require.Contains(t, logs.String(), `"system":"grpc"`)
	require.Contains(t, logs.String(), `"service":"greet.v1.GreeterService"`)
	require.Contains(t, logs.String(), `"method":"SayStreamHello"`)
	require.Contains(t, logs.String(), `"code":"NotFound"`)
	require.Contains(t, logs.String(), `"error":"rpc error: code = NotFound desc = missing"`)
}

func TestUnaryClientInterceptorLogs(t *testing.T) {
	var logs bytes.Buffer
	interceptor := grpclogger.UnaryClientInterceptor(newLogger(&logs))
	conn, err := grpc.NewClient("passthrough:///backend", grpc.WithTransportCredentials(grpc.NewInsecureCredentials()))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	err = interceptor(t.Context(), "/greet.v1.GreeterService/SayHello", nil, nil, conn, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
		return status.Error(codes.Unavailable, "unavailable")
	})

	require.Error(t, err)
	require.Equal(t, codes.Unavailable, status.Code(err))
	require.Contains(t, logs.String(), `"level":"ERROR"`)
	require.Contains(t, logs.String(), `"msg":"grpc: passthrough:///backend/greet.v1.GreeterService/SayHello"`)
	require.Contains(t, logs.String(), `"system":"grpc"`)
	require.Contains(t, logs.String(), `"service":"greet.v1.GreeterService"`)
	require.Contains(t, logs.String(), `"method":"SayHello"`)
	require.Contains(t, logs.String(), `"code":"Unavailable"`)
	require.Contains(t, logs.String(), `"error":"rpc error: code = Unavailable desc = unavailable"`)
}

func TestStreamClientInterceptorLogs(t *testing.T) {
	var logs bytes.Buffer
	interceptor := grpclogger.StreamClientInterceptor(newLogger(&logs))
	conn, err := grpc.NewClient("passthrough:///backend", grpc.WithTransportCredentials(grpc.NewInsecureCredentials()))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	streamer := func(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, status.Error(codes.InvalidArgument, "invalid")
	}

	stream, err := interceptor(t.Context(), &grpc.StreamDesc{ServerStreams: true}, conn, "/greet.v1.GreeterService/SayStreamHello", streamer)

	require.Nil(t, stream)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, logs.String(), `"level":"WARN"`)
	require.Contains(t, logs.String(), `"msg":"grpc: passthrough:///backend/greet.v1.GreeterService/SayStreamHello"`)
	require.Contains(t, logs.String(), `"system":"grpc"`)
	require.Contains(t, logs.String(), `"service":"greet.v1.GreeterService"`)
	require.Contains(t, logs.String(), `"method":"SayStreamHello"`)
	require.Contains(t, logs.String(), `"code":"InvalidArgument"`)
	require.Contains(t, logs.String(), `"error":"rpc error: code = InvalidArgument desc = invalid"`)
}

func TestStreamClientInterceptorLogsStreamOperationError(t *testing.T) {
	t.Run("recv", func(t *testing.T) {
		requireStreamClientOperationError(t,
			&grpc.StreamDesc{ServerStreams: true},
			&clientStream{recvErr: status.Error(codes.Unavailable, "recv unavailable")},
			func(stream grpc.ClientStream) error { return stream.RecvMsg(nil) },
			"ERROR",
			"recv",
			"Unavailable",
			"rpc error: code = Unavailable desc = recv unavailable",
		)
	})

	t.Run("send", func(t *testing.T) {
		requireStreamClientOperationError(t,
			&grpc.StreamDesc{ClientStreams: true},
			&clientStream{sendErr: status.Error(codes.ResourceExhausted, "send exhausted")},
			func(stream grpc.ClientStream) error { return stream.SendMsg(nil) },
			"WARN",
			"send",
			"ResourceExhausted",
			"rpc error: code = ResourceExhausted desc = send exhausted",
		)
	})
}

func TestStreamClientInterceptorDoesNotLogRecvMsgEOF(t *testing.T) {
	var logs bytes.Buffer
	interceptor := grpclogger.StreamClientInterceptor(newLogger(&logs))
	conn, err := grpc.NewClient("passthrough:///backend", grpc.WithTransportCredentials(grpc.NewInsecureCredentials()))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	streamer := func(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error) {
		return &clientStream{recvErr: io.EOF}, nil
	}

	stream, err := interceptor(t.Context(), &grpc.StreamDesc{ServerStreams: true}, conn, "/greet.v1.GreeterService/SayStreamHello", streamer)
	require.NoError(t, err)
	require.NotNil(t, stream)

	err = stream.RecvMsg(nil)

	require.ErrorIs(t, err, io.EOF)
	require.Contains(t, logs.String(), `"code":"OK"`)
	require.NotContains(t, logs.String(), `"operation":"recv"`)
	require.NotContains(t, logs.String(), `"error"`)
}

func TestCodeToLevel(t *testing.T) {
	tests := []struct {
		name  string
		level logger.Level
		code  codes.Code
	}{
		{name: "ok", code: codes.OK, level: logger.LevelInfo},
		{name: "warn", code: codes.InvalidArgument, level: logger.LevelWarn},
		{name: "error", code: codes.DeadlineExceeded, level: logger.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.level, grpclogger.CodeToLevel(tt.code))
		})
	}
}

func newLogger(logs *bytes.Buffer) *grpclogger.Logger {
	return &grpclogger.Logger{Logger: slog.New(slog.NewJSONHandler(logs, &slog.HandlerOptions{}))}
}

func requireStreamClientOperationError(t *testing.T, desc *grpc.StreamDesc, clientStream *clientStream, call func(grpc.ClientStream) error, level, op, code, msg string) {
	t.Helper()

	var logs bytes.Buffer
	interceptor := grpclogger.StreamClientInterceptor(newLogger(&logs))
	conn, err := grpc.NewClient("passthrough:///backend", grpc.WithTransportCredentials(grpc.NewInsecureCredentials()))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	streamer := func(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error) {
		return clientStream, nil
	}

	stream, err := interceptor(t.Context(), desc, conn, "/greet.v1.GreeterService/SayStreamHello", streamer)
	require.NoError(t, err)
	require.NotNil(t, stream)

	err = call(stream)

	require.ErrorIs(t, err, clientStream.err())
	require.Contains(t, logs.String(), `"level":"`+level+`"`)
	require.Contains(t, logs.String(), `"msg":"grpc: passthrough:///backend/greet.v1.GreeterService/SayStreamHello"`)
	require.Contains(t, logs.String(), `"system":"grpc"`)
	require.Contains(t, logs.String(), `"service":"greet.v1.GreeterService"`)
	require.Contains(t, logs.String(), `"method":"SayStreamHello"`)
	require.Contains(t, logs.String(), `"operation":"`+op+`"`)
	require.Contains(t, logs.String(), `"code":"`+code+`"`)
	require.Contains(t, logs.String(), `"error":"`+msg+`"`)
}

type clientStream struct {
	grpc.ClientStream
	recvErr error
	sendErr error
}

func (s *clientStream) RecvMsg(any) error {
	return s.recvErr
}

func (s *clientStream) SendMsg(any) error {
	return s.sendErr
}

func (s *clientStream) err() error {
	if s.recvErr != nil {
		return s.recvErr
	}

	return s.sendErr
}
