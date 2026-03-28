package status_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestIsErrorRecognizesWrappedCoder(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &customCoderError{code: http.StatusConflict})
	require.True(t, status.IsError(err))
}

func TestCodeRecognizesWrappedCoder(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &customCoderError{code: http.StatusConflict})
	require.Equal(t, http.StatusConflict, status.Code(err))
}

func TestFromErrorKeepsWrappedCoder(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &customCoderError{code: http.StatusConflict})
	require.Same(t, err, status.FromError(http.StatusBadRequest, err))
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
