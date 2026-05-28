package net_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/stretchr/testify/require"
)

func TestDefaultAddress(t *testing.T) {
	require.Equal(t, "tcp://:9000", net.DefaultAddress("9000"))
}

func TestNetworkAddressValue(t *testing.T) {
	require.Equal(t, "tcp://localhost:9000", net.NetworkAddress("tcp", "localhost:9000"))
}

func TestHost(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{name: "host port", addr: "localhost:9000", want: "localhost"},
		{name: "fallback", addr: "none", want: "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, net.Host(tt.addr))
		})
	}
}

func TestSplitAndJoinHostPort(t *testing.T) {
	host, port, err := net.SplitHostPort("localhost:9000")
	require.NoError(t, err)
	require.Equal(t, "localhost", host)
	require.Equal(t, "9000", port)
	require.Equal(t, "localhost:9000", net.JoinHostPort(host, port))
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

func TestListenNetworkAddress(t *testing.T) {
	network, address := net.ListenNetworkAddress("tcp://localhost:9000")
	require.Equal(t, "tcp", network)
	require.Equal(t, "localhost:9000", address)

	network, address = net.ListenNetworkAddress(":9000")
	require.Equal(t, "tcp", network)
	require.Equal(t, ":9000", address)
}
