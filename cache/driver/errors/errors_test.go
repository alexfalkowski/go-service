package errors_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	goerrors "github.com/alexfalkowski/go-service/v2/errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestErrorClassificationIdentifiesCacheErrorKinds(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		missing bool
		expired bool
	}{
		{name: "missing", err: errors.ErrMissing, missing: true},
		{name: "redis missing", err: redis.Nil, missing: true},
		{name: "wrapped redis missing", err: goerrors.Prefix("wrapped", redis.Nil), missing: true},
		{name: "expired", err: errors.ErrExpired, expired: true},
		{name: "wrapped expired", err: goerrors.Prefix("wrapped", errors.ErrExpired), expired: true},
		{name: "nil"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.missing, errors.IsMissingError(test.err))
			require.Equal(t, test.expired, errors.IsExpiredError(test.err))
		})
	}
}
