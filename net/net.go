package net

import (
	"context"
	"errors"
	"net"
	"syscall"
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

// IsConnectionRefused returns a boolean indicating whether the error is known to report connection is refused.
func IsConnectionRefused(err error) bool {
	return errors.Is(err, syscall.ECONNREFUSED)
}
