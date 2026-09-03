package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/stretchr/testify/require"
)

func TestSanitizeArgsRemovesGoTestArguments(t *testing.T) {
	args := []string{"service", "-test.v", "server", "-config", "config.yml", "-test.run=TestName", "-test-mode"}
	sanitized := os.SanitizeArgs(args)

	require.Equal(t, []string{"service", "server", "-config", "config.yml", "-test-mode"}, sanitized)
	require.Equal(t, []string{"service", "-test.v", "server", "-config", "config.yml", "-test.run=TestName", "-test-mode"}, args)
}
