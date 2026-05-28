package transport_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/stretchr/testify/require"
)

func TestConfigIsEnabled(t *testing.T) {
	require.False(t, (*transport.Config)(nil).IsEnabled())
	require.True(t, (&transport.Config{}).IsEnabled())
}
