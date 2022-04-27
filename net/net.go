package net

import (
	"context"
	"net"

	"github.com/alexfalkowski/go-service/meta"
)

// OutboundIP of the machine.
func OutboundIP(ctx context.Context) string {
	var dialer net.Dialer

	conn, err := dialer.DialContext(ctx, "udp", "8.8.8.8:80")
	if err != nil {
		meta.WithAttribute(ctx, "net.error", err.Error())

		return ""
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String() // nolint:forcetypeassert
}
