package status_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	grpcstatus "github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/http"
	httpstatus "github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestIsErrorRecognizesWrappedCoder(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &customCoderError{code: http.StatusConflict})
	require.True(t, httpstatus.IsError(err))
}

func TestCodeRecognizesWrappedCoder(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &customCoderError{code: http.StatusConflict})
	require.Equal(t, http.StatusConflict, httpstatus.Code(err))
}

func TestFromErrorKeepsWrappedCoder(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &customCoderError{code: http.StatusConflict})
	require.Same(t, err, httpstatus.FromError(http.StatusBadRequest, err))
}

func TestFromErrorWrapsCause(t *testing.T) {
	err := httpstatus.BadRequestError(test.ErrInvalid)

	require.ErrorIs(t, err, test.ErrInvalid)
	require.True(t, httpstatus.IsError(err))
	require.Equal(t, http.StatusBadRequest, httpstatus.Code(err))
}

func TestCodeRecognizesMaxBytesError(t *testing.T) {
	err := &http.MaxBytesError{Limit: 1}

	require.Equal(t, http.StatusRequestEntityTooLarge, httpstatus.Code(err))
}

func TestBadRequestErrorUsesRequestEntityTooLargeForMaxBytesError(t *testing.T) {
	err := httpstatus.BadRequestError(&http.MaxBytesError{Limit: 1})

	require.True(t, httpstatus.IsError(err))
	require.Equal(t, http.StatusRequestEntityTooLarge, httpstatus.Code(err))
}

func TestBadRequestErrorWrapsMaxBytesError(t *testing.T) {
	maxBytesErr := &http.MaxBytesError{Limit: 1}
	err := httpstatus.BadRequestError(maxBytesErr)

	require.ErrorIs(t, err, maxBytesErr)
}

func TestCodeMapsGRPCStatusCodes(t *testing.T) {
	for _, tc := range []struct {
		name string
		code codes.Code
		want int
	}{
		{name: "ok", code: codes.OK, want: http.StatusOK},
		{name: "canceled", code: codes.Canceled, want: 499},
		{name: "unknown", code: codes.Unknown, want: http.StatusInternalServerError},
		{name: "invalid argument", code: codes.InvalidArgument, want: http.StatusBadRequest},
		{name: "deadline exceeded", code: codes.DeadlineExceeded, want: http.StatusGatewayTimeout},
		{name: "not found", code: codes.NotFound, want: http.StatusNotFound},
		{name: "already exists", code: codes.AlreadyExists, want: http.StatusConflict},
		{name: "permission denied", code: codes.PermissionDenied, want: http.StatusForbidden},
		{name: "unauthenticated", code: codes.Unauthenticated, want: http.StatusUnauthorized},
		{name: "resource exhausted", code: codes.ResourceExhausted, want: http.StatusTooManyRequests},
		{name: "failed precondition", code: codes.FailedPrecondition, want: http.StatusBadRequest},
		{name: "aborted", code: codes.Aborted, want: http.StatusConflict},
		{name: "out of range", code: codes.OutOfRange, want: http.StatusBadRequest},
		{name: "unimplemented", code: codes.Unimplemented, want: http.StatusNotImplemented},
		{name: "internal", code: codes.Internal, want: http.StatusInternalServerError},
		{name: "unavailable", code: codes.Unavailable, want: http.StatusServiceUnavailable},
		{name: "data loss", code: codes.DataLoss, want: http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := grpcstatus.Error(tc.code, tc.name)

			require.Equal(t, tc.want, httpstatus.Code(err))
		})
	}
}

func TestWriteErrorReturnsWriteFailure(t *testing.T) {
	res := &test.ErrResponseWriter{}

	err := httpstatus.WriteError(res, httpstatus.BadRequestError(test.ErrInvalid))

	require.ErrorIs(t, err, test.ErrFailed)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestWriteTextReturnsWriteFailure(t *testing.T) {
	res := &test.ErrResponseWriter{}

	err := httpstatus.WriteText(res, "hello")

	require.ErrorIs(t, err, test.ErrFailed)
	require.Equal(t, http.StatusOK, res.Code)
}

type customCoderError struct {
	code int
}

func (c *customCoderError) Error() string {
	return "custom"
}

func (c *customCoderError) Code() int {
	return c.code
}
