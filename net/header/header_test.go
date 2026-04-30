package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestValidParseAuthorization(t *testing.T) {
	key, value, err := header.ParseAuthorization("Bearer token")
	require.NoError(t, err)
	require.Equal(t, "Bearer", key)
	require.Equal(t, "token", value)
}

func TestValidParseAuthorizationWithLowercaseScheme(t *testing.T) {
	for _, tc := range []struct {
		name   string
		value  string
		scheme string
	}{
		{name: "bearer", value: "bearer token", scheme: header.BearerAuthorization},
		{name: "basic", value: "basic token", scheme: header.BasicAuthorization},
	} {
		t.Run(tc.name, func(t *testing.T) {
			key, value, err := header.ParseAuthorization(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.scheme, key)
			require.Equal(t, "token", value)
		})
	}
}

func TestMissingParseAuthorization(t *testing.T) {
	_, _, err := header.ParseAuthorization(strings.Empty)
	require.ErrorIs(t, err, header.ErrInvalidAuthorization)
}

func TestNotSupportedParseAuthorization(t *testing.T) {
	_, _, err := header.ParseAuthorization("Bob token")
	require.ErrorIs(t, err, header.ErrNotSupportedAuthorization)
}
