package net

import (
	"net"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type (
	// Conn is an alias for net.Conn.
	Conn = net.Conn

	// Dialer is an alias for net.Dialer.
	Dialer = net.Dialer

	// Listener is an alias for net.Listener.
	Listener = net.Listener
)

// Listen is an alias for net.Listen.
func Listen(ctx context.Context, network, address string) (Listener, error) {
	config := &net.ListenConfig{}
	return config.Listen(ctx, network, address)
}

// SplitNetworkAddress takes an address like tcp://localhost:3000 and returns "tcp", "localhost:3000".
func SplitNetworkAddress(address string) (string, string, bool) {
	return strings.Cut(address, "://")
}

// Host from the address, if it can be split.
func Host(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

// DefaultAddress for servers in the form of tcp://:port
func DefaultAddress(port string) string {
	return "tcp://:" + port
}
