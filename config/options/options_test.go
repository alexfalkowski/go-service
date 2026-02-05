package options_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	require.Equal(t, 10*time.Minute, test.ConfigOptions.Duration("read_timeout", 5*time.Second))
	require.Equal(t, 5*time.Second, test.ConfigOptions.Duration("timeout", 5*time.Second))
	require.Equal(t, 5*time.Second, test.ConfigOptions.Duration("bob", 5*time.Second))
}
