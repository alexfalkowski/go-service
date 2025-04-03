package net

import (
	"net"

	reuse "github.com/libp2p/go-reuseport"
)

// Listener for net.
func Listener(address string) (net.Listener, error) {
	return reuse.Listen("tcp", address)
}

// Host from the address, if it can be split.
func Host(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}

	return host
}
