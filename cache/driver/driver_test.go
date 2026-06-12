package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	drivererrors "github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	"github.com/stretchr/testify/require"
)

func TestNewDriver(t *testing.T) {
	tests := []struct {
		config  *config.Config
		err     error
		name    string
		wantNil bool
	}{
		{name: "disabled", wantNil: true},
		{name: "ttlcache", config: &config.Config{Kind: "ttlcache", MaxEntries: config.DefaultMaxEntries}},
		{name: "unknown", config: &config.Config{Kind: "unknown"}, err: drivererrors.ErrNotFound, wantNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := driver.NewDriver(driver.DriverParams{Config: tt.config})
			require.ErrorIs(t, err, tt.err)
			if tt.wantNil {
				require.Nil(t, d)
				return
			}

			require.NotNil(t, d)
		})
	}
}
