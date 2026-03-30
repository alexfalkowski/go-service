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
	value, ok := os.LookupEnv(key)
	require.True(t, ok)
	require.Equal(t, "test", value)
	require.NoError(t, os.Unsetenv(key))
	value, ok = os.LookupEnv(key)
	require.False(t, ok)
	require.Empty(t, value)
}

func TestSanitizeArgs(t *testing.T) {
	args := []string{"service", "-test.v", "server", "-i", "config.yml", "-test.run=TestName", "-test-mode"}
	sanitized := os.SanitizeArgs(args)

	require.Equal(t, []string{"service", "server", "-i", "config.yml", "-test-mode"}, sanitized)
	require.Equal(t, []string{"service", "-test.v", "server", "-i", "config.yml", "-test.run=TestName", "-test-mode"}, args)
}
