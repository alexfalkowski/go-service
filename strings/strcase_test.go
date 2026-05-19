package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestToDelimited(t *testing.T) {
	require.Equal(t, "invalid argument", strings.ToDelimited("InvalidArgument", ' '))
}

func TestToLowerCamel(t *testing.T) {
	require.Equal(t, "requestId", strings.ToLowerCamel("request_id"))
}

func TestToSnake(t *testing.T) {
	require.Equal(t, "request_id", strings.ToSnake("requestID"))
}
