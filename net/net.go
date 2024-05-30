package net

import (
	"context"
	"net"
	"time"
)

// Listener for net.
func Listener(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

// DialContext for net.
func DialContext(_ context.Context, network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, time.Minute)
}
