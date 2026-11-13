package crypto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/crypto"
	"github.com/stretchr/testify/require"
)

func TestConfigIsEnabledUnlessNil(t *testing.T) {
	require.False(t, (*crypto.Config)(nil).IsEnabled())
	require.True(t, (&crypto.Config{}).IsEnabled())
}
