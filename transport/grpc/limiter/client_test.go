package limiter_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	grpclimiter "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestClientStreamInterceptorCancelsContextOnTerminalError(t *testing.T) {
	want := errors.New("terminal")

	t.Run("open", func(t *testing.T) {
		interceptor := newClient(t).StreamInterceptor()
		var streamContext context.Context
		streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
			streamContext = ctx

			return nil, want
		}

		stream, err := interceptor(clientContext(t), &grpc.StreamDesc{}, nil, "/greet.v1.GreeterService/SayStreamHello", streamer)

		require.Nil(t, stream)
		require.ErrorIs(t, err, want)
		require.ErrorIs(t, streamContext.Err(), context.Canceled)
	})

	t.Run("receive", func(t *testing.T) {
		stream := requireClientStream(t, &testClientStream{recvErr: want})

		err := stream.RecvMsg(nil)

		require.ErrorIs(t, err, want)
		require.ErrorIs(t, stream.Context().Err(), context.Canceled)
	})

	t.Run("send", func(t *testing.T) {
		stream := requireClientStream(t, &testClientStream{sendErr: want})

		err := stream.SendMsg(nil)

		require.ErrorIs(t, err, want)
		require.ErrorIs(t, stream.Context().Err(), context.Canceled)
	})
}

func newClient(t *testing.T) *grpclimiter.Client {
	t.Helper()

	client, err := grpclimiter.NewClient(fxtest.NewLifecycle(t), limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 10))
	require.NoError(t, err)

	return client
}

func requireClientStream(t *testing.T, clientStream *testClientStream) grpc.ClientStream {
	t.Helper()

	interceptor := newClient(t).StreamInterceptor()
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream.ctx = ctx

		return clientStream, nil
	}
	stream, err := interceptor(clientContext(t), &grpc.StreamDesc{}, nil, "/greet.v1.GreeterService/SayStreamHello", streamer)
	require.NoError(t, err)

	return stream
}

func clientContext(t *testing.T) context.Context {
	t.Helper()

	return meta.WithAttributes(t.Context(), meta.WithUserAgent(meta.String("test-agent")))
}

type testClientStream struct {
	grpc.ClientStream
	ctx     context.Context
	recvErr error
	sendErr error
}

func (s *testClientStream) Context() context.Context {
	return s.ctx
}

func (s *testClientStream) RecvMsg(any) error {
	return s.recvErr
}

func (s *testClientStream) SendMsg(any) error {
	return s.sendErr
}
