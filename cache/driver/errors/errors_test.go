package errors_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	goerrors "github.com/alexfalkowski/go-service/v2/errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestIsMissingError(t *testing.T) {
	require.True(t, errors.IsMissingError(errors.ErrMissing))
	require.True(t, errors.IsMissingError(redis.Nil))
	require.True(t, errors.IsMissingError(goerrors.Prefix("wrapped", redis.Nil)))
	require.False(t, errors.IsMissingError(errors.ErrExpired))
	require.False(t, errors.IsMissingError(nil))
}

func TestIsExpiredError(t *testing.T) {
	require.True(t, errors.IsExpiredError(errors.ErrExpired))
	require.True(t, errors.IsExpiredError(goerrors.Prefix("wrapped", errors.ErrExpired)))
	require.False(t, errors.IsExpiredError(errors.ErrMissing))
	require.False(t, errors.IsExpiredError(nil))
}
