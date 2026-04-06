package net

import (
	"net"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Conn is an alias for net.Conn.
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Conn = net.Conn

// Dialer is an alias for net.Dialer.
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Dialer = net.Dialer

// Listener is an alias for net.Listener.
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Listener = net.Listener

// Listen creates a listener bound to address on the given network.
//
// It uses net.ListenConfig so the operation respects ctx cancellation. If ctx is canceled before
// the listener is created, Listen returns an error from the standard library.
func Listen(ctx context.Context, network, address string) (Listener, error) {
	config := &net.ListenConfig{}
	return config.Listen(ctx, network, address)
}

// SplitNetworkAddress splits an address in the form "<network>://<address>".
//
// Example:
//
//	SplitNetworkAddress("tcp://localhost:3000") // -> ("tcp", "localhost:3000", true)
//
// It returns ok=false when the separator "://" is not present.
func SplitNetworkAddress(address string) (string, string, bool) {
	return strings.Cut(address, "://")
}

// ListenNetworkAddress resolves address into the network/address pair used by Listen.
//
// If address uses the go-service "<network>://<address>" convention, the parsed network and address are returned.
// Otherwise, the input is treated as a raw listen address and the "tcp" network is used.
func ListenNetworkAddress(address string) (string, string) {
	network, addr, ok := SplitNetworkAddress(address)
	if ok {
		return network, addr
	}

	return "tcp", address
}

// Host returns the host portion of addr if it is in host:port form.
//
// If addr cannot be parsed by net.SplitHostPort (for example it has no port, or it is malformed),
// Host returns addr unchanged.
func Host(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

// DefaultAddress returns a go-service server address in the form "tcp://:<port>".
//
// This is commonly used as a default bind address when only a port is configured.
func DefaultAddress(port string) string {
	return "tcp://:" + port
}
