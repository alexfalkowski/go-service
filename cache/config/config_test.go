package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetMaxSize(t *testing.T) {
	tests := []struct {
		cfg  *config.Config
		name string
		size bytes.Size
	}{
		{
			name: "nil",
			size: bytes.DefaultSize,
		},
		{
			name: "zero",
			cfg:  &config.Config{},
			size: bytes.DefaultSize,
		},
		{
			name: "explicit",
			cfg:  &config.Config{MaxSize: 64},
			size: bytes.Size(64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.size, tt.cfg.GetMaxSize())
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		cfg  *config.Config
		name string
		err  bool
	}{
		{
			name: "valid",
			cfg:  &config.Config{MaxEntries: config.DefaultMaxEntries},
		},
		{
			name: "valid max size",
			cfg:  &config.Config{MaxSize: bytes.MaxConfigSize, MaxEntries: config.DefaultMaxEntries},
		},
		{
			name: "valid max entries",
			cfg:  &config.Config{MaxEntries: 1},
		},
		{
			name: "negative max size",
			cfg:  &config.Config{MaxSize: -1, MaxEntries: config.DefaultMaxEntries},
			err:  true,
		},
		{
			name: "oversized max size",
			cfg:  &config.Config{MaxSize: bytes.MaxConfigSize + 1, MaxEntries: config.DefaultMaxEntries},
			err:  true,
		},
		{
			name: "zero max entries",
			cfg:  &config.Config{MaxEntries: 0},
			err:  true,
		},
		{
			name: "negative max entries",
			cfg:  &config.Config{MaxEntries: -1},
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := test.Validator.Struct(tt.cfg)
			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
