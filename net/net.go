package net

import (
	"net"
)

// Listener for net.
func Listener(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}
