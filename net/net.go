package net

import (
	"context"
	"net"

	"github.com/alexfalkowski/go-service/time"
)

// DialContext for net.
func DialContext(_ context.Context, network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, time.Timeout)
}
