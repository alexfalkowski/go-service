package errors_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/errors"
	"github.com/stretchr/testify/require"
)

func TestServerClose(t *testing.T) {
	require.NoError(t, errors.ServerError(http.ErrServerClosed))
	require.Error(t, errors.ServerError(test.ErrFailed))
}
