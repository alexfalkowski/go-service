package net

import (
	"context"
	"net"

	"github.com/alexfalkowski/go-service/meta"
)

// OutboundIP of the machine.
// nolint:forcetypeassert
func OutboundIP(ctx context.Context) string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		meta.WithAttribute(ctx, "net.error", err.Error())

		return ""
	}

	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)

	return addr.IP.String()
}
