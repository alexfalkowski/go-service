package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

// FuzzParseBearer explores Authorization header parsing while preserving supported error classification.
func FuzzParseBearer(f *testing.F) {
	for _, value := range []string{
		"Bearer token",
		"bearer token",
		strings.Empty,
		"Bearer ",
		"Bearer \t",
		"Bob token",
		"Basic token",
		"Bearer token extra",
	} {
		f.Add(value)
	}

	f.Fuzz(func(t *testing.T, value string) {
		token, err := header.ParseBearer(value)
		if err == nil {
			require.NotEmpty(t, token)
			return
		}

		require.Empty(t, token)
		require.True(
			t,
			errors.Is(err, header.ErrInvalidAuthorization) || errors.Is(err, header.ErrNotSupportedAuthorization),
			"unexpected error: %v",
			err,
		)
	})
}
