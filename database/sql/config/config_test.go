package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		cfg  *config.Config
		name string
		err  bool
	}{
		{
			name: "valid",
			cfg:  &config.Config{MaxOpenConns: 1, MaxIdleConns: 1},
		},
		{
			name: "negative conn max lifetime",
			cfg:  &config.Config{ConnMaxLifetime: -time.Second, MaxOpenConns: 1, MaxIdleConns: 1},
			err:  true,
		},
		{
			name: "negative max open conns",
			cfg:  &config.Config{MaxOpenConns: -1, MaxIdleConns: 1},
			err:  true,
		},
		{
			name: "zero max open conns",
			cfg:  &config.Config{MaxOpenConns: 0, MaxIdleConns: 1},
			err:  true,
		},
		{
			name: "negative max idle conns",
			cfg:  &config.Config{MaxOpenConns: 1, MaxIdleConns: -1},
			err:  true,
		},
		{
			name: "zero max idle conns",
			cfg:  &config.Config{MaxOpenConns: 1, MaxIdleConns: 0},
			err:  true,
		},
		{
			name: "max idle conns greater than max open conns",
			cfg:  &config.Config{MaxOpenConns: 1, MaxIdleConns: 2},
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

func TestConfigIsEnabled(t *testing.T) {
	require.False(t, (*config.Config)(nil).IsEnabled())
	require.True(t, (&config.Config{}).IsEnabled())
}
