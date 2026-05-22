package status_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/stretchr/testify/require"
)

func TestFromError(t *testing.T) {
	err := status.Error(codes.NotFound, "missing")

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
	require.Equal(t, "missing", s.Message())
}

func TestFromErrorWrapped(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", status.Error(codes.InvalidArgument, "invalid"))

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
	require.Contains(t, s.Message(), "wrapped:")
}

func TestFromErrorUnknown(t *testing.T) {
	s, ok := status.FromError(errors.New("plain error"))
	require.False(t, ok)
	require.Equal(t, codes.Unknown, s.Code())
}

func TestErrorf(t *testing.T) {
	err := status.Errorf(codes.InvalidArgument, "invalid %s", "name")

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
	require.Equal(t, "invalid name", s.Message())
}

func TestSafeError(t *testing.T) {
	cause := errors.New("secret database failure")
	err := status.SafeError(codes.Internal, cause)

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Equal(t, "grpc: internal", s.Message())
	require.ErrorIs(t, err, cause)
}

func TestSafeErrorOK(t *testing.T) {
	err := status.SafeError(codes.OK, errors.New("ignored"))

	require.NoError(t, err)
}

func TestSafeErrorf(t *testing.T) {
	cause := errors.New("secret database failure")
	err := status.SafeErrorf(codes.Internal, cause, "load tenant %s", "tenant-1")

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Equal(t, "grpc: internal", s.Message())
	require.NotContains(t, s.Message(), "tenant-1")
	require.NotContains(t, s.Message(), "secret database failure")
	require.ErrorIs(t, err, cause)

	unwrapper, ok := err.(interface{ Unwrap() error })
	require.True(t, ok)
	require.Contains(t, unwrapper.Unwrap().Error(), "load tenant tenant-1")
}

func TestSafeErrorfWithoutCause(t *testing.T) {
	err := status.SafeErrorf(codes.Internal, nil, "load tenant %s", "tenant-1")

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Equal(t, "grpc: internal", s.Message())

	unwrapper, ok := err.(interface{ Unwrap() error })
	require.True(t, ok)
	require.EqualError(t, unwrapper.Unwrap(), "load tenant tenant-1")
}

func TestSafeErrorfWithoutFormat(t *testing.T) {
	cause := errors.New("secret database failure")
	err := status.SafeErrorf(codes.Internal, cause, "")

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Equal(t, "grpc: internal", s.Message())
	require.ErrorIs(t, err, cause)

	unwrapper, ok := err.(interface{ Unwrap() error })
	require.True(t, ok)
	require.Same(t, cause, unwrapper.Unwrap())
}

func TestSafeErrorfOK(t *testing.T) {
	err := status.SafeErrorf(codes.OK, errors.New("ignored"), "ignored")

	require.NoError(t, err)
}

func TestSafeErrorWrapped(t *testing.T) {
	cause := errors.New("secret database failure")
	err := fmt.Errorf("wrapped: %w", status.SafeError(codes.Internal, cause))

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Contains(t, s.Message(), "wrapped:")
	require.Contains(t, s.Message(), "grpc: internal")
	require.NotContains(t, s.Message(), "secret database failure")
}

func TestErrorIs(t *testing.T) {
	err := status.Error(codes.NotFound, "missing")
	target := status.Error(codes.NotFound, "missing")

	require.ErrorIs(t, err, target)
}
