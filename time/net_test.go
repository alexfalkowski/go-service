package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewNetworkNilConfig(t *testing.T) {
	net, err := time.NewNetwork(nil)
	require.NoError(t, err)
	require.Nil(t, net)
}

func TestNewNetworkInvalidKind(t *testing.T) {
	_, err := time.NewNetwork(&time.Config{Kind: "invalid"})
	require.Error(t, err)
}

func requireAnyNetworkNow(t *testing.T, configs ...*time.Config) {
	t.Helper()

	var errs []error
	for _, cfg := range configs {
		n, err := time.NewNetwork(cfg)
		require.NoError(t, err)

		_, err = n.Now()
		if err == nil {
			return
		}

		t.Logf("%s %q failed: %v", cfg.Kind, cfg.Address, err)
		errs = append(errs, err)
	}

	require.NoError(t, errors.Join(errs...))
}

func requireNetworkNowError(t *testing.T, cfg *time.Config) {
	t.Helper()

	n, err := time.NewNetwork(cfg)
	require.NoError(t, err)

	_, err = n.Now()
	require.Error(t, err)
}
