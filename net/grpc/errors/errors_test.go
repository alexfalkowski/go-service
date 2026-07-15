package errors_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/errors"
	"github.com/stretchr/testify/require"
)

func TestServerError(t *testing.T) {
	t.Parallel()

	require.NoError(t, errors.ServerError(fmt.Errorf("server stopped: %w", grpc.ErrServerStopped)))

	err := errors.ServerError(test.ErrFailed)
	require.Same(t, test.ErrFailed, err)
}
