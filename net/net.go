package net

import "net"

// Listener is an alias for net.Listener.
type Listener = net.Listener

// Listen will reuse a TCP address.
func Listen(address string) (Listener, error) {
	return net.Listen("tcp", address)
}

// Host from the address, if it can be split.
func Host(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}

	return host
}
