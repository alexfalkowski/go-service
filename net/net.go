package net

import "net"

// Listener for net.
func Listener(address string) (net.Listener, error) {
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
