package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/stretchr/testify/require"
)

//nolint:usetesting
func TestEnv(t *testing.T) {
	key := "__ENV_KEY"

	require.NoError(t, os.Setenv(key, "test"))
	require.Equal(t, "test", os.Getenv(key))
	require.NoError(t, os.Unsetenv(key))
}
