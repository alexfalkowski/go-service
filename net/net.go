package net

import (
	"net"
)

// OutboundIP of the machine.
// nolint:forcetypeassert
func OutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)

	return addr.IP, nil
}

// SplitHostPort splits a network address of the form "host:port".
func SplitHostPort(hostport string) (string, string, error) {
	return net.SplitHostPort(hostport)
}
