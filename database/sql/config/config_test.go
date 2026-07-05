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
			cfg: &config.Config{
				Reader: &config.Pool{Settings: &config.PoolSettings{
					ConnMaxLifetime: time.Second,
					ConnMaxIdleTime: time.Second,
					MaxOpenConns:    2,
					MaxIdleConns:    1,
				}},
				Writer: &config.Pool{Settings: &config.PoolSettings{MaxOpenConns: 1, MaxIdleConns: 1}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireConfigValidation(t, tt.cfg, tt.err)
		})
	}
}

func TestPoolValidation(t *testing.T) {
	tests := []struct {
		cfg  *config.Config
		name string
		err  bool
	}{
		{
			name: "missing reader pool settings",
			cfg: &config.Config{
				Reader: &config.Pool{},
			},
			err: true,
		},
		{
			name: "reader pool negative conn max lifetime",
			cfg: &config.Config{
				Reader: &config.Pool{
					Settings: &config.PoolSettings{ConnMaxLifetime: -time.Second, MaxOpenConns: 1, MaxIdleConns: 1},
				},
			},
			err: true,
		},
		{
			name: "writer pool max idle conns greater than max open conns",
			cfg: &config.Config{
				Writer: &config.Pool{Settings: &config.PoolSettings{MaxOpenConns: 1, MaxIdleConns: 2}},
			},
			err: true,
		},
		{
			name: "writer pool zero max open conns",
			cfg: &config.Config{
				Writer: &config.Pool{Settings: &config.PoolSettings{MaxOpenConns: 0, MaxIdleConns: 1}},
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireConfigValidation(t, tt.cfg, tt.err)
		})
	}
}

func TestConfigIsEnabled(t *testing.T) {
	require.False(t, (*config.Config)(nil).IsEnabled())
	require.True(t, (&config.Config{}).IsEnabled())
}

func requireConfigValidation(t *testing.T, cfg *config.Config, wantErr bool) {
	t.Helper()

	err := test.Validator.Struct(cfg)
	if wantErr {
		require.Error(t, err)
		return
	}

	require.NoError(t, err)
}
