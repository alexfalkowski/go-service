package config_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/config"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	require.Equal(t, 10*time.Minute, config.Duration(test.ConfigOptions, "read_timeout", 5*time.Second))
	require.Equal(t, 5*time.Second, config.Duration(test.ConfigOptions, "timeout", 5*time.Second))
	require.Equal(t, 5*time.Second, config.Duration(test.ConfigOptions, "bob", 5*time.Second))
}
