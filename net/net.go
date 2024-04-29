package net

import (
	"context"
	"errors"
	"net"
	"syscall"

	"github.com/alexfalkowski/go-service/time"
)

// ErrInvalidPort for net.
var ErrInvalidPort = errors.New("invalid port")

// Listener for net.
func Listener(port string) (net.Listener, error) {
	if port == "" {
		return nil, ErrInvalidPort
	}

	return net.Listen("tcp", ":"+port)
}

// DialContext for net.
func DialContext(_ context.Context, network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, time.Timeout)
}

// IsConnectionRefused returns a boolean indicating whether the error is known to report connection is refused.
func IsConnectionRefused(err error) bool {
	return errors.Is(err, syscall.ECONNREFUSED)
}
