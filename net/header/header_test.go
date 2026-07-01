package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestValidParseBearer(t *testing.T) {
	t.Parallel()

	value, err := header.ParseBearer("Bearer token")
	require.NoError(t, err)
	require.Equal(t, "token", value)
}

func TestForwardedIPs(t *testing.T) {
	t.Parallel()

	require.Equal(t, [...]header.ForwardedIP{
		{HTTP: "X-Real-Ip", GRPC: "x-real-ip"},
		{HTTP: "CF-Connecting-Ip", GRPC: "cf-connecting-ip"},
		{HTTP: "True-Client-Ip", GRPC: "true-client-ip"},
		{HTTP: "X-Forwarded-For", GRPC: "x-forwarded-for"},
	}, header.ForwardedIPs)
}

func TestValidFieldName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "standard", value: "X-Request-Id", valid: true},
		{name: "http token symbols", value: "X-Rate_Limit~Policy", valid: true},
		{name: "empty", value: strings.Empty},
		{name: "space", value: "Bad Header"},
		{name: "unicode", value: "Café"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.valid, header.ValidFieldName(tt.value))
		})
	}
}

func TestValidFieldValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "plain", value: "Bearer token", valid: true},
		{name: "empty", value: strings.Empty, valid: true},
		{name: "tab", value: "a\tb", valid: true},
		{name: "newline", value: "a\nb"},
		{name: "nul", value: "a\x00b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.valid, header.ValidFieldValue(tt.value))
		})
	}
}

func TestValidParseBearerWithLowercaseScheme(t *testing.T) {
	t.Parallel()

	value, err := header.ParseBearer("bearer token")
	require.NoError(t, err)
	require.Equal(t, "token", value)
}

func TestMissingParseBearer(t *testing.T) {
	t.Parallel()

	_, err := header.ParseBearer(strings.Empty)
	require.ErrorIs(t, err, header.ErrInvalidAuthorization)
}

func TestParseBearerRejectsEmptyToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		header string
		name   string
	}{
		{name: "empty", header: "Bearer "},
		{name: "whitespace", header: "Bearer \t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := header.ParseBearer(tt.header)
			require.ErrorIs(t, err, header.ErrInvalidAuthorization)
		})
	}
}

func TestNotSupportedParseBearer(t *testing.T) {
	t.Parallel()

	_, err := header.ParseBearer("Bob token")
	require.ErrorIs(t, err, header.ErrNotSupportedAuthorization)
}

func TestParseBearerRejectsBasic(t *testing.T) {
	t.Parallel()

	_, err := header.ParseBearer("Basic token")
	require.ErrorIs(t, err, header.ErrNotSupportedAuthorization)
}
