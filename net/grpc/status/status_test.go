package status_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
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

func TestNewWithRetryInfo(t *testing.T) {
	s, err := status.New(codes.Unavailable, "retry later").WithDetails(&status.RetryInfo{
		RetryDelay: status.NewDuration(time.Second),
	})
	require.NoError(t, err)

	details := s.Details()
	require.Len(t, details, 1)

	retryInfo, ok := details[0].(*status.RetryInfo)
	require.True(t, ok)
	require.Equal(t, time.Second.Duration(), retryInfo.GetRetryDelay().AsDuration())
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

func TestSafeErrorSurfacesCauseInErrorString(t *testing.T) {
	safe := status.SafeError(codes.Internal, errors.New("secret database failure"))

	// Operator logs render err.Error(); the diagnostic cause must be visible there so failures can be actioned.
	require.Contains(t, safe.Error(), "secret database failure")

	// Returned directly to gRPC, the client-visible wire message stays the safe message.
	s, ok := status.FromError(safe)
	require.True(t, ok)
	require.Equal(t, "grpc: internal", s.Message())

	// A plain (non-safe) error keeps its client-visible message in Error() too.
	plain := status.Error(codes.NotFound, "widget missing")
	require.Contains(t, plain.Error(), "widget missing")
}

func TestLocalError(t *testing.T) {
	cause := errors.New("local limiter rejection")
	err := status.LocalError(status.SafeError(codes.ResourceExhausted, cause))

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.True(t, status.IsLocalError(err))
	require.True(t, status.IsLocalError(fmt.Errorf("wrapped: %w", err)))
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	require.Equal(t, codes.ResourceExhausted, s.Code())
	require.Equal(t, "grpc: resource exhausted", s.Message())
	require.ErrorIs(t, err, cause)
}

func TestLocalErrorConvertsPlainError(t *testing.T) {
	cause := errors.New("local failure")
	err := status.LocalError(cause)

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, codes.Unknown, s.Code())
	require.Equal(t, cause.Error(), s.Message())
	require.ErrorIs(t, err, cause)
}

func TestLocalErrorAllowsNil(t *testing.T) {
	require.NoError(t, status.LocalError(nil))
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

	unwrapper, ok := err.(unwrapper)
	require.True(t, ok)
	require.Contains(t, unwrapper.Unwrap().Error(), "load tenant tenant-1")
}

func TestSafeErrorfWithoutCause(t *testing.T) {
	err := status.SafeErrorf(codes.Internal, nil, "load tenant %s", "tenant-1")

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Equal(t, "grpc: internal", s.Message())

	unwrapper, ok := err.(unwrapper)
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

	unwrapper, ok := err.(unwrapper)
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

	// Wrapping a SafeError with fmt.Errorf before returning it to gRPC is unsupported: upstream
	// status.FromError rebuilds the wire message from the outer err.Error(), which now carries the
	// diagnostic cause. This is the deliberate "you know what you are doing" escape hatch; use SafeErrorf
	// (see TestSafeErrorf) to add internal context without surfacing the cause on the wire.
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, s.Code())
	require.Contains(t, s.Message(), "wrapped:")
	require.Contains(t, s.Message(), "secret database failure")
}

func TestErrorIs(t *testing.T) {
	err := status.Error(codes.NotFound, "missing")
	target := status.Error(codes.NotFound, "missing")

	require.ErrorIs(t, err, target)
	require.ErrorIs(t, err, status.New(codes.NotFound, "missing").Err())
}

type unwrapper interface {
	Unwrap() error
}
