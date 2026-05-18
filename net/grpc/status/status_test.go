package status_test

import (
	"errors"
	"fmt"
	"testing"

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
	err := status.SafeError(codes.Internal, "internal server error", cause)

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Equal(t, "internal server error", s.Message())
	require.ErrorIs(t, err, cause)
}

func TestSafeErrorWrapped(t *testing.T) {
	cause := errors.New("secret database failure")
	err := fmt.Errorf("wrapped: %w", status.SafeError(codes.Internal, "internal server error", cause))

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Contains(t, s.Message(), "wrapped:")
	require.Contains(t, s.Message(), "internal server error")
	require.NotContains(t, s.Message(), "secret database failure")
}

func TestErrorIs(t *testing.T) {
	err := status.Error(codes.NotFound, "missing")
	target := status.Error(codes.NotFound, "missing")

	require.ErrorIs(t, err, target)
}
