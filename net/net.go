package net

import (
	"context"
	"errors"
	"net"

	"github.com/alexfalkowski/go-service/time"
)

// ErrInvalidPort for server.
var ErrInvalidPort = errors.New("invalid port")

// Listener for server.
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
