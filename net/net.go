package net

import (
	"context"
	"net"
	"time"
)

// Listener for net.
func Listener(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}

// DialContext for net.
func DialContext(_ context.Context, network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, time.Minute)
}
