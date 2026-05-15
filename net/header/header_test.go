package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestValidParseBearer(t *testing.T) {
	value, err := header.ParseBearer("Bearer token")
	require.NoError(t, err)
	require.Equal(t, "token", value)
}

func TestForwardedIPs(t *testing.T) {
	require.Equal(t, [...]header.ForwardedIP{
		{HTTP: "X-Real-Ip", GRPC: "x-real-ip"},
		{HTTP: "CF-Connecting-Ip", GRPC: "cf-connecting-ip"},
		{HTTP: "True-Client-Ip", GRPC: "true-client-ip"},
		{HTTP: "X-Forwarded-For", GRPC: "x-forwarded-for"},
	}, header.ForwardedIPs)
}

func TestValidParseBearerWithLowercaseScheme(t *testing.T) {
	value, err := header.ParseBearer("bearer token")
	require.NoError(t, err)
	require.Equal(t, "token", value)
}

func TestMissingParseBearer(t *testing.T) {
	_, err := header.ParseBearer(strings.Empty)
	require.ErrorIs(t, err, header.ErrInvalidAuthorization)
}

func TestNotSupportedParseBearer(t *testing.T) {
	_, err := header.ParseBearer("Bob token")
	require.ErrorIs(t, err, header.ErrNotSupportedAuthorization)
}

func TestParseBearerRejectsBasic(t *testing.T) {
	_, err := header.ParseBearer("Basic token")
	require.ErrorIs(t, err, header.ErrNotSupportedAuthorization)
}
