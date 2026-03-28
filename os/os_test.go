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

func TestSanitizeArgs(t *testing.T) {
	args := []string{"service", "-test.v", "server", "-i", "config.yml", "-test.run=TestName"}
	sanitized := os.SanitizeArgs(args)

	require.Equal(t, []string{"service", "server", "-i", "config.yml"}, sanitized)
	require.Equal(t, []string{"service", "-test.v", "server", "-i", "config.yml", "-test.run=TestName"}, args)
}
