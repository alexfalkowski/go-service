package recovery_test

import (
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/recovery"
	"github.com/stretchr/testify/require"
)

func TestServerUnaryInterceptorRecoversPanic(t *testing.T) {
	want := errors.New("recovered")
	interceptor := recovery.NewServer(func(any) error { return want }).UnaryInterceptor()

	_, err := interceptor(t.Context(), nil, &grpc.UnaryServerInfo{}, func(context.Context, any) (any, error) {
		panic("panic")
	})

	require.ErrorIs(t, err, want)
}

func TestServerStreamInterceptorRecoversPanic(t *testing.T) {
	want := errors.New("recovered")
	interceptor := recovery.NewServer(func(any) error { return want }).StreamInterceptor()

	err := interceptor(nil, &test.MetaServerStream{Ctx: t.Context()}, &grpc.StreamServerInfo{}, func(any, grpc.ServerStream) error {
		panic("panic")
	})

	require.ErrorIs(t, err, want)
}
