package net

import (
	"net"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type (
	// Conn is an alias for net.Conn, re-exported for convenience.
	Conn = net.Conn

	// Dialer is an alias for net.Dialer, re-exported for convenience.
	Dialer = net.Dialer

	// Listener is an alias for net.Listener, re-exported for convenience.
	Listener = net.Listener
)

// Listen creates a listener using net.ListenConfig so it can respect ctx cancellation.
func Listen(ctx context.Context, network, address string) (Listener, error) {
	config := &net.ListenConfig{}
	return config.Listen(ctx, network, address)
}

// SplitNetworkAddress splits an address like tcp://localhost:3000 into "tcp" and "localhost:3000".
//
// It returns ok=false when the separator "://" is not present.
func SplitNetworkAddress(address string) (string, string, bool) {
	return strings.Cut(address, "://")
}

// Host returns the host portion of addr if it is in host:port form.
//
// If addr cannot be parsed by net.SplitHostPort (for example it has no port), Host returns addr unchanged.
func Host(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

// DefaultAddress returns a server address in the form tcp://:<port>.
func DefaultAddress(port string) string {
	return "tcp://:" + port
}
