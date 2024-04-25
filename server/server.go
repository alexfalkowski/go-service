package server

import (
	"errors"
	"net"
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
