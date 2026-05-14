package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/stretchr/testify/require"
)

func TestGetMaxSize(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var cfg *config.Config
		require.Equal(t, config.DefaultMaxSize, cfg.GetMaxSize())
	})

	t.Run("zero", func(t *testing.T) {
		cfg := &config.Config{}
		require.Equal(t, config.DefaultMaxSize, cfg.GetMaxSize())
	})

	t.Run("explicit", func(t *testing.T) {
		cfg := &config.Config{MaxSize: 64}
		require.Equal(t, bytes.Size(64), cfg.GetMaxSize())
	})
}
