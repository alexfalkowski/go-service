package errors_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/errors"
	"github.com/stretchr/testify/require"
)

func TestServerCloseClassifiesExpectedShutdown(t *testing.T) {
	require.NoError(t, errors.ServerError(http.ErrServerClosed))
	require.Error(t, errors.ServerError(test.ErrFailed))
}
