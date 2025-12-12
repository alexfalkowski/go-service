package net_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/stretchr/testify/require"
)

func TestDefaultAddress(t *testing.T) {
	require.Equal(t, "tcp://:9000", net.DefaultAddress("9000"))
}

func TestHost(t *testing.T) {
	require.Equal(t, "none", net.Host("none"))
}

func TestNetworkAddress(t *testing.T) {
	network, address, ok := net.SplitNetworkAddress("tcp://localhost:9000")
	require.True(t, ok)
	require.Equal(t, "tcp", network)
	require.Equal(t, "localhost:9000", address)

	network, address, ok = net.SplitNetworkAddress("no:address")
	require.False(t, ok)
	require.Equal(t, "no:address", network)
	require.Empty(t, address)
}
