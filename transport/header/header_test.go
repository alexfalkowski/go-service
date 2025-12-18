package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/stretchr/testify/require"
)

func TestValidParseAuthorization(t *testing.T) {
	key, value, err := header.ParseAuthorization("Bearer token")
	require.NoError(t, err)
	require.Equal(t, "Bearer", key)
	require.Equal(t, "token", value)
}

func TestMissingParseAuthorization(t *testing.T) {
	_, _, err := header.ParseAuthorization(strings.Empty)
	require.ErrorIs(t, header.ErrInvalidAuthorization, err)
}

func TestNotSupportedParseAuthorization(t *testing.T) {
	_, _, err := header.ParseAuthorization("Bob token")
	require.ErrorIs(t, header.ErrNotSupportedAuthorization, err)
}
